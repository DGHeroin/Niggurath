package benchmark

import (
    "context"
    "github.com/DGHeroin/Niggurath/event"
    "github.com/DGHeroin/Niggurath/rpc/client"
    "github.com/DGHeroin/Niggurath/rpc/server"
    "log"
    "sync/atomic"
    "testing"
    "time"
)

type EchoHandler struct {
}

func (h *EchoHandler) OnMessage(message string) (result string, err error) {
    atomic.AddUint64(&counter, 1)
    return message, nil
}


func TestRPC(t *testing.T) {
    handler := EchoHandler{}
    event.GetEventManager().AddListeners("echo.handler", &handler)
    isRunning = true
    srv := server.Server{}
    defer srv.Close()
    go func() {
        srv.Serve(":7878")
    }()

    cli := client.Client{}

    sender := func() {
        err := cli.Connect("127.0.0.1:7878")
        if err != nil {
            log.Println(err)
            return
        }

        for {
            if !isRunning { break }
            cli.Execute(context.Background(), "echo.handler", "hello")
        }
    }

    monitor := func() {
        ticker := time.NewTicker(time.Second)
        tick := 0
        for range ticker.C {
            log.Printf("multi rpc qps: %v(%v)", counter, size(counter))
            atomic.StoreUint64(&counter, 0)
            tick++
            if tick > 3 {
                isRunning = false
                ticker.Stop()
                break
            }
        }
    }
    for i:=0; i < rpcSenderCoroutineCount; i++ {
        go sender()
    }
    monitor()

}

func TestRPCSingleEvent(t *testing.T) {
    handler := EchoHandler{}
    event.GetEventManager().AddListener("echo.handler", &handler)
    isRunning = true
    srv := server.Server{}
    defer srv.Close()
    go func() {
        srv.Serve(":7878")
    }()

    cli := client.Client{}

    sender := func() {
        err := cli.Connect("127.0.0.1:7878")
        if err != nil {
            log.Println(err)
            return
        }
        for {
            if !isRunning { break }
            cli.ExecuteSingleEvent(context.Background(), "echo.handler", "hello")
        }
    }

    monitor := func() {
        ticker := time.NewTicker(time.Second)
        tick := 0
        for range ticker.C {
            log.Printf("single rpc qps: %v(%v)", counter, size(counter))
            atomic.StoreUint64(&counter, 0)
            tick++
            if tick > 3 {
                isRunning = false
                ticker.Stop()
                break
            }
        }
    }
    for i:=0; i < rpcSenderCoroutineCount; i++ {
        go sender()
    }
    monitor()

}