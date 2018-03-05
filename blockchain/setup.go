package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	chmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/chmgmtclient"
	resmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/resmgmtclient"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"time"
)

// FabricSetup implementation
type FabricSetup struct {
	ConfigFile      string
	OrgID           string
	ChannelID       string
	ChainCodeID     string
	initialized     bool
	ChannelConfig   string
	ChaincodeGoPath string
	ChaincodePath   string
	OrgAdmin        string
	OrgName         string
	client          chclient.ChannelClient
	admin           resmgmt.ResourceMgmtClient
	sdk             *fabsdk.FabricSDK
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func (setup *FabricSetup) Initialize() error {

	// Add parameters for the initialization
	if setup.initialized {
		return fmt.Errorf("sdk already initialized")
	}

	//TODO
	err := fmt.Errorf("")

	setup.sdk, err = fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return fmt.Errorf("failed to create sdk: %v", err)
	}

	// Channel management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	chMgmtClient, err := setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName)).ChannelMgmt()
	if err != nil {
		return fmt.Errorf("failed to add Admin user to sdk: %v", err)
	}

	// Org admin user is signing user for creating channel
	session, err := setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName)).Session()
	if err != nil {
		return fmt.Errorf("failed to get session for %s, %s: %s", setup.OrgName, setup.OrgAdmin, err)
	}
	orgAdminUser := session

	// Create channel
	req := chmgmt.SaveChannelRequest{ChannelID: setup.ChannelID, ChannelConfig: setup.ChannelConfig + "chainhero.channel.tx", SigningIdentity: orgAdminUser}
	if err = chMgmtClient.SaveChannel(req); err != nil {
		return fmt.Errorf("failed to create channel: %v", err)
	}

	// Allow orderer to process channel creation
	time.Sleep(time.Second * 5)

	// Org resource management client
	setup.admin, err = setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin)).ResourceMgmt()
	if err != nil {
		return fmt.Errorf("failed to create new resource management client: %v", err)
	}

	// Org peers join channel
	if err = setup.admin.JoinChannel(setup.ChannelID); err != nil {
		return fmt.Errorf("org peers failed to join the channel: %v", err)
	}

	fmt.Println("Initialization Successful")
	return nil
}

func (setup *FabricSetup) InstallAndInstantiateCC() error {

	// Create chaincode package for our chaincode
	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return fmt.Errorf("failed to create chaincode package: %v", err)
	}

	// Install our chaincode on org peers
	installCCReq := resmgmt.InstallCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: "1.0", Package: ccPkg}
	_, err = setup.admin.InstallCC(installCCReq)
	if err != nil {
		return fmt.Errorf("failed to install cc to org peers %v", err)
	}

	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.hf.chainhero.io"})

	// Org resource manager will instantiate our chaincode on the channel
	err = setup.admin.InstantiateCC(setup.ChannelID, resmgmt.InstantiateCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: "1.0", Args: [][]byte{[]byte("init")}, Policy: ccPolicy})
	if err != nil {
		return fmt.Errorf("failed to instantiate the chaincode: %v", err)
	}

	// Channel client is used to query and execute transactions
	setup.client, err = setup.sdk.NewClient(fabsdk.WithUser("User1")).Channel(setup.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to create new channel client: %v", err)
	}

	// Clean the client
	defer setup.client.Close()

	fmt.Println("Chaincode Installation & Instantiation Successful")
	return nil
}
