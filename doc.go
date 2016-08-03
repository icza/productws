/*

Package productws is a REST / JSON web service demo with an API to manage products.


Implementation overview

API calls are described by the callHandler type. Each API call has a value of this,
and it is registered as the handler for its associated path.
API calls include a logic function which are called by the callHandler.

The persistent layer is abstracted by the Store interface. API calls only use this interface,
so the Store implementation is completely swappable.
The Store must be set with the SetStore() function prior to starting the web service.
Package inmemstore contains an in-memory Store implementation.

Price points (prices) are modeled with a map, mapping from currency to the price value.
I chose this representation as it implicitly takes care of currencies being unique,
and also provides easy access to prices by currencies.

Price values are modeled by the Price type. The value is stored as an integer number
(as should be in all financial applications) to avoid rounding errors.
The Price.Value stores the price multiplied my Price.Multiplier.
So the real price is the quotient: Value / Multiplier.
Generally Multiplier should be a power of 10, a small value that gives integer Value
after multiplication.


API calls

The create API call creates a new product. It must be a POST request,
and it expects the request body to be the JSON product to be created.
Creation requires ID not be present.

The update API call updates an existing product. It must be a PUT request,
and it expects the request body to be the JSON product to be created.
Update requires a valid ID (to be updated).

The list API call list all existing product IDs. It must be a GET request,
and it does not require anything in the request path or body.

The details API call returns all the details of a product. It must be a GET request,
and it expects the path to contain the ID of the product whose details to return.

The setprices API call sets price points for different currencies of a product.
It must be a PUT request, and it expects the body to be a JSON product,
but only the ID and Prices fields should be present (other fields are omitted).
It merges the specified prices with the existing prices. That is, if a product
already has prices in USD and GBP currencies, and the call specifies prices
in GBP and HUF currencies, then the GBP price will be updated, HUF added
and USD left intact. If currency removal is required, update call can (should) be used.

*/

package productws
