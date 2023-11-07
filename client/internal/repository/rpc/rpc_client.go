package rpc

import (
	"context"
	"currency/client/protos/currency"
)

type rpcClient struct {
	currency.CurrencyClient
}

func NewRPCClient(cc currency.CurrencyClient) *rpcClient {
	return &rpcClient{cc}
}

func (rc *rpcClient) GetRate(c context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	return rc.CurrencyClient.GetRate(context.Background(), rr)
}
