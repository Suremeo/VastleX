package ddos

import (
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"go.uber.org/atomic"
	"strconv"
	"time"
)

type W struct{}

var Count = atomic.NewInt64(0)

func init() {
	t := time.NewTicker(10 * time.Second)
	t.Stop()
	go func() {
		for {
			<-t.C
			c := Count.Swap(0)
			if c > 100 {
				// we getting ddosed (Probably)
				log.DefaultLogger.Debug("DDOS detected, invalid packet count (last 10 seconds): " + strconv.Itoa(int(c)))
			}
		}
	}()
}

func (*W) Write(b []byte) (int, error) {
	Count.Inc()
	return 0, nil
}
