package main

import (
	"github.com/tohero/heroes-service/blockchain"
	"fmt"
)

func main() {
	_, err := blockchain.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v", err)
	}
}
