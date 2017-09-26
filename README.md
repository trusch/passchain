passchain
=========

Passchain is a tool to securely store and share passwords, tokens and other short secrets on a private blockchain.

This is a [tendermint](https://github.com/tendermint/tendermint) application and provides the following functionality:

## Features

* create, delete and show accounts
* create, delete and show secrets
* share secrets with other accounts
* share only read- or read-write access
* share secrets with group accounts

## User interface

* command line client
* all functionality is available

## Open Tasks

* write HTTP API
* write graphical UI  

## Install
```
# install via go get
go get github.com/trusch/passchain/cmd/...

# if dependency errors occur fix them by using glide:
cd $GOPATH/src/github.com/trusch/passchain
glide install
go install github.com/trusch/passchain/cmd/...

# if you need to install tendermint:
go get github.com/tendermint/tendermint/cmd/...

# if dependency errors occur fix them by using glide:
cd $GOPATH/src/github.com/tendermint/tendermint
glide install
go install github.com/tendermint/tendermint/cmd/...
```

## Walktrough

```
# start tendermint and passchain-abci
tendermint init
tendermint node --consensus.create_empty_blocks=false &
passchain-abci &

# create account and store keys in environment
passchain accounts create --id alice
export PASSCHAIN_ID=alice
export PASSCHAIN_PUBLIC_KEY=BLBoMJ+iUjremFnReUF0onvmokhV4Hmtvq+fU24oxSYIhwIePlHXZTbW28PZN66i3Nvc7txv4lXH4KwlY9O4KFo=
export PASSCHAIN_PRIVATE_KEY=7FYJr0/50gsjPvuQGmk2atj8QsdzYvbyrGrH2KTuhxY=

# create secret
passchain secrets create my-secret "this is secret"

# read secret
passchain secrets get my-secret

# share secret with bob
passchain secrets share my-secret --with bob
```

## Howto use groups
```
# create group account
passchain accounts create --id my-group

# store private key as secret with the very same name
passchain secrets create my-group "7FYJr0/50gsjPvuQGmk2atj8QsdzYvbyrGrH2KTuhxY="

# share group key with other accounts
passchain secrets share my-group --with bob

# share secret with group
passchain secrets share my-secret --with my-group

# retrieve secret with group key
passchain --as my-group secrets get my-secret

```
