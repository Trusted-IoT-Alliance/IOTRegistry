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
	"github.com/xeipuuv/gojsonschema"

	IOTRegistryTX "github.com/InternetofTrustedThings/IOTRegistry/IOTRegistryTX"
	"github.com/btcsuite/btcd/btcec"
)

/*
Notes from IOTRegistery tests

Private Key 1: 94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20
 Public Key 1: 02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc

Private Key 2: 246d4fa59f0baa3329d3908659936ac2ac9c3539dc925977759cffe3c6316e19
 Public Key 2: 03442b817ad2154766a8f5192fc5a7506b7e52cdbf4fcf8e1bc33764698443c3c9

Private Key 3: 166cc93d9eadb573b329b5993b9671f1521679cea90fe52e398e66c1d6373abf
 Public Key 3: 02242a1c19bc831cd95a9e5492015043250cbc17d0eceb82612ce08736b8d753a6

Private Key 4: 01b756f231c72747e024ceee41703d9a7e3ab3e68d9b73d264a0196bd90acedf
 Public Key 4: 020f2b95263c4b3be740b7b3fda4c2f4113621c1a7a360713a2540eeb808519cd6

Unused:

Public Key: 02cb6d65b04c4b84502015f918fe549e95cad4f3b899359a170d4d7d438363c0ce
Private Key: 60977f22a920c9aa18d58d12cb5e90594152d7aa724bcce21484dfd0f4490b58
Hyperledger address hex 10734390011641497f489cb475743b8e50d429bb
Hyperledger address Base58: EHxhLN3Ft4p9jPkR31MJMEMee9G

Owner1 key
Public Key: 0278b76afbefb1e1185bc63ed1a17dd88634e0587491f03e9a8d2d25d9ab289ee7
Private Key: 7142c92e6eba38de08980eeb55b8c98bb19f8d417795adb56b6c4d25da6b26c5

Owner2 key
Public Key: 02e138b25db2e74c54f8ca1a5cf79e2d1ed6af5bd1904646e7dc08b6d7b0d12bfd
Private Key: b18b7d3082b3ff9438a7bf9f5f019f8a52fb64647ea879548b3ca7b551eefd65
*/

var hexChars = []rune("0123456789abcdef")
var alpha = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

/*
	testing tool for creating randomized string with a certain character makeup
*/
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
		fmt.Println("randString() error: could not retrieve character list for random string generation")
		return ""
	}
	return string(b)
}

/*
	generates a signature for registering an OwnerName based on private key and message
*/
func generateRegisterOwnerSig(ownerName string, data string, privateKeyStr string) (string, error) {
	privKeyByte, err := hex.DecodeString(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("error decoding hex encoded private key (%s)", privateKeyStr)
	}
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

	message := ownerName + ":" + data
	messageBytes := sha256.Sum256([]byte(message))
	sig, err := privKey.Sign(messageBytes[:])
	if err != nil {
		return "", fmt.Errorf("error signing message (%s) with private key (%s)", message, privateKeyStr)
	}
	return hex.EncodeToString(sig.Serialize()), nil
}

/*
	generates a signature for registering a thing based on private key and message
*/
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
	messageBytes := sha256.Sum256([]byte(message))
	sig, err := privKey.Sign(messageBytes[:])
	if err != nil {
		return "", fmt.Errorf("error signing message (%s) with private key (%s)", message, privateKeyStr)
	}
	return hex.EncodeToString(sig.Serialize()), nil
}

/*
	generates a signature for registering a spec based on private key and message
*/
func generateRegisterSpecSig(specName string, ownerName string, data string, privateKeyStr string) (string, error) {
	privKeyByte, err := hex.DecodeString(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("error decoding hex encoded private key (%s)", privateKeyStr)
	}
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

	message := specName + ":" + ownerName + ":" + data
	messageBytes := sha256.Sum256([]byte(message))
	sig, err := privKey.Sign(messageBytes[:])
	if err != nil {
		return "", fmt.Errorf("error signing message (%s) with private key (%s)", message, privateKeyStr)
	}

	return hex.EncodeToString(sig.Serialize()), nil
}

/*
	checks if Init() works properly
*/
func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInit("1", "", args)
	if err != nil {
		fmt.Println("INIT", args, "failed", err)
		t.FailNow()
	}
}

