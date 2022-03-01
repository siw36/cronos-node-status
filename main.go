package main

import "github.com/siw36/cronos-node-status/router"

var source = "https://cronos.org/explorer/chain-blocks"

func main() {
	router.Serve()
}
