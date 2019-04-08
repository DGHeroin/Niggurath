package client

import (
	"context"
	"github.com/DGHeroin/Niggurath/rpc/core"
	"log"
	"sync"
	"time"
)

type worker struct {
	client        Client
	ch            *chan core.Request
	isRunning     bool
	wg            *sync.WaitGroup
	remoteAddress string
	retryDelay    int
	connected bool
}

func (w *worker) doConnect() error{
	if err := w.client.Connect(w.remoteAddress); err != nil {
		log.Printf("rpc connect error: %v", err)
		w.retryDelay = 1
		w.connected = false
		return err
	} else {
		w.connected = true
		return nil
	}
}

func (w *worker) start(remote string) {
	w.remoteAddress = remote

	w.isRunning = true
	w.connected = false
	for {
		if !w.isRunning {
			break
		}
		// try connect
		if w.connected == false {
			if err := w.doConnect(); err != nil {
				w.retryDelay++
				if w.retryDelay > 10 {
					w.retryDelay = 10 // 500ms * 120, max delay 60s
				}
				delay := time.Millisecond * 500 * time.Duration(w.retryDelay)
				log.Printf("delay reconnect: %v", delay)
				time.Sleep(delay)
				continue
			}
		}
		// handle message
		timeout := time.After(time.Millisecond * 100) // 100ms timeout
		select {
		case msg := <-*w.ch:
			w.client.ExecuteSingleEvent(context.Background(), msg.Name, msg.Message)
		case <-timeout: // wakeup but do nothing
			break
		}
	}
	w.wg.Done()
}

type Pool struct {
	ch      chan core.Request
	workers []*worker
	isInit  bool
	wg      sync.WaitGroup
}

func (p *Pool) NewWorker(num int, remote string) {
	if p.isInit {
		return
	}
	p.ch = make(chan core.Request, num)
	for i := 0; i < num; i++ {
		p.wg.Add(1)
		w := worker{ch: &p.ch, wg: &p.wg}
		p.workers = append(p.workers, &w)
		go w.start(remote)
	}
	p.isInit = true
}

func (p *Pool) Send(name, message string) {
	p.ch <- core.Request{Name: name, Message: message}
}
func (p *Pool) Close() {
	for _, w := range p.workers {
		w.isRunning = false
	}
	p.wg.Wait()
	log.Println("close all workers")
}
