package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
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

func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInit("1", "", args)
	if err != nil {
		fmt.Println("INIT", args, "failed", err)
		t.FailNow()
	}
}

//register a store type "Identites" to ledger
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

//registers a store type "Things" to ledger and an "Alias" store type for each member of string slice identities
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

//registers a store type "Spec" to ledger
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

//generates and returns SHA256 private key string
func newPrivateKeyString() (string, error) {
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", fmt.Errorf("Error generating private key\n")
	}
	privKeyBytes := privKey.Serialize()
	privKeyString := hex.EncodeToString(privKeyBytes)
	return privKeyString, nil
}

//generates and returns SHA256 public key string from private key string input
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

func testEq(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func checkQuery(t *testing.T, stub *shim.MockStub, function string, index string, expected registryTest) {
	var err error = nil
	var bytes []byte

	bytes, err = stub.MockQuery(function, []string{index})
	if err != nil {
		t.Errorf("Query (%s) failed\n", function)
		// t.FailNow()
	}
	if bytes == nil {
		t.Errorf("Query (%s) failed to get value\n", function)
		// t.FailNow()
	}
	// fmt.Printf("\nreturned from query: %s\n\n", bytes)

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(bytes, &jsonMap); err != nil {
		t.Errorf("error unmarshalling json string %s", bytes)
	}

	if function == "owner" {
		if jsonMap["OwnerName"] != expected.ownerName {
			t.Errorf("\nOwnerName got       (%s)\nOwnerName expected: (%s)\n", jsonMap["OwnerName"], expected.ownerName)
			// t.FailNow()
		}
		if jsonMap["Pubkey"] != expected.pubKeyString {
			t.Errorf("\nPubkey got       (%s)\nPubkey expected: (%s)\n", jsonMap["Pubkey"], expected.pubKeyString)
			fmt.Println("Here!")
		}
	} else if function == "thing" {
		// fmt.Printf("TYPE: (%v)\n", jsonMap["Alias"].(type))

		aliases := make([]string, len(jsonMap["Alias"].([]interface{})))
		for i, element := range jsonMap["Alias"].([]interface{}) {
			aliases[i] = element.(string)
		}
		if !(reflect.DeepEqual(aliases, expected.identities)) {
			t.Errorf("\nAlias got       (%x)\nAlias expected: (%x)\n", jsonMap["Alias"], expected.identities)
			// t.FailNow()
		}
		if jsonMap["OwnerName"] != expected.ownerName {
			t.Errorf("\nOwnerName got       (%s)\nOwnerName expected: (%s)\n", jsonMap["OwnerName"], expected.ownerName)
			// t.FailNow()
		}
		if jsonMap["Data"] != expected.data {
			t.Errorf("\nData got       (%s)\nData expected: (%s)\n", jsonMap["Data"], expected.data)
			// t.FailNow()
		}
		if jsonMap["SpecName"] != expected.specName {
			t.Errorf("\nSpecName got       (%s)\nSpecName expected: (%s)\n", jsonMap["SpecName"], expected.specName)
			// t.FailNow()
		}
	} else if function == "spec" {

	}
}

type registryTest struct {
	privateKeyString string
	pubKeyString     string
	ownerName        string
	data             string
	nonceBytes       []byte
	specName         string
	identities       []string
}

func TestIOTRegistryChaincode(t *testing.T) {
	//declaring and initializing variables for all tests
	bst := new(IOTRegistry)
	stub := shim.NewMockStub("IOTRegistry", bst)

	nonceBytes1, err := hex.DecodeString("1f7b169c846f218ab552fa82fbf86758")
	if err != nil {
		fmt.Printf("error decoding nonce hex string in TestIOTRegistryChaincode: %v", err)
	}
	// nonceBytes2, err := hex.DecodeString("bf5c97d2d2a313e4f95957818a7b3edc")
	// if err != nil {
	// 	fmt.Printf("error decoding nonce hex string in TestIOTRegistryChaincode: %v", err)
	// }
	// nonceBytes3, err := hex.DecodeString("a492f2b8a67697c4f91d9b9332e82347")
	// if err != nil {
	// 	fmt.Printf("error decoding nonce hex string in TestIOTRegistryChaincode: %v", err)
	// }
	// nonceBytes4, err := hex.DecodeString("83de17bd7a25e0a9f6813976eadf26de")
	// if err != nil {
	// 	fmt.Printf("error decoding nonce hex string in TestIOTRegistryChaincode: %v", err)
	// }

	// fmt.Printf("1: %s\n2: %s\n3: %s\n4: %s\n", nonceString1, nonceString2, nonceString3, nonceString4)
	var registryTests = []registryTest{
		{"94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20", "02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc", "Alice",
			"test data", nonceBytes1, "test spec", []string{"Foo", "Bar"}},
		// {"94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20", "02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc", "Alice",
		// 	"test data 1", nonceBytes2, "test spec 1", []string{"ident1", "ident2", "ident3"}},
		// {"94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20", "02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc", "Bob",
		// 	"test data 2", nonceBytes3, "test spec 2", []string{"ident4", "ident5", "ident6"}},
		// {"94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20", "02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc", "Cassandra",
		// 	"test data 3", nonceBytes4, "test spec 3", []string{"ident7", "ident8", "ident9"}},
	}
	for _, test := range registryTests {
		registerOwner(t, stub, test.ownerName, test.data, test.privateKeyString, test.pubKeyString)
		index := test.ownerName
		checkQuery(t, stub, "owner", index, test)

		registerThing(t, stub, test.nonceBytes, test.identities, test.ownerName, test.specName, test.data, test.privateKeyString)
		index = hex.EncodeToString(test.nonceBytes)
		checkQuery(t, stub, "thing", index, test)

		// `{"Alias":["Foo","Bar"],"OwnerName":"Alice","Data":"test data","SpecName":"test spec"}`
		registerSpec(t, stub, test.specName, test.ownerName, test.data, test.privateKeyString)
		index = test.specName
		checkQuery(t, stub, "spec", index, test)

		// `{"OwnerName":"Alice","Data":"test data"}`
	}

	// //testing private and public key generation
	// privKeyString, err := newPrivateKeyString()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// pubKeyString, err := getPubKeyString(privKeyString)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("new privKey: (%s)\nnew pubKey: %s\n", privKeyString, pubKeyString)
}
