package core

import (
    "fmt"
    "github.com/DGHeroin/Niggurath/event"
    "log"
)

type Response struct {
    Message []string
    Ok bool
}

type ResponseSingleEvent struct {
    Message string
    Ok bool
}

type Request struct {
    Name    string
    Message string
}

const HandlerName = "Handler.Execute"
const SingleHandlerName = "Handler.ExecuteSingleEvent"

type Handler struct {}

func (h *Handler) Execute(req Request, res *Response) (err error) {
    if req.Name == "" {
        err = fmt.Errorf("A name must be specified")
        return
    }

    resp, err := event.GetEventManager().DispatchEvents(req.Name, req.Message)
    errMsg := ""
    if err != nil {errMsg = err.Error()}
    if err == nil {
        res.Ok = true
    } else {
        log.Printf("Execute Filed: [%v] err:[%v]", req.Name, errMsg)
        res.Ok = false
    }
    res.Message = resp
    return
}

func (h *Handler) ExecuteSingleEvent(req Request, res *ResponseSingleEvent) (err error) {
    if req.Name == "" {
        err = fmt.Errorf("A name must be specified")
        return
    }

    resp, err := event.GetEventManager().DispatchSingleEvent(req.Name, req.Message)
    errMsg := ""
    if err != nil {errMsg = err.Error()}
    if err == nil {
        res.Ok = true
    } else {
        log.Printf("Execute Filed: [%v] err:[%v]", req.Name, errMsg)
        res.Ok = false
    }
    res.Message = resp
    return
}