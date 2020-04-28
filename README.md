# UP Voting System

This is a blockchain-based implementation of the voting protocol proposed in the paper "Verifiable 
internet elections with everlasting privacy and minimal trust" by Philipp Locher and Rolf Haenni. 
The acronym UP stands for Unconditional Privacy, a property this protocol has besides individual and 
universal verifiability. 

An important component of the prototype is the Public Bulletin Board (PBB) where all election
relevant data is stored. It blockchain-based and implemented with 
[Tendermint](https://tendermint.com/) and [Cosmos-SDK](https://docs.cosmos.network/master/) to 
satisfy the protocol's minimal trust assumptions.

Three Zero-knowledge Proofs are at the heart of the unconditional privacy property and are also
part of this implementation. They were originally implemented in the 
[UniCrypt](https://github.com/bfh-evg/unicrypt) library which served as a template.


## Build 

Requires **Go 1.13.0+**.

Make sure that your GOPATH is set and GOPATH/bin is in your PATH. 
This project makes use of go modules. 

Clone the repository:
```zsh
mkdir -p $GOPATH/src/github.com/csmuller
cd $GOPATH/src/github.com/csmuller
git clone https://github.com/csmuller/up-voting-system.git
cd up-voting-system
```

Install the PBB, voter client and administration client applications by using the included Makefile
This installs the binaries into $GOPATH/bin or $GOBIN if set.
```zsh
make install
```


## Run

Now you should be able to run the following commands:

```
pbbd help
vcli help
acli help
```

