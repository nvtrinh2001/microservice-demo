package business

import (
	"context"
	"currency/client/internal/entity"
	"currency/client/protos/currency"
	"time"

	"github.com/hashicorp/go-hclog"
)

type GetProductRequest struct {
	ID          int                 `json:"id"`
	Base        currency.Currencies `json:"base"`
	Destination currency.Currencies `json:"destination"`
}

type GetProductResponse struct {
	ID    int     `json:"id"` // Unique identifier for the product
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

type Business interface {
	GetProduct(ctx context.Context, request *GetProductRequest) (*GetProductResponse, error)
}

type ProductRepository interface {
	GetProductByID(ctx context.Context, id int) (*entity.Product, error)
}

type RateExchangeService interface {
	GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error)
}

type business struct {
	productRepository   ProductRepository
	rateExchangeService RateExchangeService
	timeout             time.Duration
	l                   hclog.Logger
}

func NewBusiness(productRepository ProductRepository, rateExchangeService RateExchangeService, timeout time.Duration, l hclog.Logger) *business {
	return &business{
		productRepository,
		rateExchangeService,
		timeout,
		l,
	}
}

func (b *business) GetProduct(c context.Context, request *GetProductRequest) (*GetProductResponse, error) {
	ctx, cancel := context.WithTimeout(c, b.timeout)
	defer cancel()

	prod, err := b.productRepository.GetProductByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	// get exchange rate
	rr := &currency.RateRequest{
		Base:        request.Base,
		Destination: request.Destination,
	}

	resp, err := b.rateExchangeService.GetRate(context.Background(), rr)
	if err != nil {
		return nil, err
	}

	prod.Price = prod.Price * resp.Rate

	b.l.Info("Recalculate price successfully", "price", prod.Price)

	return &GetProductResponse{prod.ID, prod.Name, prod.Price}, nil
}
