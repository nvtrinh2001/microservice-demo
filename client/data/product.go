package data

import "fmt"

var ErrProductNotFound = fmt.Errorf("Product not found")

type Product struct {
	ID    int     `json:"id"` // Unique identifier for the product
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

func NewProduct(id int, name string, price float32) *Product {
	return &Product{id, name, price}
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
