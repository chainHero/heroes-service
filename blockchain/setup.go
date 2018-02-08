package blockchain

import (
	"fmt"
	"os"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	"time"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/ccpackager/gopackager"
	resmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/resmgmtclient"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	chmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/chmgmtclient"
)

// FabricSetup implementation
type FabricSetup struct {
	ConfigFile      string
	OrgID           string
	ChannelID       string
	ChainCodeID     string
	Initialized     bool
	ChannelConfig 	string
	ChaincodeGoPath	string
	ChaincodePath	string
	OrgAdmin		string
	OrgName			string
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func Initialize() (*FabricSetup, error) {

	// Add parameters for the initialization
	setup := FabricSetup{
		// Channel parameters
		ChannelID:        	"chainhero",
		ChannelConfig:    	"/home/maaub/go/src/github.com/chainHero/heroes-service/fixtures/artifacts/",

		// Chaincode parameters
		ChainCodeID:      	"heroes-service",
		ChaincodeGoPath:  	os.Getenv("GOPATH"),
		ChaincodePath:    	"github.com/chainHero/heroes-service/chaincode/",
		OrgAdmin:			"Admin",
		OrgName:			"Org1",
		ConfigFile:			"config.yaml",
	}

	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk: %v", err)
	}

	// Channel management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	chMgmtClient, err := sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName)).ChannelMgmt()
	if err != nil {
		return nil, fmt.Errorf("failed to add Admin user to sdk: %v", err)
	}

	// Org admin user is signing user for creating channel
	session, err := sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName)).Session()
	if err != nil {
		return nil, fmt.Errorf("failed to get session for %s, %s: %s", setup.OrgName, setup.OrgAdmin, err)
	}
	orgAdminUser := session.Identity()

	// Create channel
	req := chmgmt.SaveChannelRequest{ChannelID: setup.ChannelID, ChannelConfig: setup.ChannelConfig + "chainhero.channel.tx", SigningIdentity: orgAdminUser}
	if err = chMgmtClient.SaveChannel(req); err != nil {
		return nil, fmt.Errorf("failed to create channel: %v", err)
	}

	// Allow orderer to process channel creation
	time.Sleep(time.Second * 5)

	// Org resource management client
	orgResMgmt, err := sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin)).ResourceMgmt()
	if err != nil {
		return nil, fmt.Errorf("failed to create new resource management client: %v", err)
	}

	// Org peers join channel
	if err = orgResMgmt.JoinChannel(setup.ChannelID); err != nil {
		return nil, fmt.Errorf("org peers failed to join the channel: %v", err)
	}

	// Create chaincode package for example cc
	fmt.Println("Start Create")
	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create chaincode package: %v", err)
	}

	// Install example cc to org peers
	fmt.Println("Start Install")
	installCCReq := resmgmt.InstallCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: "1.0", Package: ccPkg}
	_, err = orgResMgmt.InstallCC(installCCReq)
	if err != nil {
		return nil, fmt.Errorf("failed to install cc to org peers %v", err)
	}

	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.hf.chainhero.io"})

	// Org resource manager will instantiate 'example_cc' on channel
	fmt.Println("Start Instantiate")
	err = orgResMgmt.InstantiateCC(setup.ChannelID, resmgmt.InstantiateCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: "1.0", Args: [][]byte{[]byte("init")}, Policy: ccPolicy})
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate the chaincode: %v", err)
	}
	return &setup, nil
}