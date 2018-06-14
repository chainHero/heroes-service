package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chainHero/heroes-service/blockchain"
	"github.com/chainHero/heroes-service/web"
	"github.com/chainHero/heroes-service/web/controllers"
)

var installCC = flag.Bool("install-cc", true, "prepare to install chaincode.")

func main() {
	flag.Parse()

	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Channel parameters
		ChannelID:     "chainhero",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/chainHero/heroes-service/fixtures/artifacts/chainhero.channel.tx",

		// Chaincode parameters
		ChainCodeID:     "heroes-service",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/chainHero/heroes-service/chaincode/",
		OrgAdmin:        "Admin",
		OrgName:         "Org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "User1",
	}

	if *installCC {
		// Initialization of the Fabric SDK from the previously set properties
		err := fSetup.Initialize()
		if err != nil {
			fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		}

		// Install and instantiate the chaincode
		err = fSetup.InstallAndInstantiateCC()
		if err != nil {
			fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		}
	} else {
		fSetup.SetClient()
	}

	// Launch the web application listening
	app := &controllers.Application{
		Fabric: &fSetup,
	}
	web.Serve(app)
}