/*
	register a store type "Identites" to ledger by calling to Invoke()
*/
func registerOwner(t *testing.T, stub *shim.MockStub, name string, data string,
	privateKeyString string, pubKeyString string) error {

	registerName := IOTRegistryTX.RegisterOwnerTX{}
	registerName.OwnerName = name
	pubKeyBytes, err := hex.DecodeString(pubKeyString)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	registerName.PubKey = pubKeyBytes
	registerName.Data = data

	//create signature
	hexOwnerSig, err := generateRegisterOwnerSig(registerName.OwnerName, registerName.Data, privateKeyString)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	registerName.Signature, err = hex.DecodeString(hexOwnerSig)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	registerNameBytes, err := proto.Marshal(&registerName)
	registerNameBytesStr := hex.EncodeToString(registerNameBytes)
	_, err = stub.MockInvoke("3", "registerOwner", []string{registerNameBytesStr})
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

/*
	registers a store type "Things" to ledger and an "Alias" store type for each member of string slice identities by calling to Invoke()
*/
func registerThing(t *testing.T, stub *shim.MockStub, nonce []byte, identities []string,
	name string, spec string, data string, privateKeyString string) error {

	registerThing := IOTRegistryTX.RegisterThingTX{}

	registerThing.Nonce = nonce
	registerThing.Identities = identities
	registerThing.OwnerName = name
	registerThing.Spec = spec

	//create signature
	hexThingSig, err := generateRegisterThingSig(name, identities, spec, data, privateKeyString)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	registerThing.Signature, err = hex.DecodeString(hexThingSig)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	registerThing.Data = data
	registerThingBytes, err := proto.Marshal(&registerThing)
	registerThingBytesStr := hex.EncodeToString(registerThingBytes)
	_, err = stub.MockInvoke("3", "registerThing", []string{registerThingBytesStr})
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

/*
	registers a store type "Spec" to ledger by calling to Invoke()
*/
func registerSpec(t *testing.T, stub *shim.MockStub, specName string, ownerName string,
	data string, privateKeyString string) error {

	registerSpec := IOTRegistryTX.RegisterSpecTX{}

	registerSpec.SpecName = specName
	registerSpec.OwnerName = ownerName
	registerSpec.Data = data

	//create signature
	hexSpecSig, err := generateRegisterSpecSig(specName, ownerName, data, privateKeyString)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	registerSpec.Signature, err = hex.DecodeString(hexSpecSig)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	registerSpecBytes, err := proto.Marshal(&registerSpec)
	registerSpecBytesStr := hex.EncodeToString(registerSpecBytes)
	_, err = stub.MockInvoke("3", "registerSpec", []string{registerSpecBytesStr})
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

/*
	//To create new private and public keys
	privKeyString, err := newPrivateKeyString()
	if err != nil {
		fmt.Println(err)
	}
	pubKeyString, err := getPubKeyString(privKeyString)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("new privKey: (%s)\nnew pubKey: %s\n", privKeyString, pubKeyString)
*/

/*
	generates and returns SHA256 private key string
*/
func newPrivateKeyString() (string, error) {
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", fmt.Errorf("Error generating private key\n")
	}
	privKeyBytes := privKey.Serialize()
	privKeyString := hex.EncodeToString(privKeyBytes)
	return privKeyString, nil
}

/*
	generates and returns SHA256 public key string from private key string input
*/
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

/*
	tests whether two string slices are identical, returning true or false
*/
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

/*
	Checks that different queries return expected values.
*/
func checkQuery(t *testing.T, stub *shim.MockStub, function string, index string, expected registryTest) error {
	var err error = nil
	var bytes []byte

	bytes, err = stub.MockQuery(function, []string{index})
	if err != nil {
		return fmt.Errorf("Query (%s) failed\n", function)
	}
	if bytes == nil {
		return fmt.Errorf("Query (%s) failed to get value\n", function)
	}
	fmt.Printf("\n\nBYTES FROM QUERY: (%s)\n\n", bytes)
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(bytes, &jsonMap); err != nil {
		return fmt.Errorf("error unmarshalling json string %s", bytes)
	}
	fmt.Printf("JSON: %s\n", jsonMap)
	if function == "owner" {
		if jsonMap["OwnerName"] != expected.ownerName {
			return fmt.Errorf("\nOwnerName got       (%s)\nOwnerName expected: (%s)\n", jsonMap["OwnerName"], expected.ownerName)
		}
		if jsonMap["Pubkey"] != expected.pubKeyString {
			return fmt.Errorf("\nPubkey got       (%s)\nPubkey expected: (%s)\n", jsonMap["Pubkey"], expected.pubKeyString)
		}
	} else if function == "thing" {
		aliases := make([]string, len(jsonMap["Alias"].([]interface{})))
		for i, element := range jsonMap["Alias"].([]interface{}) {
			aliases[i] = element.(string)
		}
		if !(reflect.DeepEqual(aliases, expected.identities)) {
			return fmt.Errorf("\nAlias got       (%x)\nAlias expected: (%x)\n", jsonMap["Alias"], expected.identities)
		}
		if jsonMap["OwnerName"] != expected.ownerName {
			return fmt.Errorf("\nOwnerName got       (%s)\nOwnerName expected: (%s)\n", jsonMap["OwnerName"], expected.ownerName)
		}
		if jsonMap["Data"] != expected.data {
			return fmt.Errorf("\nData got       (%s)\nData expected: (%s)\n", jsonMap["Data"], expected.data)
		}
		if jsonMap["SpecName"] != expected.specName {
			return fmt.Errorf("\nSpecName got       (%s)\nSpecName expected: (%s)\n", jsonMap["SpecName"], expected.specName)
		}
	} else if function == "spec" {
		if jsonMap["OwnerName"] != expected.ownerName {
			return fmt.Errorf("\nOwnerName got       (%s)\nOwnerName expected: (%s)\n", jsonMap["OwnerName"], expected.ownerName)
		}
		if jsonMap["Data"] != expected.specSchema {
			return fmt.Errorf("\nData got       (%s)\nData expected: (%s)\n", jsonMap["Data"], expected.data)
		}
	}
	return nil
}

func validateSchema(t *testing.T, stub *shim.MockStub, info registryTest) error {
	var err error = nil

	thingBytes, err := stub.MockQuery("thing", []string{hex.EncodeToString(info.nonceBytes)})
	if err != nil {
		return fmt.Errorf("Thing query on nonce (%s) failed\n", hex.EncodeToString(info.nonceBytes))
	}
	if thingBytes == nil {
		return fmt.Errorf("no value returned from thing query on nonce (%s)\n", hex.EncodeToString(info.nonceBytes))
	}
	fmt.Printf("\n\nthingBytes from query: (%s)\n\n", thingBytes)

	var thingJSON map[string]interface{}
	if err := json.Unmarshal(thingBytes, &thingJSON); err != nil {
		return fmt.Errorf("error unmarshalling json string (%s)", thingBytes)
	}
	fmt.Printf("thingJSON: %s\n", thingJSON)

	specName := thingJSON["SpecName"].(string)
	specBytes, err := stub.MockQuery("spec", []string{specName})
	if err != nil {
		return fmt.Errorf("Query on specName (%s) failed\n", specName)
	}
	if specBytes == nil {
		return fmt.Errorf("no value returned from query on specName (%s)\n", specName)
	}
	fmt.Printf("\n\nspecBytes from query: (%s)\n\n", specBytes)
	var specJSON map[string]interface{}
	if err := json.Unmarshal(specBytes, &specJSON); err != nil {
		return fmt.Errorf("error unmarshalling json string %s", specBytes)
	}
	fmt.Printf("specJSON: %s\n", specJSON)

	fmt.Printf("doc: (%s)\n\n", thingJSON["Data"].(string))
	documentLoader := gojsonschema.NewStringLoader(thingJSON["Data"].(string))
	schemaLoader := gojsonschema.NewStringLoader(specJSON["Data"].(string))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {

			fmt.Printf("- %s\n", desc)
		}
		t.Errorf("error validating schema\n")
		return err
	}
	return nil
}

type registryTest struct {
	privateKeyString string
	pubKeyString     string
	ownerName        string
	data             string
	nonceBytes       []byte
	specName         string
	identities       []string
	specSchema       string
}

var testData = `{
	"id": 2,
	"name": "An ice sculpture",
	"price": 12.50,
	"tags": ["cold", "ice"],
	"dimensions": {
		"length": 7.0,
		"width": 12.0,
		"height": 9.5
	}
}`

var testSchema = `{
	"$schema": "http://json-schema.org/draft-04/schema#",
	"title": "Product",
	"description": "A product from catalog",
	"type": "object",
	"properties": {
		"id": {
			"description": "The unique identifier for a product",
			"type": "integer"
		},
		"name": {
			"description": "Name of the product",
			"type": "string"
		},
		"price": {
			"type": "number",
			"minimum": 0,
			"exclusiveMinimum": true
		},
		"tags": {
			"type": "array",
			"items": {
				"type": "string"
			},
			"minItems": 1,
			"uniqueItems": true
		},
		"dimensions": {
			"type": "object",
			"properties": {
				"length": {"type": "number"},
				"width": {"type": "number"},
				"height": {"type": "number"}
			},
			"required": ["length", "width", "height"]
		}
	},
	"required": ["id", "name", "price"]
}`

/*
	runs tests for four different users: Alice, Gerald, Bob, and Cassandra
*/
func TestIOTRegistryChaincode(t *testing.T) {
	//declaring and initializing variables for all tests
	bst := new(IOTRegistry)
	stub := shim.NewMockStub("IOTRegistry", bst)

	nonceBytes1, err := hex.DecodeString("1f7b169c846f218ab552fa82fbf86758")
	if err != nil {
		t.Errorf("error decoding nonce hex string in TestIOTRegistryChaincode: %v", err)
	}
	nonceBytes2, err := hex.DecodeString("bf5c97d2d2a313e4f95957818a7b3edc")
	if err != nil {
		t.Errorf("error decoding nonce hex string in TestIOTRegistryChaincode: %v", err)
	}
	nonceBytes3, err := hex.DecodeString("a492f2b8a67697c4f91d9b9332e82347")
	if err != nil {
		t.Errorf("error decoding nonce hex string in TestIOTRegistryChaincode: %v", err)
	}
	nonceBytes4, err := hex.DecodeString("83de17bd7a25e0a9f6813976eadf26de")
	if err != nil {
		t.Errorf("error decoding nonce hex string in TestIOTRegistryChaincode: %v", err)
	}
	var registryTestsSuccess = []registryTest{
		{ /*private key  1*/ "94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20",
			/*public key 1*/ "02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc",
			"Alice", testData, nonceBytes1, "test spec 1", []string{"Foo", "Bar"}, testSchema},

		{ /*private key  2*/ "246d4fa59f0baa3329d3908659936ac2ac9c3539dc925977759cffe3c6316e19",
			/*public key 2*/ "03442b817ad2154766a8f5192fc5a7506b7e52cdbf4fcf8e1bc33764698443c3c9",
			"Gerald", testData, nonceBytes2, "test spec 2", []string{"one", "two", "three"}, testSchema},

		{ /*private key  3*/ "166cc93d9eadb573b329b5993b9671f1521679cea90fe52e398e66c1d6373abf",
			/*public key 3*/ "02242a1c19bc831cd95a9e5492015043250cbc17d0eceb82612ce08736b8d753a6",
			"Bob", testData, nonceBytes3, "test spec 3", []string{"ident4", "ident5", "ident6"}, testSchema},

		{ /*private key  4*/ "01b756f231c72747e024ceee41703d9a7e3ab3e68d9b73d264a0196bd90acedf",
			/*public key 4*/ "020f2b95263c4b3be740b7b3fda4c2f4113621c1a7a360713a2540eeb808519cd6",
			"Cassandra", testData, nonceBytes4, "test spec 4", []string{"ident7", "ident8", "ident9"}, testSchema},
	}
	for _, test := range registryTestsSuccess {
		err := registerOwner(t, stub, test.ownerName, test.data, test.privateKeyString, test.pubKeyString)
		if err != nil {
			t.Errorf("%v\n", err)
			return
		}
		index := test.ownerName
		err = checkQuery(t, stub, "owner", index, test)
		if err != nil {
			t.Errorf("%v\n", err)
		}
		err = registerThing(t, stub, test.nonceBytes, test.identities, test.ownerName, test.specName, test.data, test.privateKeyString)
		if err != nil {
			t.Errorf("%v\n", err)
		}

		index = hex.EncodeToString(test.nonceBytes)
		err = checkQuery(t, stub, "thing", index, test)
		if err != nil {
			t.Errorf("%v\n", err)
		}
		err = registerSpec(t, stub, test.specName, test.ownerName, testSchema, test.privateKeyString)
		if err != nil {
			t.Errorf("%v\n", err)
		}
		index = test.specName
		err = checkQuery(t, stub, "spec", index, test)
		if err != nil {
			t.Errorf("%v\n", err)
		}
		err = validateSchema(t, stub, test)
		if err != nil {
			t.Errorf("%v\n", err)
		}

	}
}
