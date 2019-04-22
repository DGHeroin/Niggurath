package client

import (
	"context"
	"log"
	"sync"
	"time"
)

type worker struct {
	client        Client
	ch            *chan *workerRequest
	isRunning     bool
	wg            *sync.WaitGroup
	remoteAddress string
	retryDelay    int
	connected     bool
}

func (w *worker) doConnect() error {
	if err := w.client.Connect(w.remoteAddress); err != nil {
		w.connected = false
		return err
	} else {
		w.retryDelay = 0
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
				delay := time.Second * time.Duration(w.retryDelay)
				time.Sleep(delay)
				continue
			}
		}
		// handle message
		timeout := time.After(time.Millisecond * 100) // 100ms timeout
		select {
		case req := <-*w.ch:
			req.response, req.err = w.client.ExecuteSingleEvent(context.Background(), req.name, req.message)
			req.wg.Done()
		case <-timeout: // wakeup but do nothing
			break
		}
	}
	w.wg.Done()
}

type workerRequest struct {
	name     string
	message  string
	response string
	err      error
	wg sync.WaitGroup
}

type Pool struct {
	ch      chan *workerRequest
	workers []*worker
	isInit  bool
	wg      sync.WaitGroup
}

func (p *Pool) NewWorker(num int, remote string) {
	if p.isInit {
		return
	}
	p.ch = make(chan *workerRequest, num)
	for i := 0; i < num; i++ {
		p.wg.Add(1)
		w := worker{ch: &p.ch, wg: &p.wg}
		p.workers = append(p.workers, &w)
		go w.start(remote)
	}
	p.isInit = true
}

func (p *Pool) Send(name, message string) (result string, err error) {
	var (
		request = &workerRequest{name:name, message:message, err:nil}
	)
	request.wg.Add(1)
	p.ch <- request
	request.wg.Wait()
	result = request.response
	err = request.err
	return
}
func (p *Pool) Close() {
	for _, w := range p.workers {
		w.isRunning = false
	}
	p.wg.Wait()
	log.Println("close all workers")
}
