package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
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
	UserName        string
	client          *channel.Client
	admin           *resmgmt.Client
	sdk             *fabsdk.FabricSDK
	event           *event.Client
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func (setup *FabricSetup) Initialize() error {

	// Add parameters for the initialization
	if setup.initialized {
		return fmt.Errorf("sdk already initialized")
	}
	///

	// Initialize the SDK with the configuration file
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return fmt.Errorf("failed to create sdk: %v", err)
	}
	setup.sdk = sdk
	defer sdk.Close()
	///

	// The resource management client is responsible for managing channels (create/update channel)
	resourceManagerClientContext := setup.sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(setup.OrgName))
	if err != nil {
		return fmt.Errorf("failed to add Admin user to sdk: %v", err)
	}
	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return fmt.Errorf("failed to create channel management client: %s", err)
	}
	///

	// Let's create the channel
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(setup.OrgName))
	if err != nil {
		return err
	}
	adminIdentity, err := mspClient.GetSigningIdentity(setup.OrgAdmin)
	if err != nil {
		return err
	}
	req := resmgmt.SaveChannelRequest{ChannelID: setup.ChannelID, ChannelConfigPath: setup.ChannelConfig, SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := resMgmtClient.SaveChannel(req)
	if err != nil || txID.TransactionID == "" {

	}
	///

	// Create org resource management client
	orgResMgmt, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return err
	}
	setup.admin = orgResMgmt
	///

	// Org peers join channel
	if err = orgResMgmt.JoinChannel(setup.ChannelID); err != nil {
		return err
	}
	///

	fmt.Println("Initialization Successful")
	setup.initialized = true
	return nil
}

func (setup *FabricSetup) InstallAndInstantiateCC() error {

	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return err
	}

	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: "1.0", Package: ccPkg}
	_, err = setup.admin.InstallCC(installCCReq)
	if err != nil {
		return err
	}

	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.hf.chainhero.io"})

	resp, err := setup.admin.InstantiateCC(setup.ChannelID, resmgmt.InstantiateCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodeGoPath, Version: "1.0", Args: [][]byte{[]byte("init")}, Policy: ccPolicy})
	if err != nil || resp.TransactionID == "" {

	}

	// Channel client is used to query and execute transactions
	clientContext := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName))

	setup.client, err = channel.New(clientContext)
	if err != nil {
		return fmt.Errorf("failed to create new channel client: %v", err)
	}

	setup.event, err = event.New(clientContext)
	if err != nil {
		return err
	}

	fmt.Println("Chaincode Installation & Instantiation Successful")
	return nil
}
