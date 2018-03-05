# Tutorial Hyperledger Fabric SDK Go: How to build your first app?

```
Source: NEW LINK NEEDED
```

This tutorial will introduce you to the Hyperledger Fabric Go SDK and allows you to build a simple application using the blockchain principle.

**This tutorial uses Hyperledger Fabric version 1.0.5**

This is the **first part** of this tutorial. The basics SDK features will be shown, but the second part is scheduled to demonstrate a more complex application.

## 1. Prerequisites

This tutorial won’t explain in detail how Hyperledger Fabric works. I will just give some tips to understand the general behavior of the framework. If you want to get a full explanation of the tool, go to the official [documentation](http://hyperledger-fabric.readthedocs.io/en/latest/) there is a lot of work there that explains what kind of blockchain Hyperledger Fabric is.

This tutorial has been made on **Ubuntu 16.04** but the Hyperledger Fabric framework is compatible with Mac OS X, Windows and other Linux distributions.

We will use the **Go** language to design our first application because the Hyperledger Fabric has been also built in Go and the Fabric SDK Go is really simple to use. In addition, the chaincode (smart contract) can be written in Go too. So the full-stack will be only in Go! Awesome right ? However, if you are go phobic, there are other SDK like for NodeJS, Java or Python but we won't discuss about them here.

Hyperledger Fabric uses **Docker** to easily deploy a blockchain network. In addition, some components (peers) also deploys docker containers to separate data (channel). So make sure that your platform supports this kind of virtualization.

## 2. Introduction to Hyperledger Fabric

> Hyperledger Fabric is a platform for distributed ledger solutions underpinned by a modular architecture delivering high degrees of confidentiality, resiliency, flexibility and scalability. It is designed to support pluggable implementations of different components and accommodate the complexity and intricacies that exist across the economic ecosystem.

See the full explanation from the official documentation in the introduction part: [Hyperledger Fabric Blockchain](http://hyperledger-fabric.readthedocs.io/en/latest/blockchain.html)

![Blockchain Concensus](http://hyperledger-fabric.readthedocs.io/en/latest/_images/consensus.png)

## 3. Installation guide

This tutorial was made on **Ubuntu 16.04**, but there is help for Windows, Mac OS X and other Linux distributions users.

### a. Docker

**Docker version 17.03.0-ce or greater is required.**

#### Linux (Ubuntu)

First of all, in order to install docker correctly we need to install its dependencies:

```
sudo apt install apt-transport-https ca-certificates curl software-properties-common
```

Once the dependencies are installed, we can install docker:

```
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add - && \
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" && \
sudo apt update && \
sudo apt install -y docker-ce
```

Now we need to manage the current user to avoid using administration rights (`root`) access when we will use the docker command. To do so, we need to add the current user to the `docker` group:

```
sudo groupadd docker ; \
sudo gpasswd -a ${USER} docker && \
sudo service docker restart
```

Do not mind if `groupadd: group 'docker' already exists` error pop up.

To apply the changes made, you need to logout/login. You can then check your version with:

```
docker -v
```

![End of the docker installation](docs/images/finish-docker-install.png)

#### Mac OS X

Download and install the latest [`Docker.dmg`](https://docs.docker.com/docker-for-mac/install/) package for Mac OS X available on the [Docker](https://docs.docker.com/docker-for-mac/install/) website. This will install `docker-compose` as well, so you can skip the next step.

#### Linux (not Ubuntu)

See links below:
- [Debian](https://docs.docker.com/engine/installation/linux/docker-ce/debian/)
- [Fedora](https://docs.docker.com/engine/installation/linux/docker-ce/fedora/)
- [CentOS](https://docs.docker.com/engine/installation/linux/docker-ce/centos/)

#### Windows

See instructions from the Docker website: [docker.com/docker-for-windows](https://docs.docker.com/docker-for-windows/install/)

### b. Docker Compose

**Docker-compose version 1.8 or greater is required.**

We are currently unable to manage easily multiple containers at once. To solve this issue, we need **docker-compose**.

#### Linux

The installation is pretty fast:

```
sudo curl -L https://github.com/docker/compose/releases/download/1.18.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose && \
sudo chmod +x /usr/local/bin/docker-compose
```

Apply these changes by logout/login and then check its version with:

```
docker-compose version
```

![End of the docker compose installation](docs/images/finish-docker_compose-install.png)

#### Windows / Others

See instructions from the Docker-compose website: [docs.docker.com/compose/install](https://docs.docker.com/compose/install/)

### c. Go

**Go version 1.9.x or greater is required.**

#### Linux

You can either follow instructions from [golang.org](https://golang.org/dl/) or use those generics commands that will install Golang 1.9.2 and prepare your environment (generate your `GOPATH`) for Ubuntu:

```
wget https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz && \
sudo tar -C /usr/local -xzf go1.9.2.linux-amd64.tar.gz && \
rm go1.9.2.linux-amd64.tar.gz && \
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile && \
echo 'export GOPATH=$HOME/go' | tee -a $HOME/.bashrc && \
echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' | tee -a $HOME/.bashrc && \
mkdir -p $HOME/go/{src,pkg,bin}
```

To make sure that the installation works, you can logout/login (again) and run:

```
go version
```

![End of the go installation](docs/images/finish-go-install.png)


#### Windows / Mac OS X / Others

See instructions from the Golang website: [golang.org/install](https://golang.org/doc/install)

### d. Fabric SDK Go

Last but not least, the Hyperledger Fabric SDK Go will allow us to easily communicate with the Fabric's components. You don't need to install the Fabric or Fabric CA framework because the SDK automatically handles it locally. To avoid version issues, we are going to checkout to a specific commit which works with the following tutorial.

```
go get -u github.com/hyperledger/fabric-sdk-go && \
cd $GOPATH/src/github.com/hyperledger/fabric-sdk-go && \
git checkout 614551a752802488988921a730b172dada7def1d
```

Let's make sure that you have the requested dependencies:

```
cd $GOPATH/src/github.com/hyperledger/fabric-sdk-go && \
make depend-install
```

Finally, we can launch the various tests of the SDK to check its proper functioning before going further:

```
cd $GOPATH/src/github.com/hyperledger/fabric-sdk-go ; \
make
```

If you get the following error:

```
../fabric-sdk-go/vendor/github.com/miekg/pkcs11/pkcs11.go:29:18: fatal error: ltdl.h: No such file or directory
```

You need to install the package `libltdl-dev` and re-execute previous command (`make`):

```
sudo apt install libltdl-dev
```

The process take a while (depending on your network connection). During this process, a virtual network has been built and some tests have been made in order to check if your system is ready. Now we can work with our first application.

![End of the Fabric SDK Go installation](docs/images/finish-fabric-sdk-go-install.png)

> **Note**: there is more output but it's irrelevant to put it here. The only things you need to care about is the line `fabsdkgo_integration-tests_1 exited with code 0`. If you have `fabsdkgo_integration-tests_1 exited with code 1` then you have a problem. Scroll up in the messages above to find the error.

## 4. Make your first blockchain network

### a. Prepare environment

In order to make a blockchain network, we will use `docker` to build virtual computers that will handle different roles. In this tutorial we will stay as simple as possible. Hyperledger Fabric needs a lot of certificates to ensure encryption during the whole end to end process (TSL, authentications, signing blocks...). The creation of these files requires a little time and in order to go straight to the heart of the matter, we have already prepared all this for you.

Make a new directory in the `src` folder of your `GOPATH`, following our repository naming:

```
mkdir -p $GOPATH/src/github.com/chainHero/heroes-service && \
cd $GOPATH/src/github.com/chainHero/heroes-service
```

You can either follow this command line, which will force you to download subversion package. Or download it from github. 

From command line:
```
sudo apt install subversion && \
cd $GOPATH/src/github.com/chainHero/heroes-service && \
svn checkout https://github.com/chainHero/heroes-service/branches/newVersion/fixtures
```

If you want to download it from GitHub, there is no other choice for you to download the whole repository and extract the fixture folder from it:

- [Github Repository](https://github.com/chainHero/heroes-service/archive/newVersion.zip)

Alternatively, if you wanna know how to build this fixture folder and learn how to create the blockchain network, then we made a tutorial here [coming.soon](coming.soon).


### b. Test

In order to check if the network works, we will use `docker-compose` to start or stop all containers at the same time. Go inside the `fixtures` folder, and run:

```
cd $GOPATH/src/github.com/chainHero/heroes-service/fixtures ; \
docker-compose up
```

You will see a lot of logs with different colors (for your information, red isn't equal to errors).

Open a new terminal and run:

```
docker ps
```

![Docker compose up screenshot](docs/images/docker-ps.png)

You will see : two peers, the orderer, one CA containers and two dependencies (golangruntime and builder) You have successfully made a new network ready to use with the SDK. To stop the network go back to the previous terminal, press `Ctrl+C` and wait that all containers are stopped. If you want to explore more deeper, check out the official documentation about this: [Building Your First Network](http://hyperledger-fabric.readthedocs.io/en/latest/build_network.html)

> **Tips**: when the network is stopped, all containers used remain accessible. This is very useful to check logs for example. You can see them with `docker ps -a`. In order to clean up these containers, you need to delete them with `docker rm $(docker ps -aq)` or if you have used a `docker-compose` file, go where this file is and run `docker-compose down`.

> **Tips**: you can run the `docker-compose` command in background to keep the prompt. To do so, use the parameter `-d`, like this: `docker-compose up -d`. To stop containers, run in the same folder where the `docker-compose.yaml` is, the command: `docker-compose stop` (or `docker-compose down` to clean up after all containers are stopped).


## 5. Use the Fabric SDK Go

### a. Configuration

As we removed the config folder, we need to make a new config file. We will put in it everything that the Fabric SDK Go and our custom parameters for our app needs to work. The config file will contain all our custom parameters and everything else the Fabric SDK Go needs for our app to work. For the moment, we will only try to make the Fabric SDK Go work with the default chaincode:

```
cd $GOPATH/src/github.com/chainHero/heroes-service ; \
vi config.yaml
```

```
name: "chainHero-network"

# Describe what the target network is/does.
description: "The network which will host my first blockchain"

# Schema version of the content. Used by the SDK to apply the corresponding parsing rules.
version: 2

# The client section used by GO SDK.
client:
  # Which organization does this application instance belong to? The value must be the name of an org
  organization: Org1
  logging:
    level: info



# Global configuration for peer, event service and orderer timeouts
  peer:
    timeout:
      connection: 3s
      queryResponse: 45s
      executeTxResponse: 30s
  eventService:
    timeout:
      connection: 3s
      registrationResponse: 3s
  orderer:
    timeout:
      connection: 3s
      response: 5s



  cryptoconfig:
    path: "${GOPATH}/src/github.com/chainHero/heroes-service/fixtures/crypto-config"




  credentialStore:
    path: "/tmp/hfc-heroes-service-kvs"
    cryptoStore:
      path: "/tmp/heroes-service-msp"



   # BCCSP config for the client. Used by GO SDK.
  BCCSP:
    security:
     enabled: true
     default:
      provider: "SW"
     hashAlgorithm: "SHA2"
     softVerify: true
     ephemeral: false
     level: 256



  tlsCerts:
    systemCertPool: false


channels:
  chainhero:
    orderers:
      - orderer.hf.chainhero.io
    peers:
      peer0.org1.hf.chainhero.io:
        endorsingPeer: true
        chaincodeQuery: true
        ledgerQuery: true
        eventSource: true
      peer1.org1.hf.chainhero.io:
        endorsingPeer: true
        chaincodeQuery: true
        ledgerQuery: true
        eventSource: true

    chaincodes:
      # the format follows the "canonical name" of chaincodes by fabric code


# list of participating organizations in this network
organizations:
  Org1:
    mspid: org1.hf.chainhero.io
    cryptoPath:  "peerOrganizations/org1.hf.chainhero.io/users/{userName}@org1.hf.chainhero.io/msp"
    peers:
      - peer0.org1.hf.chainhero.io
      - peer1.org1.hf.chainhero.io
    certificateAuthorities:
      - ca.org1.hf.chainhero.io

# List of orderers to send transaction and channel create/update requests to. For the time being only one orderer is needed.
orderers:
  orderer.hf.chainhero.io:
    url: grpcs://localhost:7050
    grpcOptions:
      ssl-target-name-override: orderer.hf.chainhero.io
      grpc-max-send-message-length: 15
    tlsCACerts:
      path: "${GOPATH}/src/github.com/chainHero/heroes-service/fixtures/crypto-config/ordererOrganizations/hf.chainhero.io/tlsca/tlsca.hf.chainhero.io-cert.pem"

# List of peers to send various requests to, including endorsement, query and event listener registration.
peers:
  peer0.org1.hf.chainhero.io:
    url: grpcs://localhost:7051
    eventUrl: grpcs://localhost:7053
    grpcOptions:
      ssl-target-name-override: peer0.org1.hf.chainhero.io
      grpc.http2.keepalive_time: 15
    tlsCACerts:
      path: "${GOPATH}/src/github.com/chainHero/heroes-service/fixtures/crypto-config/peerOrganizations/org1.hf.chainhero.io/tlsca/tlsca.org1.hf.chainhero.io-cert.pem"

  peer1.org1.hf.chainhero.io:
    url: grpcs://localhost:8051
    eventUrl: grpcs://localhost:8053
    grpcOptions:
      ssl-target-name-override: peer1.org1.hf.chainhero.io
      grpc.http2.keepalive_time: 15
    tlsCACerts:
      path: "${GOPATH}/src/github.com/chainHero/heroes-service/fixtures/crypto-config/peerOrganizations/org1.hf.chainhero.io/tlsca/tlsca.org1.hf.chainhero.io-cert.pem"


certificateAuthorities:
  ca.org1.hf.chainhero.io:
    url: https://localhost:7054
    httpOptions:
      verify: false
    registrar:
      enrollId: admin
      enrollSecret: adminpw
    caName: ca.org1.hf.chainhero.io
```

The full configuration file is available here: [`config.yaml`](config.yaml)

### b. Initialise

We add a new folder named `blockchain` that will contain the whole interface that communicate with the network. We will see the Fabric SDK Go only in this folder.

```
mkdir $GOPATH/src/github.com/chainHero/heroes-service/blockchain
```

Now, we add a new go file named `setup.go` :

```
vi $GOPATH/src/github.com/chainHero/heroes-service/blockchain/setup.go
```

```
package blockchain

import (
        "fmt"
        "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
        "github.com/hyperledger/fabric-sdk-go/pkg/config"
        "time"
        resmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/resmgmtclient"
        chmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/chmgmtclient"
        "github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
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
        admin                   resmgmt.ResourceMgmtClient
        sdk                     *fabsdk.FabricSDK
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func (setup *FabricSetup) Initialize() error {

        // Add parameters for the initialization
        if setup.initialized {
                return fmt.Errorf("sdk already initialized")
        }
	
	var err error
        
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
```

The file is available here: [`blockchain/setup.go`](blockchain/setup.go)

At this stage, we only initialised a client that will communicate to a peer, a CA and an orderer. We also made a new channel and connected this peer to this channel. See the comments in the code for more information.

### c. Test

To make sure that the client managed to initialise all his components, we will make a simple test with the network launched. In order to do this, we need to build the go code. Since we we haven't any main file we have to add one:

```
cd $GOPATH/src/github.com/chainHero/heroes-service ; \
vi main.go
```

```

package main

import (
        "fmt"
        "os"
        "github.com/chainHero/heroes-service/blockchain"
)

func main() {
        // Definition of the Fabric SDK properties
        fSetup := blockchain.FabricSetup{
                // Channel parameters
                ChannelID:              "chainhero",
                ChannelConfig:          "" + os.Getenv("GOPATH") + "/src/github.com/chainHero/heroes-service/fixtures/artifacts/",

                // Chaincode parameters
                OrgAdmin:               "Admin",
                OrgName:                "Org1",
                ConfigFile:             "config.yaml",
        }

        // Initialization of the Fabric SDK from the previously set properties
        err := fSetup.Initialize()
        if err != nil {
                fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
        }
}

```

The file is available here: [`main.go`](main.go)

The last thing to do before starting the compilation is to use a vendor directory. In our GOPATH we have Fabric, Fabric CA, Fabric SDK Go and maybe other projects. When we will try to compile our app, there may be some conflicts (like multiple definitions of BCCSP). We will handle this by using a tool like `dep` to flatten these dependencies. Just install it and import external dependencies inside the vendor directory like this:

Create a file called Gopkg.toml and copy this inside:

```
[[constraint]]
  name = "github.com/hyperledger/fabric-sdk-go"
  revision = "614551a752802488988921a730b172dada7def1d"
```

Save the file and then execute this command: (this command may take a while to proceed)

```
dep ensure
```

Now we can make compile:

```
cd $GOPATH/src/github.com/chainHero/heroes-service ; \
go build
```

After some time, a new binary named `heroes-service` will appear at the root of the project. Try to start the binary like this:

```
cd $GOPATH/src/github.com/chainHero/heroes-service ; \
./heroes-service
```

![Screenshot app started but no network](docs/images/start-app-no-network.png)

At this point, it won't work because there is no network deployed that the SDK can talk with. Start the network and launch the app again:

```
cd $GOPATH/src/github.com/chainHero/heroes-service/fixtures ; \
docker-compose up -d ; \
cd .. ; \
./heroes-service
```

![Screenshot app started and SDK initialised](docs/images/start-app-initialized.png)
Note: You need to see "Success". If it's not the case then something went wrong.

Alright ! So we just initialised the SDK with our local network. In the next step, we will interact with a chaincode.

### d. Clean up and Makefile

The Fabric SDK generates some files, like certificates and/or temporally files. Shutting down the network won't fully clean up your environment and when you will need to start it again, these files will be reused to avoid building process. For development you can keep them to test quickly but for a real test, you need to clean up all and start from the beginning.

*How to clean up my environment ?*

- Shut down your network: `cd $GOPATH/src/github.com/chainHero/heroes-service/fixtures && docker-compose down`
- Remove MSP folder (defined in the [config](config.yaml) file, in the `fabricCA` section): `rm -rf /tmp/msp`
- Remove enrolment files (defined when we initialise the SDK, in the [setup](blockchain/setup.go) file, when we get the client):  `rm -rf /tmp/enroll_user`
- Remove some docker containers and docker images not generated by the `docker-compose` command: ```docker rm -f -v `docker ps -a --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null``` and ```docker rmi `docker images --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null```

*How to be more efficient ?*

We can automatise all these tasks in one single step. Also the build and start process can be automated. To do so, we will create a Makefile. First, ensure that you have the tool:

```
make --version
```
If `make` is not installed do (Ubuntu):

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
	@dep ensure
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
	@docker rm -f -v `docker ps -a --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"
```

The file is available here: [`Makefile`](Makefile)

Now with the task `all`:
1. the whole environment will be cleaned up,
2. then our go program will be compiled,
3. after which the network will be deployed and
4. finally the app will be up and running.

To use it, go in the root of the project and use the `make` command:

- Task `all`: `make` or `make all`
- Task `build`: `make build`
- Task `env-up`: `make env-up`
- ...

### e. Install & instanciate the chaincode

We are almost there to use the blockchain system. But for now we haven't set up any chaincode (smart contract) yet that will handle queries from our application. First, let's create a new directory named `chaincode` and add a new file named `main.go`:

```
cd $GOPATH/src/github.com/chainHero/heroes-service ; \
mkdir chaincode ; \
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
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *HeroesServiceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### HeroesServiceChaincode Init ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call")
	}

	// Put in the ledger the key/value hello/world
	err := stub.PutState("hello", []byte("world"))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return a successful message
	return shim.Success(nil)
}

// Invoke
// All future requests named invoke will arrive here.
func (t *HeroesServiceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### HeroesServiceChaincode Invoke ###########")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	// Check whether it is an invoke request
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// In order to manage multiple type of request, we will check the first argument.
	// Here we have one possible argument: query (every query request will read in the ledger without modification)
	if args[0] == "query" {
		return t.query(stub, args)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
}

// query
// Every readonly functions in the ledger will be here
func (t *HeroesServiceChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### HeroesServiceChaincode query ###########")

	// Check whether the number of arguments is sufficient
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Like the Invoke function, we manage multiple type of query requests with the second argument.
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

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown query action, check the second argument.")
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(HeroesServiceChaincode))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
	}
}
```

> The chaincode isn't really related to the application, we can have one repository for the app and another for the chaincode. For your information, in a near future, the chaincode could be written in other languages.

For now, the chaincode does nothing extraordinary, just put the key/value `hello`/`world` in the ledger at initialisation. In addition, there is one function that we can call by an invoke: `query hello`. This function gets the state of the ledger, i.e. `hello` and give it in response. We will test this in the next step, after successfully install and instantiate the chaincode.

In order to install and instantiate the chaincode, we need to add some code in the application. Edit the [`blockchain/setup.go`](blockchain/setup.go) with those following lines:

```
package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	"time"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/ccpackager/gopackager"
	resmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/resmgmtclient"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	chmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/chmgmtclient"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
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
	admin 		resmgmt.ResourceMgmtClient
	sdk 		*fabsdk.FabricSDK
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
```

The file is available here: [`blockchain/setup.go`](blockchain/setup.go)

> **Tips**: take care of the chaincode version, if you want to update your chaincode, increment this number. Otherwise the network will keep the same chaincode.

We need now to modify our setup structure in the main file

You can directly edit the file by modifying those lines inside

```
	fSetup := blockchain.FabricSetup{
		// Channel parameters
		ChannelID:        	"chainhero",
		ChannelConfig:    	"" + os.Getenv("GOPATH") + "/src/github.com/chainHero/heroes-service/fixtures/artifacts/",

		// Chaincode parameters
		ChainCodeID:      	"heroes-service",
		ChaincodeGoPath:  	os.Getenv("GOPATH"),
		ChaincodePath:    	"github.com/chainHero/heroes-service/chaincode/",
		OrgAdmin:		"Admin",
		OrgName:		"Org1",
		ConfigFile:		"config.yaml",
	}
```

The file is available here: [`blockchain/setup.go`](blockchain/setup.go)


We can test this, just with the `make` command setup in the previous step:

```
cd $GOPATH/src/github.com/chainHero/heroes-service ; \
make
```

![Screenshot Chaincode installed and instantiated](docs/images/install-and-instantiate-cc.png)

> **Tips**: the installation and the instantiation don't need to be run at every start of the application. Only when we update the chaincode (and the chaincode version). A solution is to provide an argument when we run the application to tell to do this additional procedure before move on. Since in this tutorial we will clean up the environment every time we don't really care about that.

### f. Query the chaincode

Like a database, the chaincode is plugged and ready to answer. Let's try the `hello` query.

We will put all query functions in a new file named `query.go` in the `blockchain` folder:

```
cd $GOPATH/src/github.com/chainHero/heroes-service ; \
vi blockchain/query.go
```

```
package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	"fmt"
)

// QueryHello query the chaincode to get the state of hello
func (setup *FabricSetup) QueryHello() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, "hello")

	response, err := setup.client.Query(chclient.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}})
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	return string(response.Payload), nil
}
```

The file is available here: [`blockchain/query.go`](blockchain/query.go)

You can add the call to this new function in the [`main.go`](main.go):

```
func main() {

[...]

	// Query the chaincode
	response, err := fSetup.QueryHello()
	if err != nil {
		fmt.Printf("Unable to query hello on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Response from the query hello: %s\n", response)
	}
}
```

The file is available here: [`main.go`](main.go)

Let's try:

```
cd $GOPATH/src/github.com/chainHero/heroes-service ; \
make
```

![Screenshot Query Hello](docs/images/query-hello.png)

### g. Change the ledger state

The next thing to do in order to make a basic tour of the Fabric SDK Go, is to make a request to the chaincode in order to change the ledger state.

First, we will add this ability in the chaincode. Edit the [`chaincode/main.go`](chaincode/main.go) file:

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
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *HeroesServiceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### HeroesServiceChaincode Init ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call")
	}

	// Put in the ledger the key/value hello/world
	err := stub.PutState("hello", []byte("world"))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return a successful message
	return shim.Success(nil)
}

// Invoke
// All future requests named invoke will arrive here.
func (t *HeroesServiceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### HeroesServiceChaincode Invoke ###########")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	// Check whether it is an invoke request
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
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

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
}

// query
// Every readonly functions in the ledger will be here
func (t *HeroesServiceChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### HeroesServiceChaincode query ###########")

	// Check whether the number of arguments is sufficient
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Like the Invoke function, we manage multiple type of query requests with the second argument.
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

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown query action, check the second argument.")
}

// invoke
// Every functions that read and write in the ledger will be here
func (t *HeroesServiceChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### HeroesServiceChaincode invoke ###########")

	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Check if the ledger key is "hello" and process if it is the case. Otherwise it returns an error.
	if args[1] == "hello" && len(args) == 3 {

		// Write the new value in the ledger
		err := stub.PutState("hello", []byte(args[2]))
		if err != nil {
			return shim.Error("Failed to update state of hello")
		}

		// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
		err = stub.SetEvent("eventInvoke", []byte{})
		if err != nil {
			return shim.Error(err.Error())
		}

		// Return this value in response
		return shim.Success(nil)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown invoke action, check the second argument.")
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(HeroesServiceChaincode))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
	}
}
```

The file is available here: [`chaincode/main.go`](chaincode/main.go)

From the application side, we add a new interface to make the invocation of the chaincode. Add a file named `invoke.go` in the `blockchain` folder:

```
cd $GOPATH/src/github.com/chainHero/heroes-service ; \
vi blockchain/invoke.go
```

```
package blockchain

import (
	"fmt"
	"time"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
)

// InvokeHello
func (setup *FabricSetup) InvokeHello(value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "invoke")
	args = append(args, "hello")
	args = append(args, value)

	eventID := "eventInvoke"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")

	// Register a notification handler on the client
	notifier := make(chan *chclient.CCEvent)
	rce, err := setup.client.RegisterChaincodeEvent(notifier, setup.ChainCodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("failed to register chaincode evet: %v", err)
	}

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(chclient.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}, TransientMap:transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	// Unregister the notification handler previously created on the client
	err = setup.client.UnregisterChaincodeEvent(rce)

	return response.TransactionID.ID, nil
}
```

The file is available here: [`blockchain/invoke.go`](blockchain/invoke.go)

You can then add the call to this function in the [`main.go`](main.go):

```
func main() {

[...]

	// Query the chaincode
	response, err := fSetup.QueryHello()
	if err != nil {
		fmt.Printf("Unable to query hello on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Response from the query hello: %s\n", response)
	}

	// Invoke the chaincode
	txId, err := fSetup.InvokeHello("chainHero")
	if err != nil {
		fmt.Printf("Unable to invoke hello on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Successfully invoke hello, transaction ID: %s\n", txId)
	}

	// Query again the chaincode
	response, err = fSetup.QueryHello()
	if err != nil {
		fmt.Printf("Unable to query hello on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Response from the query hello: %s\n", response)
	}
}
```

Let's try:

```
cd $GOPATH/src/github.com/chainHero/heroes-service ; \
make
```

// TODO executer la fonction précedente et copié l'output
![Screenshot Invoke Hello](docs/images/invoke-hello.png)

## 6. Make this in a web application

We also can make this usable for any user. The best choice is a web application and we are lucky because the Go language natively provides a web server handling HTTP requests and also templating for HTML.

For now, we have only two different actions: the query and the invocation of the hello value. Let's make two HTML pages for each action. We add a [`web`](web) directory with three other directories:
- [`web/templates`](web/templates): contains all HTML pages (templates)
- [`web/assets`](web/assets): contains all CSS, Javascript, Fonts, Images...
- [`web/controllers`](web/controllers): contains all functions that will render templates

We use the MVC (Model-View-Controller) to make it more readable. The Model will be the blockchain part, the View are templates and Controller are provided by functions in the [`controllers`](web/controllers) directory.

Populate each with the appropriate code (we also added Bootstrap to make the result a little prettier:

- [`web/templates/layout.html`](web/templates/layout.html)
- [`web/templates/home.html`](web/templates/home.html)
- [`web/templates/request.html`](web/templates/request.html)
- [`web/controllers/controller.go`](web/controllers/controller.go)
- [`web/controllers/home.go`](web/controllers/home.go)
- [`web/controllers/request.go`](web/controllers/request.go)
- [`web/app.go`](web/app.go)
- [`web/assets`](web/assets)

And finaly, we change the [`main.go`](main.go), in order to use the web interface instead of directly query the blockchain.

- [`main.go`](main.go)

Run the app and go to [localhost:3000/home.html](http://localhost:3000/home.html):

```
cd $GOPATH/src/github.com/chainHero/heroes-service ; \
make
```

The `home` page make a query in in the blockchain to get the value of the `hello` key and display it.

![Screenshot Web Home Hello World](docs/images/web-home-hello-world.png)

The `request` page has a form to change the `hello` value. 

![Screenshot Web Request Write](docs/images/web-request-write.png)

After a successful submission the transaction ID is given.

![Screenshot Web Request Success](docs/images/web-request-success.png)

We can see the change by going back to the `home` page.

![Screenshot Web Home Hello Superman](docs/images/web-home-hello-superman.png)

**It's the end for the first part. A more complex application is coming.**

## 7. References

- [Hyperledger website](https://www.hyperledger.org/)
- [Hyperledger Fabric online documentation](http://hyperledger-fabric.readthedocs.io/en/latest/)
- [Hyperledger Fabric on github](https://github.com/hyperledger/fabric)
- [Hyperledger Fabric Certificate Authority on github](https://github.com/hyperledger/fabric-ca)
- [Hyperledger Fabric SDK Go on github](https://github.com/hyperledger/fabric-sdk-go)
- [Fabric SDK Go tests](https://github.com/hyperledger/fabric-sdk-go/blob/master/test/integration/end_to_end_test.go)
- [CLI](https://github.com/securekey/fabric-examples/tree/master/fabric-cli/): An example CLI for Fabric built with the Go SDK.
