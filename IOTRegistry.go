/*
Copyright (c) 2016 Skuchain,Inc

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"crypto/sha256"

	"errors"

	"github.com/btcsuite/btcd/btcec"
	proto "github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/skuchain/IOTRegistry/IOTRegistryStore"
	IOTRegistryTX "github.com/skuchain/IOTRegistry/IOTRegistryTX"
)

// This chaincode implements the ledger operations for the proofchaincode

// ProofChainCode example simple Chaincode implementation
type IOTRegistry struct {
}

func (t *IOTRegistry) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) < 1 {
		fmt.Printf("Invalid Init Arg")
		return nil, fmt.Errorf("Invalid Init Arg: (%s)", args)
	}

	//not sure if next seven lines are necessary
	counterSeed := sha256.Sum256([]byte(args[0]))
	err := stub.PutState("CounterSeed", counterSeed[:])
	if err != nil {
		fmt.Printf("Error initializing CounterSeed")
		return nil, fmt.Errorf("Error initializing CounterSeed: (%s)", counterSeed)
	}
	return nil, nil
}

func verify(pubKeyBytes []byte, sigBytes []byte, message string) (err error) {
	//deserialize public key bytes into a public key object
	creatorKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		return fmt.Errorf("Invalid Creator key")
	}

	//DER is a standard for serialization
	//parsing DER signature from bitcoin curve into a signature object
	signature, err := btcec.ParseDERSignature(sigBytes, btcec.S256())
	if err != nil {
		fmt.Printf("Bad Creator signature encoding")
		return fmt.Errorf("Bad Creator signature encoding")
	}

	messageBytes := sha256.Sum256([]byte(message))

	//try to verify the signature
	success := signature.Verify(messageBytes[:], creatorKey)
	if !success {
		fmt.Printf("Invalid Creator Signature")
		return fmt.Errorf("Invalid Creator Signature")
	}
	return nil
}

func (t *IOTRegistry) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) == 0 {
		fmt.Println("Insufficient arguments found")
		return nil, fmt.Errorf("Insufficient arguments found")
	}
	argsBytes, err := hex.DecodeString(args[0])
	if err != nil {
		fmt.Printf("Invalid argument (%s) expected hex", args[0])
		return nil, fmt.Errorf("Invalid argument (%s) expected hex", args[0])
	}

	switch function {
	case "registerOwner":
		//declare and initialize RegisterIdentity struct
		registerNameArgs := IOTRegistryTX.RegisterIdentityTX{}
		err = proto.Unmarshal(argsBytes, &registerNameArgs)
		if err != nil {
			fmt.Printf("Invalid argument expected RegisterNameTX protocol buffer %s", err.Error())
		}

		if len(registerNameArgs.OwnerName) == 0 {
			return nil, fmt.Errorf("length of OwnerName (%s) is zero", registerNameArgs.OwnerName)
		}
		if len(registerNameArgs.PubKey) == 0 {
			return nil, fmt.Errorf("length of Pubkey (%s) is zero", registerNameArgs.PubKey)
		}
		if len(registerNameArgs.Signature) == 0 {
			return nil, fmt.Errorf("length of Signature (%s) is zero", registerNameArgs.Signature)
		}
		//check if name is available
		registerNameBytes, err := stub.GetState("OwnerName: " + registerNameArgs.OwnerName)
		if err != nil {
			fmt.Println("Could not get OwnerName State")
			return nil, fmt.Errorf("Could not get OwnerName (%s) State", registerNameArgs.OwnerName)
		}

		//if name unavailable
		if len(registerNameBytes) != 0 {
			fmt.Println("OwnerName is unavailable")
			return nil, fmt.Errorf("OwnerName (%s) is unavailable", registerNameArgs.OwnerName)
		}

		creatorKeyBytes := registerNameArgs.PubKey
		creatorSig := registerNameArgs.Signature
		message := registerNameArgs.OwnerName + ":" + registerNameArgs.Data

		err = verify(creatorKeyBytes, creatorSig, message)
		if err != nil {
			return nil, fmt.Errorf("Error verifying signature (%s)", creatorSig)
		}

		//marshall into store type. Then put that variable into the state
		store := IOTRegistryStore.Identities{}
		store.OwnerName = registerNameArgs.OwnerName
		store.Pubkey = registerNameArgs.PubKey
		storeBytes, err := proto.Marshal(&store)
		if err != nil {
			fmt.Println(err)
		}

		err = stub.PutState("OwnerName: "+registerNameArgs.OwnerName, storeBytes)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}

	case "registerThing":
		registerThingArgs := IOTRegistryTX.RegisterThingTX{}
		err = proto.Unmarshal(argsBytes, &registerThingArgs)
		if err != nil {
			fmt.Printf("Invalid argument expected RegisterThingTX protocol buffer %s", err.Error())
		}

		if len(registerThingArgs.OwnerName) == 0 {
			return nil, fmt.Errorf("length of OwnerName (%s) is zero", registerThingArgs.OwnerName)
		}
		if len(registerThingArgs.Nonce) == 0 {
			return nil, fmt.Errorf("length of Nonce (%s) is zero", registerThingArgs.Nonce)
		}
		if len(registerThingArgs.Signature) == 0 {
			return nil, fmt.Errorf("length of Signature (%s) is zero", registerThingArgs.Signature)
		}

		//check if nonce already exists
		nonceCheckBytes, err := stub.GetState("Nonce: " + hex.EncodeToString(registerThingArgs.Nonce))
		if err != nil {
			fmt.Println("Could not get Nonce State")
			return nil, errors.New("Could not get Nonce State")
		}

		//if nonce exists
		if len(nonceCheckBytes) != 0 {
			fmt.Println("Nonce is unavailable")
			return nil, fmt.Errorf("Nonce is unavailable %s", hex.EncodeToString(registerThingArgs.Nonce))
		}

		//check if owner is valid id (name exists in registry)
		checkIDBytes, err := stub.GetState("OwnerName: " + registerThingArgs.OwnerName)
		if err != nil {
			fmt.Println("Failed to look up Owner Name")
			return nil, fmt.Errorf("Failed to look up Owner Name")
		}

		//if owner is not registered name
		if len(checkIDBytes) == 0 {
			fmt.Println("OwnerName is not registered")
			return nil, fmt.Errorf("OwnerName is not registered %s", registerThingArgs.OwnerName)
		}

		//check if any identities exist
		for _, identity := range registerThingArgs.Identities {
			aliasCheckBytes, err := stub.GetState("OwnerName: " + identity)
			if err != nil {
				fmt.Printf("Could not get identity: (%s) State", identity)
				return nil, errors.New("Could not get Identity State")
			}
			//throw error if any of the identities already exist
			if len(aliasCheckBytes) != 0 {
				fmt.Printf("OwnerName: (%s) is already in registry", identity)
				return nil, fmt.Errorf("OwnerName: (%s) is already in registry", identity)
			}
		}

		//retrieve state associated with owner name to get public key
		ownerRegistration := IOTRegistryStore.Identities{}
		err = proto.Unmarshal(checkIDBytes, &ownerRegistration)
		if err != nil {
			return nil, err
		}

		ownerPubKeyBytes := ownerRegistration.Pubkey

		ownerSig := registerThingArgs.Signature

		//TODO review later
		message := registerThingArgs.OwnerName
		for _, identity := range registerThingArgs.Identities {
			message += ":" + identity
		}
		message += ":" + registerThingArgs.Data
		message += ":" + registerThingArgs.Spec
		err = verify(ownerPubKeyBytes, ownerSig, message)
		if err != nil {
			return nil, errors.New("Error verifying signature")
		}

		for _, identity := range registerThingArgs.Identities {

			alias := IOTRegistryStore.Alias{}
			alias.Nonce = registerThingArgs.Nonce
			aliasStoreBytes, err := proto.Marshal(&alias)

			if err != nil {
				return nil, fmt.Errorf("Error marshalling alias (%v) into bytes", alias)
			}
			stub.PutState("Alias:"+identity, aliasStoreBytes)
		}

		store := IOTRegistryStore.Things{}
		store.Alias = registerThingArgs.Identities
		store.OwnerName = registerThingArgs.OwnerName
		store.Data = registerThingArgs.Data
		store.SpecName = registerThingArgs.Spec
		storeBytes, err := proto.Marshal(&store)
		if err != nil {
			fmt.Println(err)
		}
		err = stub.PutState("Thing: "+hex.EncodeToString(registerThingArgs.Nonce), storeBytes)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}
	case "registerSpec":
		specArgs := IOTRegistryTX.RegisterSpecTX{}
		err = proto.Unmarshal(argsBytes, &specArgs)
		if err != nil {
			fmt.Printf("Invalid argument expected RegisterSpecTX protocol buffer %s", err.Error())
		}

		if len(specArgs.OwnerName) == 0 {
			return nil, fmt.Errorf("length of OwnerName (%s) is zero", specArgs.OwnerName)
		}
		if len(specArgs.SpecName) == 0 {
			return nil, fmt.Errorf("length of Nonce (%s) is zero", specArgs.SpecName)
		}
		if len(specArgs.Signature) == 0 {
			return nil, fmt.Errorf("length of Signature (%s) is zero", specArgs.Signature)
		}

		//check if spec already exists
		specNameCheckBytes, err := stub.GetState("Spec: " + specArgs.SpecName)
		if err != nil {
			fmt.Println("Could not get Spec State")
			return nil, errors.New("Could not get Spec State")
		}

		//if spec already exists
		if len(specNameCheckBytes) != 0 {
			fmt.Println("SpecName is unavailable")
			return nil, fmt.Errorf("SpecName (%s) is unavailable", specArgs.SpecName)
		}

		//check if owner is valid id (name exists in registry)
		checkIDBytes, err := stub.GetState("OwnerName: " + specArgs.OwnerName)
		if err != nil {
			fmt.Println("Failed to look up OwnerName")
			return nil, fmt.Errorf("Failed to look up OwnerName (%s)", specArgs.OwnerName)
		}

		//if owner is not registered name
		if len(checkIDBytes) == 0 {
			fmt.Println("OwnerName is not registered")
			return nil, fmt.Errorf("OwnerName is not registered %s", specArgs.OwnerName)
		}

		//retrieve state associated with owner name to get public key
		ownerRegistration := IOTRegistryStore.Identities{}
		err = proto.Unmarshal(checkIDBytes, &ownerRegistration)
		if err != nil {
			return nil, err
		}

		ownerPubKeyBytes := ownerRegistration.Pubkey

		ownerSig := specArgs.Signature

		/*ownerRegistration, a struct of type IOTRegistryStore.Identities, has a field specName.
		should we make sure that the specName input as an argument to invoke spec is equal to this value?
		That is, should we check that specArgs.SpecName == ownerRegistration.SpecName?
		for example:
		if specArgs.SpecName != ownerRegistration.SpecName {
			return nil, fmt.Errorf("mismatched SpecName values (%s) and (%s)",
									specArgs.SpecName, ownerRegistration.SpecName)
		}

		Alternatively, should we just retrieve specName from the ownerRegistration on the ledger?
		If we did this, there would be no need for a field in ownerRegistration for specName.
			-whether a user should be able to register more than one spec seems relevant to this question.
				-if a user can register multiple specs, maybe we can make this a slice of strings.
		*/
		//TODO review later
		message := specArgs.SpecName + ":" + specArgs.OwnerName + ":" + specArgs.Data
		err = verify(ownerPubKeyBytes, ownerSig, message)
		if err != nil {
			return nil, errors.New("Error verifying signature")
		}

		store := IOTRegistryStore.Spec{}
		store.OwnerName = specArgs.OwnerName
		store.Data = specArgs.Data
		storeBytes, err := proto.Marshal(&store)
		if err != nil {
			fmt.Println(err)
		}
		err = stub.PutState("Spec: "+specArgs.SpecName, storeBytes)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}
	}
	return nil, nil
}

