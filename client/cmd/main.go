package main

import (
	"context"
	"currency/client/internal/business"
	"currency/client/internal/entity"
	"currency/client/internal/repository/rpc"
	"currency/client/internal/repository/storage"
	"currency/client/internal/transport/api"
	"currency/client/protos/currency"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"google.golang.org/grpc"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9091", "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")

var productList = []*entity.Product{
	{
		ID:    1,
		Name:  "espresso",
		Price: 2.45,
	},
	{
		ID:    2,
		Name:  "americano",
		Price: 1.99,
	},
}

func main() {
	env.Parse()

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "client",
			Level: hclog.LevelFromString(*logLevel),
		},
	)

	// create a logger for the server from the default logger
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	// storage
	memStore := storage.NewMemStore(productList)

	// create client
	cc := currency.NewCurrencyClient(conn)
	rpcClient := rpc.NewRPCClient(cc)

	// business
	business := business.NewBusiness(memStore, rpcClient, time.Duration(2)*time.Second, l)

	// create the handlers
	apiHandler := api.NewAPI(business, l)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/products/{id:[0-9]+}", apiHandler.GetProduct).Queries("base", "{[A-Z]{3}}", "dest", "{[A-Z]{3}}")
	getR.HandleFunc("/products/{id:[0-9]+}", apiHandler.GetProduct)

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create a new server
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      ch(sm),            // set the default handler
		ErrorLog:     sl,                // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Info("Starting server", "bind_address", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	l.Info("Shutting down server with", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.Shutdown(ctx)
}
