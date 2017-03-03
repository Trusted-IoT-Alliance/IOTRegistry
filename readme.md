## IOTRegistry

This is a hyperledger chaincode which provides functionality for registering IOT devices, owners, and specifications to a blockchain. This chaincode can be used to invoke transactions on and query the ledger, and includes functionality to generate signatures for each kind of transaction. The three transaction types are registerOwner, registerThing, and registerSpec.  

## Implementation

The primary file for this chaincode is IOTRegistry.go. Invoke() and Query() are the central methods of the chaincode.  

### Invoke
Invoke is used to create a transaction. The three kinds of transactions are "registerOwner", "registerThing", and "registerSpec".  
  
Invoke receives a collection of arguments marshalled into a protobuffer, which are formatted according to the kind of transaction to be performed. The input struct for each transaction is defined in IOTRegistryTX/IOTRegistry.pb.go.  
  
For each kind of transaction, Invoke() does the following:  
1. Unmarshals the protobuffer into the appropriate structure  
2. Performs various checks to ensure that the input and attempted transaction are valid.  
3. Creates a store struct with the data to be put on the ledger, and  
4. Marshals that data into a protobuffer and puts the state to the ledger.  
  
Register thing is moderately more complicated than register owner or register spec.  For register thing, two kinds of states are put to the ledger. One kind is an "alias" state, which  is put for each member of the alias string slice input to an invocation of a register thing transaction. The other kind is a "thing" state, one of which kind is put to the ledger for each valid call to invoke  a register thing transaction. This allows for an owner of a thing to have multiple aliases.  

### Query
Query retrieves a state from the ledger and returns the data in JSON format.  
  
### Signature Generation
The three signature generation functions are in IOTRegistery_test.go:  
generateRegisterNameSig, generateRegisterThingSig, and generateRegisterSpecSig.  
  
## Usage

It is worth noting that the import paths may require some setup. Essentially, the folder paths to the cloned repository should be the same as the paths in the import statements at the top of IOTRegistry.go and IOTRegistry_test.go.  
  
IOTRegistry_test.go is a good place to look in order to understand how interaction with this chaincode can occur.  
In particular, TestIOTRegistryChaincode(), located around line 375 of IOTRegistry_test.go, displays exactly how to interact with the chaincode.  
  
First, bst takes the value of a new IOTRegistry type. Stub is declared, which is the primary means of interfacing with the ledger. Then, for each struct of type registryTest, a full test is run which includes registering an owner, a thing, and a spec, and performing a query for each transaction to validate the output.  
  
## Authors

Zaki Manian and Robert Nowell

 
