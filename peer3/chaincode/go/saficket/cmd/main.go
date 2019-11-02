package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/chaincode/go/saficket"
)

func main() {
	err := shim.Start(new(saficket.SmartContract))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
