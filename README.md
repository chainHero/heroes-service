# Tutorial Hyperledger Fabric SDK Go: How to build your first app?

This tutorial will introduce you to the Hyperledger Fabric Go SDK and allow you to build a simple application using the blockchain principle. The application will be a tool that allows people to make a request for a get some help. A hero can see all requests and accept some. When he achieved the service, the person who made the request can thank and reward the hero. The kind of things is like a smart contract between a simple user and a hero.

## 1. Prerequisites

This tutorial won’t explain in detail how Hyperledger Fabric works, I will just give some tips to understand the general behavior of the framework. If you want to get a full explanation of the tool, go to the official [documentation](http://hyperledger-fabric.readthedocs.io/en/latest/) there is a lot of work that explains to you what kind of blockchain is Hyperledger Fabric.

In the technical part, this tutorial has been made on **Ubuntu 16.04**. The Hyperledger Fabric framework is compatible with Mac OSX and Windows too, but I can’t guarantee that all the stuff can work.

We will use the **Go** language to design a first application, because the Hyperledger Fabric has been built also in Go and the Fabric SDK Go is really simple to use. In addition, the chaincode can be write in Go too, so the full-stack will be only in Go! There are other SDK if you want to, like for NodeJS, Java or Python.

Hyperledger Fabric uses **Docker** to easily deploy a blockchain network. In addition, in the v1.0, some component (peers) also deploys docker containers to separate data (channel). So make sure that the platform supports this kind of virtualization (we will install Docker in the installation part).

## 2. Introduction to Hyperledger Fabric

*TODO - Explain Fabric and Fabric CA*

## 3. Installation guide

This installation guide was made in **Ubuntu 16.04**.

### a. Docker

The required **version for docker is 1.12 or greater**, this version is already available in the package manager on Ubuntu. Just install it with this command line:

```
sudo apt install docker.io
```

In addition, we need **docker-compose 1.8+** to manage multiple containers at one. You can also use your package manager that hold the right version:

```
sudo apt install docker-compose
```

Now we need to manage the current user to avoid using root access when we will use docker. To do so, we need to add the user to the docker group:

```
sudo groupadd docker
sudo gpasswd -a ${USER} docker
sudo service docker restart
```

In order to apply this change, you need to logout/login. Then check versions (and that works)  with:

```
docker --version
docker-compose version
```

![End of the docker installation](docs/images/finish-docker-install.png)

### b. Go

Hyperledger Fabric required a **Go version 1.7.x** or more and we have only Go version 1.6.x in package manager. So this time we need to use the official installation method. You can follow instructions from golang.org (https://golang.org/dl/) or use this generics commands that will install Golang 1.8.3 and prepare your environment (generate your GOPATH and add variables):

```
wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz && \
sudo tar -C /usr/local -xzf go1.8.3.linux-amd64.tar.gz && \
rm go1.8.3.linux-amd64.tar.gz && \
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile && \
echo 'export GOPATH=$HOME/go' | tee -a $HOME/.bashrc && \
echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' | tee -a $HOME/.bashrc && \
mkdir -p $HOME/go/{src,pkg,bin}
```

To make sure that the installation works, you can logout/login (again) and run:

```
go verison
```

![End of the Go installation](docs/images/finish-go-install.png)

### c. Hyperledger Fabric & CA

Now we can install the main framework: Hyperledger Fabric. We will fix the commit level to the v1.0.0-rc1 because the Fabric SDK Go is compatible with it. All the code is available in a mirror on github, just check out (and optionally build binaries):

```
mkdir -p $GOPATH/src/github.com/hyperledger && \
cd $GOPATH/src/github.com/hyperledger && \
git clone https://github.com/hyperledger/fabric.git && \
cd fabric && \
git checkout v1.0.0-rc1
```

Same for the Hyperledger Fabric CA part:

```
cd $GOPATH/src/github.com/hyperledger && \
git clone https://github.com/hyperledger/fabric-ca.git && \
cd fabric-ca && \
git checkout v1.0.0-rc1
```

We won’t use directly the framework, but this is useful to have the framework locally in your GOPATH to compile your app.

### d. Fabric SDK Go

Finally, we will install the Hyperledger Fabric SDK Go that will allow us to easily communicate with the Fabric framework. To do so, we will use the built in function provide by golang:

```
go get -u github.com/hyperledger/fabric-sdk-go/pkg/fabric-client
go get -u github.com/hyperledger/fabric-sdk-go/pkg/fabric-ca-client
```

If you get the following error:

```
../fabric-sdk-go/vendor/github.com/miekg/pkcs11/pkcs11.go:29:18: fatal error: ltdl.h: No such file or directory
```

You need to install the package “libltdl-dev” and re-execute previous command (`go get ...`):

```
sudo apt install libltdl-dev
```

Then you can go inside the new fabric-sdk-go directory in your GOPATH and install dependencies and check out if all is ok:

```
cd $GOPATH/src/github.com/hyperledger/fabric-sdk-go && make
```

The installation can take a while (depending on your network connection), but at the end you should see Integration tests passed. During this process, a virtual network has been built and some test are made with the SDK in order to check if your system is ready. Now we can work with our first application.

![End of the Fabric SDK Go installation](docs/images/finish-fabric-sdk-go-install.png)

## 4. Make your first blockchain network

### a. Prepare environment

In order to make a blockchain network we will use docker to build virtual computers that will handle different roles. In this tutorial we will stay simple as possible. To do so, we will directly get the network use of testing in the Fabric SDK Go. Hyperledger Fabric needs a lot of certificates to ensure encryption in the end to end manner (SSL, TSL …).

Make a new directory in source folder of your GOPATH place, we will name it ‘heroes-service’:

```
mkdir -p $GOPATH/src/github.com/tohero/heroes-service && \
cd $GOPATH/src/github.com/tohero/heroes-service
```

Now, we can copy the environment of the Fabric SDK Go placed in the test folder:

```
cp -r $GOPATH/src/github.com/hyperledger/fabric-sdk-go/test/fixtures ./
```

We can clean up a little bit to make it more simple. We remove the default chaincode, we will make our later. We also remove some files used by the test script of the SDK:

```
rm -rf fixtures/{config,src,.env,latest-env.sh}
```

### b. Built a Docker compose file

In order to make it work, we have to edit the `docker-compose.yaml` file. This is the configuration file for docker-compose, it tells what containers need to be created and started and with a custom configuration for each. Take your favorite text editor and copy paste content from this repository:

```
cd $GOPATH/src/github.com/tohero/heroes-service && \
vi fixtures/docker-compose.yaml
```

see [fixtures/docker-compose.yaml](fixtures/docker-compose.yaml)

Now if we use docker-compose we will setup 2 fabric certificate authorities with 1 peer for each. Peers will have all roles: ledger, endorer and commiter. In addition, an orderer is also created with the `solo` ordering (no consensus is made).

In our example, one organisation will be the heroes and the other concern people who made request. So lets say that peer 0 of the organisation 1 is superman and peer 0 of the organisation 2 is John, a normal citizen.

### c. Test

In order to check if the network works, we will use command provide by docker-compose to start or stop all containers at the same time. Go inside the `fixtures` folder, and run:

```
cd $GOPATH/src/github.com/tohero/heroes-service/fixtures && \
docker-compose up
```

You will see a lot of logs with different colors (red isn't equal to errors).

Open a new terminal and run:

```
docker ps
```

![Docker compose up screenshot](docs/images/docker-ps.png)

You will see the two peers, the orderer and the two CA. To stop the network go back to the previous terminal, and press `Ctrl+C`, wait that all containers are stopped. You have succefully made a new network ready to use with th SDK. If you want to explore more deeper, check out the official documentation about this: [Building Your First Network](http://hyperledger-fabric.readthedocs.io/en/latest/build_network.html)

![Docker compose up screenshot](docs/images/docker-compose-up.png)

> **Tips**: when the network is down, all containers used remains accessible, but not active, to check logs for example. You can see them with `docker ps -a`. In order to clean up this containers, you need to delete them with `docker rm $(docker ps -aq)`.

> **Tips**: you can run the `docker-compose` command in background to keep the prompt. To do so, use the parameter `-d`, like this: `docker-compose up -d`. To stop containers, run in the same place where the `docker-compose.yaml` is: `docker-compose down` (this will also remove all containers).

## 5. Build a simple application using SDK

### a. Fabric SDK Go: configuration

Like we remove the config folder, we need to make a new config file. We will put everything that the Fabric SDK Go need and our custom parameters for the app. For now we will only try to make the Fabric SDK Go work with the default chaincode, that with we just put the blockchain configuration:

```
cd $GOPATH/src/github.com/tohero/heroes-service && \
vi config.yaml
```

```
client:
 peers:
  # peer0
  - host: "localhost"
    port: 7051
    eventHost: "localhost"
    eventPort: 7053
    primary: true
    tls:
      # Certificate location absolute path
      certificate: "$GOPATH/src/github.com/tohero/heroes-service/fixtures/channel/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/cacerts/org1.example.com-cert.pem"
      serverHostOverride: "peer0.org1.example.com"

 tls:
  enabled: true

 security:
  enabled: true
  hashAlgorithm: "SHA2"
  level: 256

 tcert:
  batch:
    size: 200

 orderer:
  host: "localhost"
  port: 7050
  tls:
    # Certificate location absolute path
    certificate: "$GOPATH/src/github.com/tohero/heroes-service/fixtures/channel/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/cacerts/example.com-cert.pem"
    serverHostOverride: "orderer.example.com"

 logging:
  level: info

 fabricCA:
  tlsEnabled: true
  id: "Org1MSP"
  name: "ca-org1"
  homeDir: "/tmp/"
  mspDir: "msp"
  serverURL: "https://localhost:7054"
  certfiles :
    - "$GOPATH/src/github.com/tohero/heroes-service/fixtures/tls/fabricca/ca/ca_root.pem"
  client:
   keyfile: "$GOPATH/src/github.com/tohero/heroes-service/fixtures/tls/fabricca/client/client_client1-key.pem"
   certfile: "$GOPATH/src/github.com/tohero/heroes-service/fixtures/tls/fabricca/client/client_client1.pem"

 cryptoconfig:
  path: "$GOPATH/src/github.com/tohero/heroes-service/fixtures/channel/crypto-config"
```

The full configuration file is available here: [config.yaml](config.yaml)

### b. Fabric SDK Go: initialize

We add a new folder named `blockchain` that will contain the whole interface that comunicate with the network. We will see the Fabric SDK go only in this folder.

```
mkdir $GOPATH/src/github.com/tohero/heroes-service/blockchain
```

Now add a new go file named `setup.go` :

```
vi $GOPATH/src/github.com/tohero/heroes-service/blockchain/setup.go
```

```
package blockchain

import (
	api "github.com/hyperledger/fabric-sdk-go/api"
	fsgConfig "github.com/hyperledger/fabric-sdk-go/pkg/config"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	fcutil "github.com/hyperledger/fabric-sdk-go/pkg/util"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/events"
	"fmt"
	"os"
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
```

The full file is available here: [blockchain/setup.go](blockchain/setup.go)

At this stage we only initialize a client that will comunicate to a peer, a CA and an orderer. We also make a new channel and make the peer join this channel. See the comments in the code for more information.

### c. Fabric SDK Go: test

To make sure that the client arrive to initialize all his components, we will make a simple test with the network up. In order to make this, we need to build the go code, but we haven't any main file. Let's add one:

```
cd $GOPATH/src/github.com/tohero/heroes-service && \
vi main.go
```

```
package main

import (
	"github.com/tohero/heroes-service/blockchain"
	"fmt"
	"os"
	"runtime"
	"path/filepath"
)

// Fix empty GOPATH with golang 1.8 (see https://github.com/golang/go/blob/1363eeba6589fca217e155c829b2a7c00bc32a92/src/go/build/build.go#L260-L277)
func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		def := filepath.Join(home, "go")
		if filepath.Clean(def) == filepath.Clean(runtime.GOROOT()) {
			// Don't set the default GOPATH to GOROOT,
			// as that will trigger warnings from the go tool.
			return ""
		}
		return def
	}
	return ""
}

func main() {
	// Setup correctly the GOPATH in the environment
	if goPath := os.Getenv("GOPATH"); goPath == "" {
		os.Setenv("GOPATH", defaultGOPATH())
	}

	// Initialize the Fabric SDK
	_, err := blockchain.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v", err)
	}
}
```

The full file is available here: [main.go](main.go)

Like you can see, we fix the GOPATH in the environment if it's not set. We will need this futur for compile the chaincode (we will see this in the next step).

The last thing to do before start the compilation is to use a vendor directory. In our GOPATH we have Fabric, Fabric CA and Fabric SDK Go. All is ok but when we will try to compile our app, there will be some conflict (like multiple definitions of BCCSP). We will handle this by using a tool like `govendor` to flatten these dependencies. Just install it and import external dependencies inside the vendor directory like this:

```
go get -u github.com/kardianos/govendor && \
cd $GOPATH/src/github.com/tohero/heroes-service && \
govendor init && govendor add +external
```

Now we can make the compilation:

```
cd $GOPATH/src/github.com/tohero/heroes-service && \
go build
```

After some time, a new binary named `heroes-service` will appear at the root of the project. Try to start the binary like this:

```
cd $GOPATH/src/github.com/tohero/heroes-service && \
./heroes-service
```

![Screenshot app started but no network](docs/images/start-app-no-network.png)

But this won't work because there is no network that the SDK can talk to. Try to start the network and launch the app again:

```
cd $GOPATH/src/github.com/tohero/heroes-service/fixtures && \
docker-compose up -d && \
cd .. && \
./heroes-service
```

![Screenshot app started and SDK initialized](docs/images/start-app-initialized.png)

Great! We arrive to initialize the SDK to our network created locally. Next step is to interact with a chaincode.

### d. Fabric SDK Go: clean up and Makefile

The Fabric SDK generate some file, like certificates or temporally files. Put down the network won't fully clean up your environment and when you will start again, this files will be reused to avoid recreation. For development you can keep them to test quicly but for a real test, you need to clean up all and start from the begining.

*How clean up my environment ?*

- Put down your network: `cd $GOPATH/src/github.com/tohero/heroes-service/fixtures && docker-compose down`
- Remove MSP folder (defined in the [config](config.yaml) file, in the `fabricCA` section): `rm -rf /tmp/msp`
- Remove enrolment files (defined when we initialize the SDK, in the [setup](blockchain/setup.go) file, when we get the client):  `rm -rf /tmp/enroll_user`

*How to be more productive ?*

We can automatize all these tasks in a single one, same for the build and start. To do so, I propose to use a Makefile. ensure that you have the tool:

```
sudo apt install make
```

Then create a file named `Makefile` at the root of the project with this content:

```
.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@govendor sync
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environnement ..."
	@cd fixtures && docker-compose up --force-recreate -d
	@echo "Sleep 15 seconds in order to let the environment setup correctly"
	@sleep 15
	@echo "Environnement up"

env-down:
	@echo "Stop environnement ..."
	@cd fixtures && docker-compose down
	@echo "Environnement down"

##### RUN
run:
	@echo "Start app ..."
	@./heroes-service

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/enroll_user /tmp/msp heroes-service
	@echo "Clean up done"
```

The full file is available here: [Makefile](Makefile)

Now with the task `all`, first there will be a cleanup of the environment, then the compilation phase, then the network up and finally run the app.

To use it, go to the root of the project and use the `make` command:

- Task `all`: `make` or `make all`
- Task `build`: `make build`
- Task `env-up`: `make env-up`
- ...

### e. Fabric SDK Go: install & instanciate a chaincode

We are very close to use the blockchain system. But for now we haven't setup a chaincode (smart contract) thht will handle queries from our application. First, let's create a new directory named `chaincode` and add a new file named `main.go`:

```
cd $GOPATH/src/github.com/tohero/heroes-service && \
mkdir chaincode && \
vi chaincode/main.go
```

```
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
	if len(args) >= 1 {
		return shim.Error("The number of arguments is insufficient, you need to provide the function to invoke.")
	}

	// In order to manage multiple type of request, we will check the first argument.
	// Here we have only one possible argument: query (every query request will read in the ledger without modification)
	if args[0] == "query" {
		return t.query(stub, args)
	}

	// If the argument given match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
}

// query
// Every readonly functions in the ledger will be here
func (t *HeroesServiceChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check that the number of argument is sufficient to possibly find the function concern by the query
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient, you need to provide the function to query.")
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
	return shim.Error("Unknown query action, check the second argument, must be one of 'index' or 'count'")
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(HeroesServiceChaincode))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
	}
}
```

The full file is available here: [chaincode/main.go](chaincode/main.go)

For now, the chaincode does nothing extrodinary, just put the key/value `hello`/`world` in the ledger at the initialization. In addition, there is one function that we can call by an invoke: `query hello`. This function get the state in the ledger of `hello` and give it in response. We will test this in the next step, after successfully install and instantiate the chaincode.