/*
Copyright (c) 2016 Skuchain,Inc

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

import (
	"encoding/hex"
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
	fmt.Printf("entering INIT\n")
	if len(args) < 1 {
		fmt.Printf("Invalid Init Arg")
		return nil, errors.New("Invalid Init Arg")
	}

	counterSeed := sha256.Sum256([]byte(args[0]))

	err := stub.PutState("CounterSeed", counterSeed[:])

	if err != nil {
		fmt.Printf("Error initializing CounterSeed")
		return nil, errors.New("Error initializing CounterSeed")
	}

	return nil, nil
}

// func (p *Pop) CreateOutput(amount int, assetType string, data string, creatorKeyBytes []byte, creatorSig []byte) error {
// 	//uncertain where does creatorKeyBytes is public key in bytes

// //need pubkey in bytes, signature in bytes, message string
// 	//deserialize public key bytes into a public key object
// 	creatorKey, err := btcec.ParsePubKey(creatorKeyBytes, btcec.S256())

// 	if err != nil {
// 		return fmt.Errorf("Invalid Creator key")
// 	}
// 	//DER is a standard for serialization
// 	//parsing DER signature from bitcoin curve into a signature object
// 	signature, err := btcec.ParseDERSignature(creatorSig, btcec.S256())
// 	if err != nil {
// 		fmt.Printf("Bad Creator signature encoding %+v", p)
// 		return fmt.Errorf("Bad Creator signature encoding %+v", p)
// 	}
// 	//FIXME add Value to the signature
// 	message := hex.EncodeToString(p.Counter) + ":" + p.Address + ":" + strconv.FormatInt(int64(amount), 10) + ":" + assetType + ":" + data

// 	//here we're using the Sum256 hash. I don't remember the distinction from the normal one.
// 	messageBytes := sha256.Sum256([]byte(message))

// 	//try to verify the signature (most likely failure is that the wrong thing has been signed (maybe the counterseed changed or the message you signed and the message you verified are not the same))
// 	success := signature.Verify(messageBytes[:], creatorKey)
// 	if !success {
// 		fmt.Printf("Invalid Creator Signature %+v", p)
// 		return fmt.Errorf("Invalid Creator Signature %+v", p)
// 	}

// 	output := OTX.New(creatorKey, amount, assetType, data, p.Counter)

// 	p.Outputs = append(p.Outputs, *output)
// 	newCounter := sha256.Sum256(p.Counter)
// 	p.Counter = newCounter[:]
// 	return nil
// }

func verify(pubKeyBytes []byte, sigBytes []byte, message string) (success bool, err error) {
	//deserialize public key bytes into a public key object
	creatorKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		return false, fmt.Errorf("Invalid Creator key")
	}

	//DER is a standard for serialization
	//parsing DER signature from bitcoin curve into a signature object
	signature, err := btcec.ParseDERSignature(sigBytes, btcec.S256())
	if err != nil {
		fmt.Printf("Bad Creator signature encoding")
		return false, fmt.Errorf("Bad Creator signature encoding")
	}

	messageBytes := sha256.Sum256([]byte(message))

	//try to verify the signature
	success = signature.Verify(messageBytes[:], creatorKey)
	if !success {
		fmt.Printf("Invalid Creator Signature")
		return false, fmt.Errorf("Invalid Creator Signature")
	}
	return success, nil
}

func (t *IOTRegistry) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) == 0 {
		fmt.Println("Insufficient arguments found")
		return nil, errors.New("Insufficient arguments found")
	}

	argsBytes, err := hex.DecodeString(args[0])
	if err != nil {
		fmt.Println("Invalid argument expected hex")
		return nil, errors.New("Invalid argument expected hex")
	}

	switch function {
	case "registerName":
		//declare and initialize RegisterIdentity struct
		registerNameArgs := IOTRegistryTX.RegisterIdentityTx{}
		err = proto.Unmarshal(argsBytes, &registerNameArgs)
		if err != nil {
			fmt.Printf("Invalid argument expected RegisterNameTX protocol buffer %s", err.Error())
		}

		//check if name is available
		registerNameBytes, err := stub.GetState("Name: " + registerNameArgs.OwnerName)
		if err != nil {
			fmt.Println("Could not get Name State")
			return nil, errors.New("Could not get Name State")
		}

		//if name unavailable
		if len(registerNameBytes) != 0 {
			fmt.Println("Name is unavailable")
			return nil, errors.New("Name is unavailable")
		}

		creatorKeyBytes := registerNameArgs.PubKey

		creatorSig := registerNameArgs.Signature

		message := registerNameArgs.OwnerName + ":" + registerNameArgs.Data

		success, err := verify(creatorKeyBytes, creatorSig, message)
		if err != nil {
			return nil, errors.New("Error verifying signature")
		}

		//marshall into store type. Then put that variable into the state
		store := IOTRegistryStore.Identities{}
		store.OwnerName = registerNameArgs.OwnerName
		store.Pubkey = registerNameArgs.PubKey

		storeBytes, err := proto.Marshal(&store)
		if err != nil {
			fmt.Println(err)
		}

		err = stub.PutState("IdentityName: "+registerNameArgs.OwnerName, storeBytes)
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

		//check if nonce already exists
		nonceCheckBytes, err := stub.GetState("Nonce: " + hex.EncodeToString(registerThingArgs.Nonce))
		if err != nil {
			fmt.Println("Could not get Nonce State")
			return nil, errors.New("Could not get Nonce State")
		}

		//if nonce exists
		if len(nonceCheckBytes) != 0 {
			fmt.Println("Nonce is unavailable")
			return nil, errors.New("Nonce is unavailable")
		}

		//check if owner is valid id (name exists in registry)
		validIDCheckBytes, err := stub.GetState("Name: " + registerThingArgs.OwnerName)
		//looks like OwnerName defined in .proto but not in pb.go
		if err != nil {
			fmt.Println("Could not get OwnerName State")
			return nil, errors.New("Could not get OwnerName State")
		}

		if len(validIDCheckBytes) == 0 {
			fmt.Println("OwnerName is not in registry")
			return nil, errors.New("OwnerName is not in registry")
		}

		for _, identity := range registerThingArgs.Identities {
			aliasCheckBytes, err := stub.GetState("Name: " + identity)
			if err != nil {
				fmt.Printf("Could not get identity: (%s) State", identity)
				return nil, errors.New("Could not get Identity State")
			}
			if len(aliasCheckBytes) != 0 {
				fmt.Printf("Identity: (%s) is already in registry", identity)
				return nil, errors.New("OwnerName is not in registry")
			}
		}

		creatorKeyBytes := registerThingArgs.OwnerName

		creatorSig := registerThingArgs.Signature

		message := registerThingArgs.OwnerName
		for _, identity := range registerThingArgs.Identities {
			message += ":" + identity
		}

		success, err := verify(creatorKeyBytes, creatorSig, message)
		if err != nil {
			return nil, errors.New("Error verifying signature")
		}

		for _, identity := range registerThingArgs.Identities {
			store := IOTRegistryStore.Things{}
			store.Alias = identity
			store.OwnerName = registerThingArgs.OwnerName
			store.Data = registerThingArgs.Data

			storeBytes, err := proto.Marshal(&store)
			if err != nil {
				fmt.Println(err)
			}

			err = stub.PutState("IdentityName: "+identity, storeBytes)
			if err != nil {
				fmt.Printf(err.Error())
				return nil, err
			}
		}
	}

	// change from argument transaction store into
	//new object of store protobuff type
	//

	// case "create":
	// 	createArgs := TuxedoPopsTX.CreateTX{}
	// 	err = proto.Unmarshal(argsBytes, &createArgs)
	// 	if err != nil {
	// 		fmt.Println("Invalid argument expected CreateTX protocol buffer")
	// 		return nil, fmt.Errorf("Invalid argument expected CreateTX protocol buffer %s", err.Error())
	// 	}

	// 	popcodebytes, err := stub.GetState(createArgs.Address)

	// 	if err != nil {
	// 		fmt.Println("Could not get Popcode State")
	// 		return nil, errors.New("Could not get Popcode State")
	// 	}
	// 	popcode := Pop.Pop{}

	// 	if len(popcodebytes) == 0 {
	// 		addrBytes, err := hex.DecodeString(createArgs.Address)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("Invalid popcode address %s ", createArgs.Address)
	// 		}
	// 		hasher := sha256.New()
	// 		hasher.Write(counterseed)
	// 		hasher.Write(addrBytes)
	// 		hashedCounterSeed := []byte{}
	// 		hashedCounterSeed = hasher.Sum(hashedCounterSeed)
	// 		popcode.Counter = hashedCounterSeed[:]
	// 		popcode.Address = hex.EncodeToString(addrBytes)

	// 		err = popcode.CreateOutput(int(createArgs.Amount), createArgs.Type, createArgs.Data, createArgs.CreatorPubKey, createArgs.CreatorSig)
	// 		if err != nil {
	// 			fmt.Printf(err.Error())
	// 			return nil, err
	// 		}

	// 		antiReplayDigest := sha256.Sum256(createArgs.CreatorSig) // WARNING Assumes the Creator sig is not malleable without private key. Need to check if all maleability vectors are checked

	// 		if txCache.Cache[string(antiReplayDigest[:])] {
	// 			fmt.Printf("Already recieved transaction")
	// 			return nil, fmt.Errorf("Already recieved transaction")
	// 		}
	// 		if len(txCache.Cache) > 100 {
	// 			nextseed := sha256.Sum256(counterseed)
	// 			counterseed = nextseed[:]
	// 			txCache.Cache = make(map[string]bool)
	// 		}

	// 	} else {
	// 		err := popcode.FromBytes(popcodebytes)
	// 		if err != nil {
	// 			fmt.Println("Popcode Deserialization error")
	// 			return nil, errors.New("Popcode Deserialization Failure")
	// 		}
	// 		err = popcode.CreateOutput(int(createArgs.Amount), createArgs.Type, createArgs.Data, createArgs.CreatorPubKey, createArgs.CreatorSig)
	// 		if err != nil {
	// 			fmt.Printf(err.Error())
	// 			return nil, err
	// 		}

	// 	}

	// 	sigHash := sha256.Sum256(createArgs.CreatorSig[:])
	// 	cacheIndex := hex.EncodeToString(sigHash[:])
	// 	txCache.Cache[cacheIndex] = true
	// 	err = stub.PutState(createArgs.Address, popcode.ToBytes())
	// 	if err != nil {
	// 		fmt.Printf(err.Error())
	// 		return nil, err
	// 	}

	// case "transfer":
	// 	transferArgs := TuxedoPopsTX.TransferOwners{}
	// 	err = proto.Unmarshal(argsBytes, &transferArgs)
	// 	if err != nil {
	// 		fmt.Println("Invalid argument expected TransferOwners protocol buffer")
	// 		return nil, fmt.Errorf("Invalid argument expected TransferOwners protocol buffer %s", err.Error())
	// 	}
	// 	popcodebytes, err := stub.GetState(transferArgs.Address)
	// 	if err != nil {
	// 		fmt.Println("Could not get Popcode State")
	// 		return nil, errors.New("Could not get Popcode State")
	// 	}
	// 	if len(popcodebytes) == 0 {
	// 		fmt.Println("No value found in popcode")
	// 		return nil, errors.New("No value found in popcode")
	// 	}

	// 	popcode := Pop.Pop{}
	// 	popcode.FromBytes(popcodebytes)
	// 	err = popcode.SetOwner(int(transferArgs.Output), int(transferArgs.Threshold), transferArgs.Data, transferArgs.Owners, transferArgs.PrevOwnerSigs, transferArgs.PopcodePubKey, transferArgs.PopcodeSig)
	// 	if err != nil {
	// 		fmt.Printf(err.Error())
	// 		return nil, err
	// 	}
	// 	err = stub.PutState(transferArgs.Address, popcode.ToBytes())
	// 	if err != nil {
	// 		fmt.Printf(err.Error())
	// 		return nil, err
	// 	}
	// case "unitize":
	// 	unitizeArgs := TuxedoPopsTX.Unitize{}
	// 	err = proto.Unmarshal(argsBytes, &unitizeArgs)
	// 	if err != nil {
	// 		fmt.Println("Invalid argument expected Unitize protocol buffer")
	// 		return nil, fmt.Errorf("Invalid argument expected Unitize protocol buffer %s", err.Error())
	// 	}
	// 	popcodeKeyDigest := sha256.Sum256(unitizeArgs.PopcodePubKey)
	// 	sourceAddress := hex.EncodeToString(popcodeKeyDigest[:20])
	// 	sourcePopcodeBytes, err := stub.GetState(sourceAddress)
	// 	if err != nil {
	// 		fmt.Println("Could not get Popcode State")
	// 		return nil, errors.New("Could not get Popcode State")
	// 	}
	// 	if len(sourcePopcodeBytes) == 0 {
	// 		fmt.Println("No value found in popcode")
	// 		return nil, errors.New("No value found in popcode")
	// 	}
	// 	sourcePopcode := Pop.Pop{}
	// 	err = sourcePopcode.FromBytes(sourcePopcodeBytes)
	// 	if err != nil {
	// 		fmt.Println("Could not get Popcode State")
	// 		return nil, errors.New("Could not get Popcode State")
	// 	}
	// 	destAddress := unitizeArgs.DestAddress
	// 	destPopcodeBytes, err := stub.GetState(destAddress)
	// 	if err != nil {
	// 		return nil, errors.New("Could not get Popcode State")
	// 	}
	// 	destPopcode := Pop.Pop{}
	// 	if len(destPopcodeBytes) == 0 {
	// 		destAddressBytes, err := hex.DecodeString(destAddress)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("Invalid address %s", destAddress)
	// 		}
	// 		hasher := sha256.New()
	// 		hasher.Write(sourcePopcode.Counter)
	// 		hasher.Write(destAddressBytes)
	// 		hashedCounterSeed := []byte{}
	// 		hashedCounterSeed = hasher.Sum(hashedCounterSeed)
	// 		destPopcode.Address = unitizeArgs.DestAddress
	// 		destPopcode.Counter = hashedCounterSeed[:]
	// 	} else {
	// 		err = destPopcode.FromBytes(destPopcodeBytes)
	// 		if err != nil {
	// 			fmt.Println("Dest Popcode Deserialization error")
	// 			return nil, errors.New("Dest Popcode Deserialization Failure")
	// 		}
	// 	}
	// 	convertedAmounts := make([]int, len(unitizeArgs.DestAmounts))
	// 	for i, destAmount := range unitizeArgs.DestAmounts {
	// 		convertedAmounts[i] = int(destAmount)
	// 	}
	// 	sourcePopcode.UnitizeOutput(int(unitizeArgs.SourceOutput), convertedAmounts, unitizeArgs.Data, &destPopcode, unitizeArgs.OwnerSigs, unitizeArgs.PopcodePubKey, unitizeArgs.PopcodeSig)
	// 	err = stub.PutState(sourceAddress, sourcePopcode.ToBytes())
	// 	if err != nil {
	// 		fmt.Printf(err.Error())
	// 		return nil, err
	// 	}
	// 	err = stub.PutState(destAddress, destPopcode.ToBytes())
	// 	if err != nil {
	// 		fmt.Printf(err.Error())
	// 		return nil, err
	// 	}
	// case "combine":
	// 	combineArgs := TuxedoPopsTX.Combine{}

	// 	err = proto.Unmarshal(argsBytes, &combineArgs)
	// 	if err != nil {
	// 		fmt.Println("Invalid argument expected Combine protocol buffer")
	// 		return nil, fmt.Errorf("Invalid argument expected Combine protocol buffer %s", err.Error())
	// 	}

	// 	popcode := Pop.Pop{}
	// 	popcodeBytes, err := stub.GetState(combineArgs.Address)
	// 	if err != nil {
	// 		fmt.Println("Could not get Popcode State")
	// 		return nil, errors.New("Could not get Popcode State")
	// 	}
	// 	if len(popcodeBytes) == 0 {
	// 		fmt.Println("No value found in popcode")
	// 		return nil, errors.New("No value found in popcode")
	// 	}
	// 	popcode.FromBytes(popcodeBytes)

	// 	sources := make([]Pop.SourceOutput, len(combineArgs.Sources))

	// 	for i, v := range combineArgs.Sources {
	// 		sources[i] = v
	// 	}

	// 	popcode.CombineOutputs(sources, combineArgs.OwnerSigs, combineArgs.PopcodePubKey, combineArgs.PopcodeSigs, int(combineArgs.Amount), combineArgs.Type, combineArgs.Data, combineArgs.CreatorPubKey, combineArgs.CreatorSig)
	// 	err = stub.PutState(combineArgs.Address, popcode.ToBytes())
	// 	if err != nil {
	// 		fmt.Printf(err.Error())
	// 		return nil, err
	// 	}

	// default:
	// 	fmt.Printf("Invalid function type")
	// 	return nil, fmt.Errorf("Invalid function type")
	// }
	// txCacheBytes, err = proto.Marshal(&txCache)
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// 	return nil, err
	// }
	// err = stub.PutState("TxCache", txCacheBytes)
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// 	return nil, err
	// }
	// err = stub.PutState("CounterSeed", counterseed)
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// 	return nil, err
	// }
	return nil, nil
}

func (t *IOTRegistry) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// fmt.Printf("function: %s", function)
	// switch function {
	// case "balance":
	// 	if len(args) != 1 {
	// 		return nil, fmt.Errorf("No argument specified")
	// 	}
	// 	counterseed, err := stub.GetState("CounterSeed")

	// 	if err != nil {
	// 		fmt.Printf(err.Error())
	// 		return nil, err
	// 	}

	// 	address := args[0]
	// 	popcode := Pop.Pop{}
	// 	popcodeBytes, err := stub.GetState(address)
	// 	if err != nil {
	// 		fmt.Printf(err.Error())
	// 		return nil, err
	// 	}
	// 	if len(popcodeBytes) == 0 {
	// 		addrBytes, _ := hex.DecodeString(address)
	// 		hasher := sha256.New()
	// 		hasher.Write(counterseed)
	// 		hasher.Write(addrBytes)
	// 		hashedCounterSeed := []byte{}
	// 		hashedCounterSeed = hasher.Sum(hashedCounterSeed)
	// 		popcode.Address = args[0]
	// 		popcode.Counter = hashedCounterSeed
	// 		return popcode.ToJSON(), nil
	// 	}
	// 	popcode.FromBytes(popcodeBytes)
	// 	return popcode.ToJSON(), nil

	// }

	return nil, nil
}

func main() {
	err := shim.Start(new(IOTRegistry))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s\n", err)
	}
}
