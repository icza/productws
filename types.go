/*

Types used in the product demo.

*/

package productws

import (
	"errors"
)

// Default currency, must be present in all products.
const DefaultCurrency = "USD"

// Price models a price value.
// The value is stored as an integer number (as should be in all financial applications)
// to avoid rounding errors.
// The Value stores the price multiplied my Multiplier.
// So the real price is the quotient: Value / Multiplier
// Generally Multiplier should be a power of 10, a small value that gives
// integer Value after multiplication.
//
// For example to model the price of 1.99:
//     p := Price{Value: 199, Multiplier: 100}
type Price struct {
	Value      int64 // Value of the price, multiplied by Multiplier
	Multiplier int64 // Multiplier: Value = price * Multiplier
}

// ID is the type of product IDs.
type ID int64

// Product models a product.
type Product struct {
	ID   ID     // Unique product ID
	Name string // Name of the product
	Desc string // Description of the product

	// Optional tags of the product
	Tags []string `json:",omitempty"`

	// Price points, mapped from currency (e.g. "USD")
	Prices map[string]Price
}

// Validate validates a product.
// Checks mandatory fields (Id field is not checked).
// Returns an empty string if product is valid, else an error message.
func (p *Product) Validate() string {
	if p.Name == "" {
		return "Name must be specified!"
	}
	if p.Desc == "" {
		return "Desc must be specified!"
	}
	if len(p.Prices) == 0 {
		return "Prices must be specified!"
	}

	for _, v := range p.Prices {
		if v.Value < 0 {
			return "Price Value must be non-negative!"
		}
		if v.Multiplier < 1 {
			return "Price Multiplier must be positive!"
		}
	}
	// Price for default currency must be present
	if _, ok := p.Prices[DefaultCurrency]; !ok {
		return "Price for \"" + DefaultCurrency + "\" currency must be specified!"
	}

	return ""
}

// Clone clones the product, returns a new product identical to, but independent
// from this one.
func (p *Product) Clone() *Product {
	p2 := new(Product)
	*p2 = *p

	// Clone Tags
	if p.Tags != nil {
		p2.Tags = append([]string(nil), p.Tags...)
	}
	// Clone Prices
	p2.Prices = map[string]Price{}
	for k, v := range p.Prices {
		p2.Prices[k] = v
	}

	return p2
}

// JSONResp is the wrapper for all JSON responses.
type JSONResp struct {
	// Operation (name of the API call)
	Op string

	// Tells if the call was completed successfully
	Success bool

	// Optional error message
	Error string `json:",omitempty"`

	// Optional data
	Data interface{} `json:",omitempty"`
}

// Errors returned by this store implementation
var (
	ErrInvalidId = errors.New("Invalid Product ID")
)

// Store defines the interface for the persistent layer,
// where products are stored.
//
// Implementation may choose where and how the products are stored,
// e.g. it may be an in-memory store, a file store, SQL store etc.
type Store interface {
	// AllIDs returns the list of all product IDs.
	AllIDs() ([]ID, error)

	// Save saves a new Product.
	// If ID of the product is 0, it is saved anew.
	// Else it updates an existing product.
	// ErrInvalidId should be returned if p.Id is not 0 but no product exists with that ID.
	Save(p *Product) error

	// Load loads a Product.
	// ErrInvalidId should be returned if no product exists with the specified ID.
	Load(id ID) (*Product, error)
}
