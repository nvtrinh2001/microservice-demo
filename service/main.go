package main

import (
	"currency/protos/currency"
	"currency/server"
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	grpcServer := grpc.NewServer()

	currencyServer := server.NewCurrencyServer(log, currency.UnimplementedCurrencyServer{})

	currency.RegisterCurrencyServer(grpcServer, currencyServer)

	// add grpc service to the reflection list when the client wants to achieve info of the server
	// e.g. using grpcurl list
	reflection.Register(grpcServer)

	// create a TCP socket for inbound server connections
	listeningSocket, err := net.Listen("tcp", fmt.Sprintf(":%d", 9092))
	if err != nil {
		log.Error("Unable to create a new listener", "error", err)
		os.Exit(1)
	}

	grpcServer.Serve(listeningSocket)
}
