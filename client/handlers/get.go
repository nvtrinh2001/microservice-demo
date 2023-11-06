package handlers

import (
	"context"
	"currency/client/data"
	"currency/client/protos/currency"
	"net/http"
)

// ListSingle handles GET requests
func (p *Products) GetProduct(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)

	prod, err := data.GetProductByID(id)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Error("Error fetching product", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Error("Error fetching product", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	// get exchange rate
	rr := &currency.RateRequest{
		Base:        currency.Currencies(currency.Currencies_value["EUR"]),
		Destination: currency.Currencies(currency.Currencies_value["GBP"]),
	}

	resp, err := p.cc.GetRate(context.Background(), rr)
	if err != nil {
		p.l.Error("Error getting new rate", "error", err)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	p.l.Info("Response from currency service", "info", resp)

	prod.Price = prod.Price * resp.Rate

	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Error serializing product", "error", err)
	}
}
