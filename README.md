# Forecasting blockchain solution using Hyperledger Fabric 1.0

Create a network to jump start development of your decentralized application.

The network can be deployed to multiple docker containers on one host for development or to multiple hosts for testing 
or production.

Scripts of this starter generate crypto material and config files, start the network and deploy your chaincodes. 
Developers can use admin web app of **REST API server**

What's left is to develop your chaincodes and place them into the [chaincode](./chaincode) folder.

Most of the plumbing work is taken care of by this starter.

## Members and Components

Network consortium consists of:

- Orderer organization `walmartlabs.com`
- Peer organization org1 `walmart` 
- Peer organization org2 `unilever` 


They transact with each other on the following channels:


  - `walmart-unilever`

Each organization starts several docker containers:

- **peer0** (ex.: `peer0.unilever.walmartlabs.com`) with the anchor [peer](https://github.com/hyperledger/fabric/tree/release/peer) runtime
- **peer1** `peer1.unilever.walmartlabs.com` with the secondary peer
- **ca** `ca.unilever.walmartlabs.com` with certificate authority server [fabri-ca](https://github.com/hyperledger/fabric-ca)
- **api** `api.unilever.walmartlabs.com` with [fabric-rest](https://gecgithub01.walmart.com/a0b013g/hyperfabric-rest.git) API server
- **www** `www.unilever.walmartlabs.com` with a simple http server to serve members' certificate files during artifacts generation and setup
- **cli** `cli.unilever.walmartlabs.com` with tools to run commands during setup

## Local deployment

Deploy docker containers of all member organizations to one host, for development and testing of functionality. 

All containers refer to each other by their domain names and connect via the host's docker network. The only services 
that need to be available to the host machine are the `api` so you can connect to admin web apps of each member; 
thus their `4000` ports are mapped to non conflicting `4000, 4001, 4002` ports on the host.

Generate artifacts:
```bash
./network.sh -m generate
```

Generated crypto material of all members, block and tx files are placed in shared `artifacts` folder on the host.

Please wait for 7 mins before making the containers online as the certificates are generated 7 mins in future.

Start the fabric docker containers of all members, this will start the blockchain network:
```bash
./network.sh -m up
```

Once the fabric container are up, next step is to start the API servers for each organization.

```bash
./network.sh -m api-up
```

Once the API servers are up, next step would be to install all the smart contracts (chaincode) to there respective nodes and join the channels.

```bash
./network.sh -m install
./network.sh -m join
```

After all containers are up, browse to each member's admin web app to transact on their behalf: 

- Walmart [http://localhost:4000/admin](http://localhost:4000/admin)
- Unilever [http://localhost:4001/admin](http://localhost:4001/admin)

Tail logs of each member's docker containers by passing its name as organization `-o` argument:
```bash
# orderer
./network.sh -m logs -m walmartlabs.com

# members
./network.sh -m logs -m walmart
./network.sh -m logs -m unilever
```
Stop all:
```bash
./network.sh -m down
```
Remove dockers:
```bash
./network.sh -m clean
```

