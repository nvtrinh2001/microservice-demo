package handlers

import (
	"currency/client/protos/currency"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type GenericError struct {
	Message string
}

// Products handler for getting and updating products
type Products struct {
	l  hclog.Logger
	cc currency.CurrencyClient
}

// NewProducts returns a new products handler with the given logger
func NewProducts(l hclog.Logger, cc currency.CurrencyClient) *Products {
	return &Products{l, cc}
}

// getProductID returns the product ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getProductID(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}
