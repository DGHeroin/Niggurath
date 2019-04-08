package benchmark

import (
	"github.com/DGHeroin/Niggurath/event"
	"log"
	"sync/atomic"
	"testing"
	"time"
)

type Handler struct {}

func (h *Handler) OnMessage(message string) (result string, err error) {
	atomic.AddUint64(&counter, 1)
	return
}

func aTestQPS(t *testing.T) {
	handler := Handler {}
	event.GetEventManager().AddListeners("rpc.test", &handler)
	isRunning = true
	producer := func() {
		for {
			if !isRunning {
				break
			}
			mgr := event.GetEventManager()
			mgr.DispatchEvents("rpc.test", "")
		}
	}

	monitor := func() {
		ticker := time.NewTicker(time.Second)
		tick := 0
		for range ticker.C {
			log.Printf("multi listeners qps: %v", counter)
			atomic.StoreUint64(&counter, 0)
			tick++
			if tick > 5 {
				isRunning = false
				ticker.Stop()
				break
			}
		}
	}
	go producer()
	monitor()
}


func cTestSingleListenerQPS(t *testing.T) {
	handler := Handler {}
	event.GetEventManager().AddSingleListener("rpc.test", &handler)
	mgr := event.GetEventManager()
	isRunning = true
	producer := func() {
		for {
			if !isRunning {
				break
			}
			mgr.DispatchSingleEvent("rpc.test", "")
		}
	}

	monitor := func() {
		ticker := time.NewTicker(time.Second)
		tick := 0
		for range ticker.C {
			log.Printf("single listener qps: %v", counter)
			atomic.StoreUint64(&counter, 0)
			tick++
			if tick > 5 {
				isRunning = false
				ticker.Stop()
				break
			}
		}
	}
	go producer()

	monitor()

}