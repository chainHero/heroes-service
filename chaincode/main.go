package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// HeroesServiceChaincode implementation of Chaincode
type HeroesServiceChaincode struct {
}

// Init of the chaincode
// This function is call only one when the chaincode is instantiate.
// So the goal is to prepare the ledger to handle futures requests.
func (t *HeroesServiceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### HeroesServiceChaincode Init ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check that the request concern an init
	if function != "init" {
		return shim.Error("Unknown function call")
	}

	// Put in the ledger the key/value hello/wolrd
	err := stub.PutState("hello", []byte("world"))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return a successful message
	return shim.Success(nil)
}

// Invoke
// Every futures requests named invoke will arrive here.
func (t *HeroesServiceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### HeroesServiceChaincode Invoke ###########")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	// Check that the request concern an invoke
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	// Check that the number of argument is sufficient to possibly find the function to invoke
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient, you need to provide the function to invoke.")
	}

	// In order to manage multiple type of request, we will check the first argument.
	// Here we have one possible argument: query (every query request will read in the ledger without modification)
	if args[0] == "query" {
		return t.query(stub, args)
	}

	// The update argument will manage all update in the ledger
	if args[0] == "invoke" {
		return t.invoke(stub, args)
	}

	// If the argument given match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
}

// query
// Every readonly functions in the ledger will be here
func (t *HeroesServiceChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check that the number of argument is sufficient to possibly find the function concern by the query
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient, you need to provide the function for the query.")
	}

	// Like in the Invoke function, we manage multiple type of query request with the second argument.
	// We also have only one possible argument: hello
	if args[1] == "hello" {

		// Get the state of the value matching the key hello in the ledger
		state, err := stub.GetState("hello")
		if err != nil {
			return shim.Error("Failed to get state of hello")
		}

		// Return this value in response
		return shim.Success(state)
	}

	// If the argument given match any function, we return an error
	return shim.Error("Unknown query action, check the second argument.")
}

// invoke
// Every functions that read and write in the ledger will be here
func (t *HeroesServiceChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient, you need to provide the function for the invoke.")
	}

	if args[1] == "hello" && len(args) == 3 {

		//
		err := stub.PutState("hello", []byte(args[2]))
		if err != nil {
			return shim.Error("Failed to update state of hello")
		}

		// Return this value in response
		return shim.Success(nil)
	}

	// If the argument given match any function, we return an error
	return shim.Error("Unknown invoke action, check the second argument.")
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(HeroesServiceChaincode))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
	}
}
