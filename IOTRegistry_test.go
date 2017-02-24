package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"

	proto "github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"github.com/btcsuite/btcd/btcec"
	IOTRegistryTX "github.com/skuchain/IOTRegistry/IOTRegistryTX"
)

// Notes from Testing popcode
// Public Key: 02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc
// Private Key: 94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20
// Hyperledger address hex 74ded2036e988fc56e3cff77a40c58239591e921
// Hyperledger address Base58: 8sDMfw2Ti7YumfTkbf7RHMgSSSxuAmMFd2GS9wnjkUoX

// Notes from Testing popcode2
// Public Key: 02cb6d65b04c4b84502015f918fe549e95cad4f3b899359a170d4d7d438363c0ce
// Private Key: 60977f22a920c9aa18d58d12cb5e90594152d7aa724bcce21484dfd0f4490b58
// Hyperledger address hex 10734390011641497f489cb475743b8e50d429bb
// Hyperledger address Base58: EHxhLN3Ft4p9jPkR31MJMEMee9G

//Owner1 key
// Public Key: 0278b76afbefb1e1185bc63ed1a17dd88634e0587491f03e9a8d2d25d9ab289ee7
// Private Key: 7142c92e6eba38de08980eeb55b8c98bb19f8d417795adb56b6c4d25da6b26c5

// Owner2 key
// Public Key: 02e138b25db2e74c54f8ca1a5cf79e2d1ed6af5bd1904646e7dc08b6d7b0d12bfd
// Private Key: b18b7d3082b3ff9438a7bf9f5f019f8a52fb64647ea879548b3ca7b551eefd65

var hexChars = []rune("0123456789abcdef")
var alpha = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//testing tool for creating randomized string with a certain character makeup
func randString(n int, kindOfString string) string {
	b := make([]rune, n)
	if kindOfString == "hex" {
		for i := range b {
			b[i] = hexChars[rand.Intn(len(hexChars))]
		}
	} else if kindOfString == "alpha" {
		for i := range b {
			b[i] = alpha[rand.Intn(len(alpha))]
		}
	} else {
		fmt.Println("Error retrieving character list for random string generation")
		return ""
	}
	return string(b)
}

func generateRegisterNameSig(ownerName string, data string, privateKeyStr string) (string, error) {
	privKeyByte, err := hex.DecodeString(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("error decoding hex encoded private key (%s)", privateKeyStr)
	}
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

	message := ownerName + ":" + data
	fmt.Println("Signed Message")
	fmt.Println(message)
	messageBytes := sha256.Sum256([]byte(message))
	sig, err := privKey.Sign(messageBytes[:])
	if err != nil {
		return "", fmt.Errorf("error signing message (%s) with private key (%s)", message, privateKeyStr)
	}
	return hex.EncodeToString(sig.Serialize()), nil
}

func generateRegisterThingSig(ownerName string, identities []string, spec string, data string, privateKeyStr string) (string, error) {
	privKeyByte, err := hex.DecodeString(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("error decoding hex encoded private key (%s)", privateKeyStr)
	}
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

	message := ownerName
	for _, identity := range identities {
		message += ":" + identity
	}
	message += ":" + data
	message += ":" + spec
	fmt.Println("Signed Message")
	fmt.Println(message)
	messageBytes := sha256.Sum256([]byte(message))
	sig, err := privKey.Sign(messageBytes[:])
	if err != nil {
		return "", fmt.Errorf("error signing message (%s) with private key (%s)", message, privateKeyStr)
	}
	return hex.EncodeToString(sig.Serialize()), nil
}

func generateRegisterSpecSig(specName string, ownerName string, data string, privateKeyStr string) (string, error) {
	privKeyByte, err := hex.DecodeString(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("error decoding hex encoded private key (%s)", privateKeyStr)
	}
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

	message := specName + ":" + ownerName + ":" + data
	fmt.Println("Signed Message")
	fmt.Println(message)
	messageBytes := sha256.Sum256([]byte(message))
	sig, err := privKey.Sign(messageBytes[:])
	if err != nil {
		return "", fmt.Errorf("error signing message (%s) with private key (%s)", message, privateKeyStr)
	}
	return hex.EncodeToString(sig.Serialize()), nil
}

func checkQuery(t *testing.T, stub *shim.MockStub, function string, index string, value string) {
	var err error = nil
	var bytes []byte

	bytes, err = stub.MockQuery(function, []string{index})
	if err != nil {
		fmt.Println("Query", index, "failed", err)
		t.FailNow()
	}
	if bytes == nil {
		fmt.Println("Query", index, "failed to get value")
		t.FailNow()
	}
	fmt.Printf("\nreturned from query: %s\n\n", bytes)
	if string(bytes) != value {
		fmt.Printf("json string (%s) returned from (%s) function query. Want (%s)", string(bytes), function, value)
		t.FailNow()
	}
}

func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInit("1", "", args)
	if err != nil {
		fmt.Println("INIT", args, "failed", err)
		t.FailNow()
	}
}

