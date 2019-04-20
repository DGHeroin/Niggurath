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

func TestEventQueueQPS(t *testing.T) {
    handler := Handler {}
    event.GetEventManager().AddListeners("rpc.test.multi", &handler)
    isRunning = true
    producer := func() {
        for {
            if !isRunning {
                break
            }
            mgr := event.GetEventManager()
            mgr.DispatchEvents("rpc.test.multi", "hello")
        }
    }

    monitor := func() {
        ticker := time.NewTicker(time.Second)
        tick := 0
        for range ticker.C {
            log.Printf("multi listeners qps: %v(%v)", counter, size(counter))
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


func TestSingleListenerQPS(t *testing.T) {
    handler := Handler {}
    event.GetEventManager().AddSingleListener("rpc.test.single", &handler)
    mgr := event.GetEventManager()
    isRunning = true
    producer := func() {
        for {
            if !isRunning {
                break
            }
            mgr.DispatchSingleEvent("rpc.test.single", "")
        }
    }

    monitor := func() {
        ticker := time.NewTicker(time.Second)
        tick := 0
        for range ticker.C {
            log.Printf("single listener qps: %v(%v)", counter, size(counter))
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