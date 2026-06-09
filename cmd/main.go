package main

import (
	"github.com/ethereal3x/apc/server"
)

func main() {
	initApp()
	rs := newGrpcServer()
	hs := newGatewayServer()
	server.RunGrpcGatewayService(rs, hs)
}
