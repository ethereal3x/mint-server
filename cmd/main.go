package main

import (
	"log"

	"github.com/ethereal3x/apc/server"
	"github.com/ethereal3x/mint-server/internal/app"
)

func main() {
	application, err := app.New("")
	if err != nil {
		log.Fatalf("init app: %v", err)
	}
	rs := newGrpcServer(application)
	hs := newGatewayServer(application)
	server.RunGrpcGatewayService(rs, hs)
}
