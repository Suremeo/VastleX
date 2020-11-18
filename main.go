package main

import (
<<<<<<< Updated upstream
	"github.com/vastlellc/vastlex/log"
	"github.com/vastlellc/vastlex/vastlex"
=======
	"github.com/VastleLLC/VastleX/vastlex"
	"github.com/VastleLLC/VastleX/vastlex/logging"
>>>>>>> Stashed changes
)

// main starts the proxy.
func main() {
	log.DefaultLogger.Fatal(vastlex.Start())
}