package server

import (
	"context"

	"currency/protos/currency"

	"github.com/hashicorp/go-hclog"
)

type CurrencyServer struct {
	log hclog.Logger
	currency.UnimplementedCurrencyServer
}

func NewCurrencyServer(l hclog.Logger, unimplementedCurrencyServer currency.UnimplementedCurrencyServer) *CurrencyServer {
	return &CurrencyServer{l, unimplementedCurrencyServer}
}

func (c *CurrencyServer) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", rr.GetBase(), "dest", rr.GetDestination())
	return &currency.RateResponse{
		Rate: 0.5,
	}, nil
}
