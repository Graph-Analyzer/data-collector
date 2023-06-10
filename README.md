# Graph Analyzer - Data Collector

The data collector feeds graphs from various sources (currently GEXF & Jalapeno) into the Neo4j database.
It also handles the Neo4j GDS graph creation.

## Run Locally

### Checkout the project

```zsh
git clone https://github.com/Graph-Analyzer/data-collector.git
```

### Code Formatting
Code formatting is done using [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports).
Installation instructions can be found in the official documentation.

Run the code formatting with the following command:

```zsh
goimports -w .
```

### Build

```zsh
go build
```

## Protobuf / gRPC

To update the gRPC files based on the Protobuf definitions, follow [this](https://grpc.io/docs/languages/go/quickstart/) guide
to install all locally required dependencies and run:

```zsh
protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.
```

## Config

The configuration can be set statically or dynamically.

### Env File

Create a `.env` file locally, based on the provided `.env.example`.

### Flags

Pass the appropriate `--flags`. A complete list of flags can be found in the usage section.

## Usage

```zsh
# ./data-collector --help
data-collector is used to import network graphs to neo4j

Usage:
  data-collector [flags]
  data-collector [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  gexf        Use GEXF file format
  help        Help about any command
  jalapeno    Use jalapeño gateway

Flags:
      --config string           Config file (default is .env)
  -h, --help                    help for data-collector
      --neo4j-host string       Neo4j Host (including protocol neo4j/bolt)
      --neo4j-password string   Neo4j Password
      --neo4j-port string       Neo4j Port
      --neo4j-realm string      Neo4j Realm
      --neo4j-user string       Neo4j Username
      --network-name string     Display name of the network (ASCII, 50 characters max) (default "default")
  -v, --verbose count           Increase output verbosity. Example: --verbose=2 or -vv

Use "data-collector [command] --help" for more information about a command.
```

## GEXF
Make sure that you use a GEXF file with weighted edges. Undefined weights will be interpreted as 0.

### File
```zsh
# ./data-collector gexf -f testgraph/graph_sa.gexf         
INFO[2022-12-05T21:29:44+01:00] Database connection established              
INFO[2022-12-05T21:29:45+01:00] Using GEXF file testgraph/graph_sa.gexf      
INFO[2022-12-05T21:29:46+01:00] Parsed 8 nodes from input file               
INFO[2022-12-05T21:29:47+01:00] Parsed 26 edges from input file              
INFO[2022-12-05T21:29:47+01:00] Finished importing GEXF file   
```

### GRPC
```zsh
# ./data-collector gexf -l
INFO[2023-04-18T16:33:02+02:00] Database connection established
INFO[2023-04-18T16:33:03+02:00] Using GEXF gRPC listener
```

## Jalapeno

```zsh
# ./data-collector jalapeno                                                                                                                                                   12s
INFO[2022-12-05T21:30:20+01:00] Database connection established              
INFO[2022-12-05T21:30:20+01:00] Using Jalapeno                               
INFO[2022-12-05T21:30:20+01:00] Subscribe to LsNodesEdges                    
INFO[2022-12-05T21:30:20+01:00] Subscribe to LsNodes                         
INFO[2022-12-05T21:30:21+01:00] Received 8 nodes from JAGW getLsNodes        
INFO[2022-12-05T21:30:21+01:00] Received 26 edges from JAGW getLsNodeEdges   
INFO[2022-12-05T21:30:22+01:00] Start processing node updates                
INFO[2022-12-05T21:30:22+01:00] Start processing edge updates                
... (waiting for gRPC stream events)
```

### Mocking Jalapeno gRPC connection

```zsh
# Terminal 1
npm install -g camouflage-server
cd camouflage
camouflage --config config.yml
```

```zsh
# Terminal 2
go build && ./data-collector jalapeno  --jagw-host 127.0.0.1 --jagw-request-port 4312  --jagw-subscription-port 4312
```

See `camouflage/grpc/mocks` for how to define new responses

## Authors

- [@lribi](https://github.com/lribi)
- [@pesc](https://github.com/pesc)

## License

This project is licensed under the [MIT](https://github.com/Graph-Analyzer/api/blob/data-collector/LICENSE)
License.

### Third Party Licenses

Third party licenses can be found in `THIRD-PARTY-LICENSES.txt`.
Regenerate them with this command.

```zsh
# Make sure you have a compiled version of the project in the project root
go install github.com/uw-labs/lichen@latest

./scripts/go-licenses.sh
```