//register an owner to ledger
func registerOwner(t *testing.T, stub *shim.MockStub, name string, data string,
	privateKeyString string, pubKeyString string) {

	registerName := IOTRegistryTX.RegisterIdentityTX{}
	registerName.OwnerName = name
	pubKeyBytes, err := hex.DecodeString(pubKeyString)
	if err != nil {
		fmt.Println(err)
	}
	registerName.PubKey = pubKeyBytes
	registerName.Data = data

	//create signature
	hexOwnerSig, err := generateRegisterNameSig(registerName.OwnerName, registerName.Data, privateKeyString)
	if err != nil {
		fmt.Println(err)
	}
	registerName.Signature, err = hex.DecodeString(hexOwnerSig)
	if err != nil {
		fmt.Println(err)
	}
	registerNameBytes, err := proto.Marshal(&registerName)
	registerNameBytesStr := hex.EncodeToString(registerNameBytes)
	_, err = stub.MockInvoke("3", "registerOwner", []string{registerNameBytesStr})
	if err != nil {
		fmt.Println(err)
	}
}

//registers a thing to ledger
func registerThing(t *testing.T, stub *shim.MockStub, nonce []byte, identities []string,
	name string, spec string, data string, privateKeyString string) {

	registerThing := IOTRegistryTX.RegisterThingTX{}

	registerThing.Nonce = nonce
	registerThing.Identities = identities
	registerThing.OwnerName = name
	registerThing.Spec = spec

	//create signature
	hexThingSig, err := generateRegisterThingSig(name, identities, spec, data, privateKeyString)
	if err != nil {
		fmt.Println(err)
	}
	registerThing.Signature, err = hex.DecodeString(hexThingSig)
	if err != nil {
		fmt.Println(err)
	}

	registerThing.Data = data
	registerThingBytes, err := proto.Marshal(&registerThing)
	registerThingBytesStr := hex.EncodeToString(registerThingBytes)
	_, err = stub.MockInvoke("3", "registerThing", []string{registerThingBytesStr})
	if err != nil {
		fmt.Println(err)
	}
}

//registers a spec to ledger
func registerSpec(t *testing.T, stub *shim.MockStub, specName string, ownerName string,
	data string, privateKeyString string) {

	registerSpec := IOTRegistryTX.RegisterSpecTX{}

	registerSpec.SpecName = specName
	registerSpec.OwnerName = ownerName
	registerSpec.Data = data

	//create signature
	hexSpecSig, err := generateRegisterSpecSig(specName, ownerName, data, privateKeyString)
	if err != nil {
		fmt.Println(err)
	}
	registerSpec.Signature, err = hex.DecodeString(hexSpecSig)
	if err != nil {
		fmt.Println(err)
	}

	registerSpecBytes, err := proto.Marshal(&registerSpec)
	registerSpecBytesStr := hex.EncodeToString(registerSpecBytes)
	_, err = stub.MockInvoke("3", "registerSpec", []string{registerSpecBytesStr})
	if err != nil {
		fmt.Println(err)
	}
}

func newPrivateKeyString() (string, error) {
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", fmt.Errorf("Error generating private key\n")
	}
	privKeyBytes := privKey.Serialize()
	privKeyString := hex.EncodeToString(privKeyBytes)
	return privKeyString, nil
}

func getPubKeyString(privKeyString string) (string, error) {
	privKeyBytes, err := hex.DecodeString(privKeyString)
	if err != nil {
		return "", fmt.Errorf("error decoding private key string (%s)", privKeyString)
	}
	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
	pubKeyBytes := pubKey.SerializeCompressed()
	pubkKeyString := hex.EncodeToString(pubKeyBytes)
	return pubkKeyString, nil
}

func TestIOTRegistryChaincode(t *testing.T) {
	//testing private and public key generation
	privKeyString, err := newPrivateKeyString()
	if err != nil {
		fmt.Println(err)
	}
	pubKeyString, err := getPubKeyString(privKeyString)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("new privKey: (%s)\nnew pubKey: %s\n", privKeyString, pubKeyString)

	//declaring and initializing variables for all tests
	bst := new(IOTRegistry)
	stub := shim.NewMockStub("IOTRegistry", bst)

	//test register owner
	registerOwner(t, stub, "Alice", "Test Data",
		/*private key: */ "94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20",
		/*public key: */ "02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc")
	checkQuery(t, stub, "owner", "Alice", `{"OwnerName":"Alice","Pubkey":"AspKjH3FCQ+STN4iZK8kDXb21YpdLRXIxfWdlccL2eTc"}`)

	//test register thing
	nonceString := randString(32, "hex")
	fmt.Printf("nonce:%s\n", nonceString)
	nonceBytes, err := hex.DecodeString(nonceString)
	if err != nil {
		fmt.Printf("error decoding nonce hex string in TestIOTRegistry Chaincode: %v", err)
	}
	spec := "Test spec"
	identities := []string{"Foo", "Bar"}
	fmt.Printf("len identities: %d\n", len(identities))
	registerThing(t, stub, nonceBytes, identities, "Alice", spec, "Test Data",
		/*private key: */ "94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20")
	checkQuery(t, stub, "thing", nonceString, `{"Alias":["Foo","Bar"],"OwnerName":"Alice","Data":"Test Data","Spec":"Test spec"}`)

	// test register spec
	registerSpec(t, stub, "Test spec 1", "Alice", "Test data",
		/*private key: */ "94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20")
	checkQuery(t, stub, "spec", "Test spec 1", `{"OwnerName":"Alice","Data":"Test data"}`)
}
