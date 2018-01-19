#!/usr/bin/env bash

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

InstallDependancies() {
    echo "Install dependancies"
    echo -e ${GREEN}INSTALLATION

    echo -e ${GREEN}Step I - Docker
    echo -e ${NC}

    sudo apt-get update
    sudo apt-get install apt-transport-https ca-certificates curl software-properties-common
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    sudo apt-get update
    sudo apt-get install docker-ce

    sudo groupadd docker
    sudo gpasswd -a ${USER} docker
    sudo service docker restart


    echo -e ${GREEN}Step II - Docker-compose
    echo -e ${NC}

    sudo curl -L https://github.com/docker/compose/releases/download/1.18.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/doc$
    sudo chmod +x /usr/local/bin/docker-compose

    echo -e ${GREEN}Step III - Golang
    echo -e ${NC}

    wget https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz && \
    sudo tar -C /usr/local -xzf go1.9.2.linux-amd64.tar.gz && \
    rm go1.9.2.linux-amd64.tar.gz && \
    echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile && \
    echo 'export GOPATH=$HOME/go' | tee -a $HOME/.bashrc && \
    echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' | tee -a $HOME/.bashrc && \
    mkdir -p $HOME/go/src $HOME/go/pkg $HOME/go/bin

    echo -e ${GREEN}Dependancies successfully installed, reboot your computer & relaunch this script.
}

InstallBlockChain()
{
    echo -e ${GREEN}Step III - Hyperledger Fabric \& Certificate Authority \(CA\)
    echo -e ${NC}

    mkdir -p $GOPATH/src/github.com/hyperledger && \
    cd $GOPATH/src/github.com/hyperledger && \
    git clone https://github.com/hyperledger/fabric.git && \
    cd fabric && \
    git checkout v1.0.5


    echo -e ${GREEN}Step IV - Fabric SDK Go
    echo -e ${NC}

    go get -u github.com/hyperledger/fabric-sdk-go
    go get github.com/hyperledger/fabric-sdk-go/pkg/fabric-client
    go get github.com/hyperledger/fabric-sdk-go/pkg/fabric-ca-client
    cd $GOPATH/src/github.com/hyperledger/fabric-sdk-go && make depend-install
    cd $GOPATH/src/github.com/hyperledger/fabric-sdk-go && make


    echo -e ${GREEN}Installation DONE
    echo -e ${NC}
}

if docker --version &> /dev/null
then
    InstallBlockChain
else
    InstallDependancies
fi


