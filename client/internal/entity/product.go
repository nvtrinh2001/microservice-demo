package entity

import "fmt"

var ErrProductNotFound = fmt.Errorf("Product not found")

type Product struct {
	ID    int     `json:"id"` // Unique identifier for the product
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