func ownerNameToJSON(ownerName string, pubKey []byte) ([]byte, error) {
	type JSONIdentities struct {
		OwnerName string
		Pubkey    string
	}
	jsonOwner := JSONIdentities{}
	jsonOwner.OwnerName = ownerName
	jsonOwner.Pubkey = hex.EncodeToString(pubKey)

	jsonstring, err := json.Marshal(jsonOwner)
	if err != nil {
		return nil, err
	}
	return jsonstring, nil
}

// func ThingToJSON() []byte {
// 	type JSONThings struct {
// 		Alias     []string
// 		OwnerName string
// 		Data      string
// 		SpecName  string
// 	}
// 	jsonThing := JSONThings{}
// 	jsonThing.Address = p.Address
// 	jsonThing.Counter = hex.EncodeToString(p.Counter)
// 	for _, o := range p.Outputs {
// 		jsonThing.Outputs = append(jsonThing.Outputs, string(o.ToJSON()))
// 	}
// 	jsonstring, err := json.Marshal(jsonThing)
// 	if err != nil {
// 		return nil
// 	}
// 	return jsonstring
// }

// func SpecToJSON(address string, ) []byte {
// 	type JSONSpec struct {
// 		OwnerName string
// 		Data      string
// 	}
// 	jsonSpec := JSONSpec{}
// 	jsonSpec.Address = p.Address
// 	jsonSpec.Counter = hex.EncodeToString(p.Counter)
// 	for _, o := range p.Outputs {
// 		jsonSpec.Outputs = append(jsonSpec.Outputs, string(o.ToJSON()))
// 	}
// 	jsonstring, err := json.Marshal(jsonSpec)
// 	if err != nil {
// 		return nil
// 	}
// 	return jsonstring
// }

