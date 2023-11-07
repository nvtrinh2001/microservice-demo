package api

import (
	"currency/client/internal/business"
	"currency/client/internal/entity"
	"currency/client/pkg"
	"currency/client/protos/currency"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type api struct {
	business.Business
	l hclog.Logger
}

type GenericError struct {
	Message string
}

func NewAPI(b business.Business, l hclog.Logger) *api {
	return &api{b, l}
}

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

func (api *api) GetProduct(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)
	base := r.URL.Query().Get("base")
	dest := r.URL.Query().Get("dest")

	prod, err := api.Business.GetProduct(
		r.Context(),
		&business.GetProductRequest{
			ID:          id,
			Base:        currency.Currencies(currency.Currencies_value[base]),
			Destination: currency.Currencies(currency.Currencies_value[dest]),
		})

	switch err {
	case nil:

	case entity.ErrProductNotFound:
		api.l.Error("Unable to fetch product", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		api.l.Error("Unable to fetching product", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = pkg.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		api.l.Error("Unable to serializing product", err)
	}

	api.l.Info("Successfully get the product")
}
