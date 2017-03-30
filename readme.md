# IOTRegistry

This is a chaincode or smart-contract which provides functionality for interacting with a hyperledger blockchain to store and retrieve information about internet of things devices on the blockchain. This chaincode allows for 
1. The creation of a user (called a registrant)
2. Registration of IOT devices
3. Registration of specifications of IOT devices.  


## Chaincode Overview

When this chaincode is deployed to a running Hyperledger instance, it can be used in order to interact with a blockchain. In essence, the chaincode governs the rules of how information related to Internet of Things devices can be stored on and retrieved from a hyperledger blockchain.  
  
The chaincode is entirely contained in IOTRegistry.go. Invoke() and Query() are the central methods of the chaincode. Invoke is used to store information on the blockchain and query is used to retrieve information from the blockchain.  

## Storing Information on the Blockchain

Cababilities to store information on the blockchain include creating a user (called a registrant), registering an IOT device, and registering a specification of a device to the blockchain.  

### Invoke

Invoke is the chaincode method used to store information on the blockchain. Invoke stores information to the blockchain by completing a certain kind of *transaction* which culminates in a call to *PutState* which puts a particular state to the blockchain. Invoke receives the following parameters:
1. A chaincode stub interface type called stub. More information on this struct can be found here: https://github.com/hyperledger/fabric/blob/master/core/chaincode/shim/chaincode.go
2. function (string)  
This argument is used in a switch to determine which kind of transaction to perform. Valid inputs include "registerOwner", register  
3. args []string  
This is a collection of arguments marshalled into a protobuffer, which are formatted according to the kind of transaction to be performed. The input struct for each transaction is defined in IOTRegistryTX/IOTRegistry.pb.go.  

### Transactions
The three kinds of transactions are "createRegistrant", "registerThing", and "registerSpec".  

For each kind of transaction, Invoke() does the following:  
1. Unmarshals the protobuffer into the appropriate structure  
2. Performs various checks to ensure that the input and attempted transaction are valid.  
3. Creates a store struct with the data to be put on the ledger, and  
4. Marshals that data into a protobuffer and puts the state to the ledger.  

#### createRegistrant
In order to register IOT devices and specifications to the blockchain, a valid user (called a registrant) must first exist. createRegistrant accomplishes this.    

<img src="https://github.com/InternetofTrustedThings/IOTRegistry/blob/master/images/createRegistrant.png" 
alt="main" border="10"/>

createRegistrant does the following:
1. Unmarshall arguments protobuf into createRegistrant struct type
	createRegistrant struct looks like this:  

<img src="https://github.com/InternetofTrustedThings/IOTRegistry/blob/master/images/createRegistrantTX.png" 
alt="main" border="10"/>
2. Check that inputs exist for name, public key, and signature.
3. Verify that the registrant to be created does not already exist
4. Create what should be the message represented by the signature input as an argument
5. Use public key and message to verify the signature
6. marshall arguments into type createRegistrantStore, which looks like this:  
<img src="https://github.com/InternetofTrustedThings/IOTRegistry/blob/master/images/createRegistrantStore.png" 
alt="main" border="10"/>  

The createRegistrantStore struct holds the information that is committed to the blockchain through a call to 
```
stub.PutState("RegistrantPubkey:"+<registrantPublicKey>, storeBytes)
```
In this way, the registrant information is stored on the blockchain as a byte slice, with the public key of the registrant serving as the lookup index.  

#### registerThing

Once a registrant has been created, IOT devices can be registered to the blockchain. For this purpose, registerThing is the relevant invoke function.

For register thing, two kinds of states are put to the ledger. One kind is an "alias" state, which  is put for each member of the alias string slice input to an invocation of a register thing transaction. The other kind is a "thing" state, one of which kind is put to the ledger for each valid call to invoke  a register thing transaction. This allows for an owner of a thing to have multiple aliases.  

registerThing stores information on the blockchain about a particular IOT device. The struct which represents the stored information on the blockchain looks like this:  
<img src="https://github.com/InternetofTrustedThings/IOTRegistry/blob/master/images/registerThingStore.png" 
alt="main" border="10"/>  
This struct has the following elements:
1. A list of aliases of the registrant. This list of aliases allows for multiple registrants to be associated with a single IOT device.
2. The public key of the registrant. This can be used to look up registrant information such as name.
3. Arbitrary data which can be used to describe the device.
4. The name of the specification which governs the device. This is a schema which defines the formatting of the data argument.

a registerThing transaction involves the following steps:
1. unmarshal the collection of arguments into a struct of type RegisterThingTX, which looks like this:  
<img src="https://github.com/InternetofTrustedThings/IOTRegistry/blob/master/images/registerThingTX.png" 
alt="main" border="10"/>  
2. check that the arguments contain a public key, a nonce (which can be thought of as a random string of characters that serves as an identifier for the device), and a signature.
3. Check that:  
	a. the nonce does not already exist on then ledger as an identifier for an IOT device  
	b. that the registrant exists on the ledger (has been created through a createRegistrant transaction), and  
	c. whether any of the aliases supplied already exist as registrants on the ledger.
4. Recreate the signed message and verify input signature with registrant public key
5. Next, store relevant information on the ledger:  
<img src="https://github.com/InternetofTrustedThings/IOTRegistry/blob/master/images/registerThingStates.png" 
alt="main" border="10"/>  
5a. Put to the blockchain an alias for each member of registerThingArgs.Aliases (alternate public keys connected to the device)  
5b. Put to the blockchain a thing with the information contained in the registerThingStoreType.


#### registerSpec

An IOT device can have a specification that provides information about the device. In particular, the specification defines the schema which governs the data field of the struct representing the device. 

<img src="https://github.com/InternetofTrustedThings/IOTRegistry/blob/master/images/registerSpec.png" 
alt="main" border="10"/>  

Registering a spec involves the following steps:
1. unmarshall arguments into s registerSpecTX struct, which looks like this:
<img src="https://github.com/InternetofTrustedThings/IOTRegistry/blob/master/images/registerSpecTX.png" 
alt="main" border="10"/>  
2. Verify that necessary arguments were input, that the spec does not already exist, and that the signature is valid.
3. Marshal arguments into protobuf of type Spec, which looks like this:
<img src="https://github.com/InternetofTrustedThings/IOTRegistry/blob/master/images/registerSpecStore.png" 
alt="main" border="10"/>  

  
### Query
Query retrieves a state from the ledger and returns data in JSON.  
  
### Signature Generation
The three signature generation functions are in IOTRegistery_test.go:  
generateRegisterNameSig, generateRegisterThingSig, and generateRegisterSpecSig.  
  
## Testing
  
IOTRegistry_test.go is a good place to look in order to understand how interaction with this chaincode can occur.   
  
First, bst takes the value of a new IOTRegistry type. Stub is declared, which is the primary means of interfacing with the ledger. Then, for each struct of type registryTest, a full test is run which includes registering an owner, a thing, and a spec, and performing a query for each transaction to validate the output.  
  

## Authors

Zaki Manian and Robert Nowell

 


 
