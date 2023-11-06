package data

import "fmt"

var ErrProductNotFound = fmt.Errorf("Product not found")

type Product struct {
	ID    int     `json:"id"` // Unique identifier for the product
	Price float32 `json:"price"`
}

func NewProduct(id int, price float32) *Product {
	return &Product{id, price}
}

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func GetProductByID(id int) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	return productList[i], nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

var productList = []*Product{
	&Product{
		ID:    1,
		Price: 2.45,
	},
	&Product{
		ID:    2,
		Price: 1.99,
	},
}
