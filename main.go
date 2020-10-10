package main

import (
	"github.com/vastlellc/vastlex/log"
	"github.com/vastlellc/vastlex/vastlex"
)

// main starts the proxy.
func main() {
	err := vastlex.Start()
	if err != nil {
		log.FatalError("VastleX crashed!", err)
	}
}
