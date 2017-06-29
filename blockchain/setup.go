package blockchain

import (
	api "github.com/hyperledger/fabric-sdk-go/api"
	fsgConfig "github.com/hyperledger/fabric-sdk-go/pkg/config"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	fcutil "github.com/hyperledger/fabric-sdk-go/pkg/util"
	"fmt"
)

// FabricSetup implementation
type FabricSetup struct {
	Client           api.FabricClient
	Channel          api.Channel
	EventHub         api.EventHub
	Initialized      bool
	ChannelID        string
	ChannelConfig    string
}

// Initialize reads configuration from file and sets up client, chain and event hub
func Initialize() (*FabricSetup, error) {

	setup := FabricSetup{
		ChannelID:       "mychannel",
		ChannelConfig:   "../fixtures/channel/mychannel.tx",
	}

	// Initialize the config
	// This will read the config.yaml, in order to tell to
	// the SDK all options and how contact a peer
	configImpl, err := fsgConfig.InitConfig("../config.yaml")
	if err != nil {
		return nil, fmt.Errorf("Initialize the config failed: %v", err)
	}

	// Initialize blockchain cryptographic service provider (BCCSP)
	// This tool manage certificates and keys
	err = bccspFactory.InitFactories(configImpl.GetCSPConfig())
	if err != nil {
		return nil, fmt.Errorf("Failed getting ephemeral software-based BCCSP [%s]", err)
	}

	// This will make a user access (here the admin) to interact with the network
	// To do so, this will contact the Fabric CA to check if the user has access
	// and give it to him (enrollment)
	client, err := fcutil.GetClient("admin", "adminpw", "/tmp/enroll_user", configImpl)
	if err != nil {
		return nil, fmt.Errorf("Create client failed: %v", err)
	}
	setup.Client = client

	// Get channel
	// TODO
	channel, err := fcutil.GetChannel(setup.Client, setup.ChannelID)
	if err != nil {
		return nil, fmt.Errorf("Create channel (%s) failed: %v", setup.ChannelID, err)
	}
	setup.Channel = channel

	// Get orderer user
	// TODO
	ordererUser, err := fcutil.GetPreEnrolledUser(
		client,
		"ordererOrganizations/example.com/users/Admin@example.com/keystore",
		"ordererOrganizations/example.com/users/Admin@example.com/signcerts",
		"ordererAdmin",
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to get the orderer user failed: %v", err)
	}

	// Get organisation user
	// TODO
	orgUser, err := fcutil.GetPreEnrolledUser(
		client,
		"peerOrganizations/org1.example.com/users/Admin@org1.example.com/keystore",
		"peerOrganizations/org1.example.com/users/Admin@org1.example.com/signcerts",
		"peerorg1Admin",
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to get the organisation user failed: %v", err)
	}

	// Create and join channel
	// TODO
	if err := fcutil.CreateAndJoinChannel(client, ordererUser, orgUser, channel, setup.ChannelConfig); err != nil {
		return nil, fmt.Errorf("CreateAndJoinChannel return error: %v", err)
	}

	client.SetUserContext(orgUser)

	setup.Initialized = true

	return &setup, nil
}
