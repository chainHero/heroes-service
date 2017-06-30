package blockchain

import (
	api "github.com/hyperledger/fabric-sdk-go/api"
	fsgConfig "github.com/hyperledger/fabric-sdk-go/pkg/config"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	fcutil "github.com/hyperledger/fabric-sdk-go/pkg/util"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/events"
	"fmt"
)

// FabricSetup implementation
type FabricSetup struct {
	Client           api.FabricClient
	Channel          api.Channel
	EventHub         api.EventHub
	Initialized      bool
	ChannelId        string
	ChannelConfig    string
}

// Initialize reads configuration from file and sets up client, chain and event hub
func Initialize() (*FabricSetup, error) {

	// Add parameters for the initialization
	setup := FabricSetup{
		// Channel parameters
		ChannelId:        "mychannel",
		ChannelConfig:    "fixtures/channel/mychannel.tx",
	}

	// Initialize the config
	// This will read the config.yaml, in order to tell to
	// the SDK all options and how contact a peer
	configImpl, err := fsgConfig.InitConfig("config.yaml")
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

	// Make a new instance of channel pre-configured with the info we have provided,
	// but for now we can't use this channel because we need to create and
	// make some peer join it
	channel, err := fcutil.GetChannel(setup.Client, setup.ChannelId)
	if err != nil {
		return nil, fmt.Errorf("Create channel (%s) failed: %v", setup.ChannelId, err)
	}
	setup.Channel = channel

	// Get an orderer user that will be used to validate an order of proposal
	// The authentication will be made with local certificates
	ordererUser, err := fcutil.GetPreEnrolledUser(
		client,
		"ordererOrganizations/example.com/users/Admin@example.com/keystore",
		"ordererOrganizations/example.com/users/Admin@example.com/signcerts",
		"ordererAdmin",
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to get the orderer user failed: %v", err)
	}

	// Get an organisation user (admin) that will be used to sign proposal
	// The authentication will be made with local certificates
	orgUser, err := fcutil.GetPreEnrolledUser(
		client,
		"peerOrganizations/org1.example.com/users/Admin@org1.example.com/keystore",
		"peerOrganizations/org1.example.com/users/Admin@org1.example.com/signcerts",
		"peerorg1Admin",
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to get the organisation user failed: %v", err)
	}

	// Initialize the channel "mychannel" base on the genesis block
	// locate in fixtures/channel/mychannel.tx and join the peer given
	// in the configuration file to this channel
	if err := fcutil.CreateAndJoinChannel(client, ordererUser, orgUser, channel, setup.ChannelConfig); err != nil {
		return nil, fmt.Errorf("CreateAndJoinChannel return error: %v", err)
	}

	// Give the organisation user to the client for next proposal
	client.SetUserContext(orgUser)

	// Setup Event Hub
	// This will allow use to listen for some event from the chaincode
	// and make some actions. We won't use it for now.
	eventHub, err := getEventHub(client)
	if err != nil {
		return nil, err
	}
	if err := eventHub.Connect(); err != nil {
		return nil, fmt.Errorf("Failed eventHub.Connect() [%s]", err)
	}
	setup.EventHub = eventHub

	// Tell that the initialization is done
	setup.Initialized = true

	return &setup, nil
}

// getEventHub initilizes the event hub
func getEventHub(client api.FabricClient) (api.EventHub, error) {
	eventHub, err := events.NewEventHub(client)
	if err != nil {
		return nil, fmt.Errorf("Error creating new event hub: %v", err)
	}
	foundEventHub := false
	peerConfig, err := client.GetConfig().GetPeersConfig()
	if err != nil {
		return nil, fmt.Errorf("Error reading peer config: %v", err)
	}
	for _, p := range peerConfig {
		if p.EventHost != "" && p.EventPort != 0 {
			fmt.Printf("EventHub connect to peer (%s:%d)\n", p.EventHost, p.EventPort)
			eventHub.SetPeerAddr(fmt.Sprintf("%s:%d", p.EventHost, p.EventPort),
				p.TLS.Certificate, p.TLS.ServerHostOverride)
			foundEventHub = true
			break
		}
	}

	if !foundEventHub {
		return nil, fmt.Errorf("No EventHub configuration found")
	}

	return eventHub, nil
}