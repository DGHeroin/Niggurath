package client

import (
    "github.com/DGHeroin/Niggurath/rpc/core"
    "net/rpc"
    "net/rpc/jsonrpc"
    "log"
    "context"
)

type Client struct {
    client *rpc.Client
    remote string
    ok bool
}

func (c *Client) Connect(remote string) (err error) {
    c.remote = remote
    c.client, err = jsonrpc.Dial("tcp", remote)
    if err != nil {
        return err
    }
    c.ok = true
    return nil
}

func (c *Client) Clone() (*Client, error) {
    cli := &Client{}
    if err := cli.Connect(c.remote); err == nil {
        return cli, nil
    } else {
        return nil, err
    }
}

func (c *Client) Close() (err error) {
    return c.client.Close()
}

func (c *Client) Execute(ctx context.Context, name string, message string) (msg []string, err error) {
    if !c.ok {
        err = c.Connect(c.remote)
        if err != nil {
            log.Println("Try auto connect failed")
            return
        }
    }

    var (
        request = &core.Request{Name:name,Message:message}
        response = new(core.Response)
    )

    err = c.client.Call(core.HandlerName, request, response)
    if err != nil {
        c.ok = false
        return
    }
    c.ok = true
    msg = response.Message
    return
}

func (c *Client) ExecuteSingleEvent(ctx context.Context, name string, message string) (msg string, err error) {
    if !c.ok {
        err = c.Connect(c.remote)
        if err != nil {
            log.Println("Try auto connect failed")
            return
        }
    }

    var (
        request = &core.Request{Name:name,Message:message}
        response = new(core.ResponseSingleEvent)
    )

    err = c.client.Call(core.SingleHandlerName, request, response)
    if err != nil {
        c.ok = false
        return
    }
    c.ok = true
    msg = response.Message
    return
}