func (t *IOTRegistry) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// fmt.Printf("function: %s\n", function)
	switch function {
	case "owner":
		if len(args) != 1 {
			return nil, fmt.Errorf("No argument specified")
		}

		owner := IOTRegistryStore.Identities{}

		ownerName := args[0]
		ownerBytes, err := stub.GetState("OwnerName: " + ownerName)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}

		//Owner does not exist \/
		if len(ownerBytes) == 0 {
			return nil, fmt.Errorf("OwnerName (%s) does not exist\n", ownerName)
		}
		err = proto.Unmarshal(ownerBytes, &owner)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}
		jsonBytes, err := ownerNameToJSON(owner.OwnerName, owner.Pubkey)
		return jsonBytes, err
	case "thing":
		if len(args) != 1 {
			return nil, fmt.Errorf("No argument specified")
		}
		thing := IOTRegistryStore.Things{}
		thingNonce := args[0]
		thingBytes, err := stub.GetState("Thing: " + thingNonce)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}

		if len(thingBytes) == 0 {
			return nil, fmt.Errorf("Thing (%s) does not exist\n", thingNonce)
		}

		err = proto.Unmarshal(thingBytes, &thing)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}
		return json.Marshal(thing)
	case "spec":
		if len(args) != 1 {
			return nil, fmt.Errorf("no argument specified")
		}

		spec := IOTRegistryStore.Spec{}
		specName := args[0]

		specBytes, err := stub.GetState("Spec: " + specName)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}

		if len(specBytes) == 0 {
			return nil, fmt.Errorf("spec (%s) does not exist\n", specName)
		}

		err = proto.Unmarshal(specBytes, &spec)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}
		return json.Marshal(spec)
	}
	return nil, nil
}

func main() {
	err := shim.Start(new(IOTRegistry))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s\n", err)
	}
}
