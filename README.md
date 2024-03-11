# Tcp server with PoW DDOS protection

This is a simple tcp server with PoW DDOS protection. It is a simple server that listens for incoming connections and
checks if the client has done a PoW before accepting the connection.

## PoW algorithm explanation

I was comparing [Hashcash](https://en.wikipedia.org/wiki/Hashcash)
to [Guided_tour_puzzle_protocol](https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol)

I choose Hashcash becasue of the simplicity both in implementing and understanding, Guided tour puzzle protocol involves
a lot of rounds of rpc calls.

## How to build

For building docker images:
`make dockerbuild`
`make dockerbuildclient`

## How to use

1. It can be used just as go binaries.
   Server: `go run ./cmd/server/main.go`
   Client: `go run ./cmd/client/main.go localhost:8080`
2. Can be used as docker containers. After you have docker images:
   `docker-compose up` - will start both server and client, you can find the quote in client logs.
   `docker compose run client` - will restart the client

## Possible improvements

1. I didn't have time to cover the code with unit tests. I would like to add tests for the server and the PoW algorithm.
2. For distributed environment in-memory cache can be replaced with Redis.
3. A lot of parameters are hardcoded. It would be nice to move them to the configuration file.
