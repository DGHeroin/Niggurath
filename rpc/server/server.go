package server

import (
	"github.com/DGHeroin/Niggurath/rpc/core"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server struct {
	listener net.Listener
}

func (s *Server) Serve(address string) (err error) {
	rpc.Register(&core.Handler{})

	s.listener, err = net.Listen("tcp", address)
	if err != nil {
		log.Println("start falied:", err)
		return err
	}
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go jsonrpc.ServeConn(conn)
	}
}
