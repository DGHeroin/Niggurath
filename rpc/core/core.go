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

type Request struct {
    Name    string
    Message string
}

const HandlerName = "Handler.Execute"

type Handler struct {}

func (h *Handler) Execute(req Request, res *Response) (err error) {
    if req.Name == "" {
        err = fmt.Errorf("A name must be specified")
        return
    }

    resp, err := event.GetEventManager().DispatchEvent(req.Name, req.Message)
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