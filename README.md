# Tutorial Hyperledger Fabric SDK Go: How to build your first app?

This tutorial will introduce you to the Hyperledger Fabric Go SDK and allow you to build a simple application using the blockchain principle. The application will be a tool that allows people to make a request for a get some help. A hero can see all requests and accept some. When he achieved the service, the person who made the request can thank and reward the hero. The kind of things is like a smart contract between a simple user and a hero.

## 1. Prerequisites

This tutorial won’t explain in detail how Hyperledger Fabric works, I will just give some tips to understand the general behavior of the framework. If you want to get a full explanation of the tool, go to the official [documentation](http://hyperledger-fabric.readthedocs.io/en/latest/) there is a lot of work that explains to you what kind of blockchain is Hyperledger Fabric.

In the technical part, this tutorial has been made on **Ubuntu 16.04**. The Hyperledger Fabric framework is compatible with Mac OSX and Windows too, but I can’t guarantee that all the stuff can work.

We will use the **Go** language to design a first application, because the Hyperledger Fabric has been built also in Go and the Fabric SDK Go is really simple to use. There are other SDK if you want to, like for NodeJS, Java or Python.

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

In order to make it work, we have to edit the docker-composer.yaml file. This is the configuration file for docker-compose, it tells what containers need to be created and started and with a custom configuration for each. Take your favorite text editor and copy paste content from this repository:

```
vi fixtures/docker-composer.yaml
```

see [fixtures/docker-composer.yaml](fixtures/docker-compose.yaml)

Now if we use docker-compose we will setup 2 fabric certificate authorities with 1 peer for each. Peers will have all roles: ledger, endorer and commiter. In addition, an orderer is also created with the `solo` ordering (no consensus is made).

In our example, one organisation will be the heroes and the other concern people who made request. So lets say that peer 0 of the organisation 1 is superman and peer 0 of the organisation 2 is John, a normal citizen.

### c. Test

*TODO*

## 5. Build a simple application using SDK

### a. Configuration of the Fabric SDK

Like we remove the config folder, we need to make a new config file:

```
vi config.yaml
```

see [config.yaml](config.yaml)
