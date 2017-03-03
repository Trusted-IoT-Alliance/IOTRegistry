## IOTRegistry

This is a hyperledger chaincode which provides functionality for registering IOT devices,   
owners, and specifications to a blockchain. This chaincode can be used to invoke transactions  
on and query the ledger, and includes functionality to generate signatures for each kind of  
transaction. The three transaction types are registerOwner, registerThing, and registerSpec.  

## Implementation

The chaincode is located in IOTRegistry.go. Invoke() and Query() are the central methods of the chaincode. 

### Invoke
Invoke is used to create a transaction.  
The three kinds of transactions are "registerOwner", "registerThing", and "registerSpec".  
  
Invoke receives a collection of arguments marshalled into a protobuffer,  
which are formatted according to the kind of transaction to be performed.  
The arguments struct for each transaction is defined in IOTRegistryTX/IOTRegistry.pb.go.  


The three signature generation functions are in IOTRegistery_test.go. There are a few structs that are relied upon to perform chaincode operations.  

1. Transaction Structs
These structs are located under IOTRegistryTX and 

//the dual error formatting is useful for swagger debugging.

## Installation and Usage

IOTRegistry_test.go is a good place to look in order to understand how interaction with this chaincode can occur.
It is worth noting that the import paths may require some setup. Essentially, the folder paths to the cloned repository should be the same as the paths in the import statements at the top of IOTRegistry.go and IOTRegistry_test.go.  


## Authors

Zaki Manian and Robert Nowell

 