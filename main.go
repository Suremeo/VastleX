package main

import (
	"github.com/VastleLLC/VastleX/vastlex"
	"github.com/VastleLLC/VastleX/vastlex/logging"
)

// main starts the proxy.
func main() {
	log.DefaultLogger.Fatal(vastlex.Start())
}
