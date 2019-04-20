package main

import (
    "flag"
    "github.com/DGHeroin/Niggurath/event"
    "github.com/DGHeroin/Niggurath/rpc/client"
    "github.com/DGHeroin/Niggurath/rpc/server"
    "log"
    "sync/atomic"
    "time"
)

var (
    pool client.Pool
)

func startTestServer() {
    log.Printf("listen on: %v", serverAddress)
    srv := server.Server{}
    srv.Serve(serverAddress)
}

type Handler struct{}

func (h *Handler) OnMessage(message string) (result string, err error) {
    atomic.AddUint64(&procQPS, 1)
    return
}

var (
    isRunning               = true
    procQPS                 = uint64(0)
    sendQPS                 = uint64(0)
    isServer                = false
    isClient                = false
    serverAddress           = ":7878"
    remoteAddress           = "127.0.0.1:7878"
    runSeconds              = 5
    workerNum               = 10
)

func runServer() {
    event.GetEventManager().AddSingleListener("h", &Handler{})
    go startTestServer()
}

func runClient() {
    pool.NewWorker(workerNum, remoteAddress)
    sender := func() {
        for {
            if !isRunning {
                break
            }
            pool.Send("h", "")
            atomic.AddUint64(&sendQPS, 1)
        }
    }

    for i := 0; i < workerNum; i++ {
        go sender()
    }
}

func monitor() {
    ticker := time.NewTicker(time.Second)
    tick := 0
    for range ticker.C {
        log.Printf("process qps: %v send qps:%v", procQPS, sendQPS)
        atomic.StoreUint64(&procQPS, 0)
        atomic.StoreUint64(&sendQPS, 0)
        tick++
        if runSeconds == 0 { // run forever
        } else {
            if tick > runSeconds {
                isRunning = false
                ticker.Stop()
                break
            }
        }

    }
}

func main() {

    flag.BoolVar(&isServer, "s", false, "run as server")
    flag.BoolVar(&isClient, "c", false, "run as client")
    flag.StringVar(&serverAddress, "l", ":7878", "server listen address")
    flag.StringVar(&remoteAddress, "r", "127.0.0.1:7878", "connect to remote address")
    flag.IntVar(&runSeconds, "t", 5, "count time")
    flag.IntVar(&workerNum, "workerNum", 10, "number of workers")
    flag.Parse()

    if isServer {
        log.Println("run as server")
        go runServer()
        monitor()
        return
    }

    if isClient {
        log.Println("run as client")
        go runClient()
        defer pool.Close()
        monitor()
        return
    }

    log.Println("run as sever and client")
    go runServer()
    go runClient()
    defer pool.Close()

    monitor()
}
