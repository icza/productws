/*

Contains the handlers implementing the product demo REST API operations.

*/

package productws

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Constants for the operations (names of API calls)
const (
	opCreate    = "create"    // Create a new product
	opList      = "list"      // Getting a list of all products
	opDetails   = "details"   // Getting details of a product.
	opUpdate    = "update"    // Update a product
	opSetPrices = "setprices" // Set price points for different currencies for a product
)

// Store implementation to use
var store Store

// SetStore sets the Store used by the API calls.
// Must be done prior to starting the web service.
func SetStore(st Store) {
	store = st
}

// General messages sent in response
const (
	MsgGeneralStoreErr = "Product store unavailable" // General error message concerning Store errors.
	MsgInvalidIDErr    = "Invalid ID!"               // Error saying no product for the ID
)

// createUpdateLogic implements creating a new product and updating a product.
// Expects the request body to be the JSON product to be created or updated.
// Creation requires ID not be present, update requires a valid ID (to be updated).
func createUpdateLogic(w http.ResponseWriter, r *http.Request, ch *callHandler) *JSONResp {
	p := new(Product)
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		log.Printf("Error decoding %s request: %v", ch.op, err)
		http.Error(w, "Can't decode input JSON", http.StatusBadRequest)
		return nil
	}

	// Id must not be specified when creating a new product:
	if ch.op == opCreate && p.ID != 0 {
		return &JSONResp{Error: "ID must not be specified!"}
	}
	// Id must be specified when updating an existing product:
	if ch.op == opUpdate && p.ID == 0 {
		return &JSONResp{Error: "ID must be specified!"}
	}
	if msg := p.Validate(); msg != "" {
		return &JSONResp{Error: msg}
	}

	if err := store.Save(p); err != nil {
		log.Printf("Error saving product: %v", err)
		if err == ErrInvalidId {
			return &JSONResp{Error: MsgInvalidIDErr}
		}
		return &JSONResp{Error: MsgGeneralStoreErr}
	}

	return &JSONResp{Success: true, Data: struct{ ID ID }{p.ID}}
}

// listLogic implements getting a list of all products.
// Does not require anything in the request path or body.
func listLogic(http.ResponseWriter, *http.Request, *callHandler) *JSONResp {
	ids, err := store.AllIDs()
	if err != nil {
		log.Printf("Error getting all product IDs: %v", err)
		return &JSONResp{Error: MsgGeneralStoreErr}
	}

	return &JSONResp{Success: true, Data: ids}
}

// detailsLogic implements getting details about a product.
// Requires the path to contain the ID of the product whose details to return.
// Path must be like
//     /details/id
func detailsLogic(w http.ResponseWriter, r *http.Request, ch *callHandler) *JSONResp {
	// Get id from path which is like /details/id
	parts := strings.Split(r.URL.Path, "/")
	var id ID
	if len(parts) >= 3 {
		if id_, err := strconv.ParseInt(parts[2], 10, 64); err == nil {
			id = ID(id_)
		}
	}
	if id == 0 {
		log.Printf("Invalid path: %v", r.URL.Path)
		http.Error(w, "Path must be like /details/id", http.StatusBadRequest)
		return nil
	}

	var p *Product
	var err error
	if p, err = store.Load(id); err != nil {
		log.Printf("Error loading product with id %d: %v", id, err)
		if err == ErrInvalidId {
			return &JSONResp{Error: MsgInvalidIDErr}
		}
		return &JSONResp{Error: MsgGeneralStoreErr}
	}

	return &JSONResp{Success: true, Data: p}
}

// setPricesLogic implements setting price points for different currencies for a product.
// Requires the body to be a JSON product, but only the ID and Prices fields should be present
// (other fields are omitted).
func setPricesLogic(w http.ResponseWriter, r *http.Request, ch *callHandler) *JSONResp {
	p := new(Product)
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		log.Printf("Error decoding %s request: %v", ch.op, err)
		http.Error(w, "Can't decode input JSON", http.StatusBadRequest)
		return nil
	}

	if p.ID == 0 {
		return &JSONResp{Error: "ID must be specified!"}
	}
	// Validate prices
	if len(p.Prices) == 0 {
		return &JSONResp{Error: "Prices must be specified!"}
	}
	for _, v := range p.Prices {
		if v.Value < 0 {
			return &JSONResp{Error: "Price Value must be non-negative!"}
		}
		if v.Multiplier < 1 {
			return &JSONResp{Error: "Price Multiplier must be positive!"}
		}
	}

	// First get existing product
	var p2 *Product
	var err error
	if p2, err = store.Load(p.ID); err != nil {
		log.Printf("Error loading product with id %d: %v", p.ID, err)
		if err == ErrInvalidId {
			return &JSONResp{Error: MsgInvalidIDErr}
		}
		return &JSONResp{Error: MsgGeneralStoreErr}
	}

	// Merge changes into the product:
	for k, v := range p.Prices {
		p2.Prices[k] = v
	}

	// And finally save updated product
	if err := store.Save(p2); err != nil {
		log.Printf("Error saving product: %v", err)
		if err == ErrInvalidId {
			return &JSONResp{Error: MsgInvalidIDErr}
		}
		return &JSONResp{Error: MsgGeneralStoreErr}
	}

	return &JSONResp{Success: true, Data: struct{ ID ID }{p2.ID}}
}

// callLogic is a function type of call logic implementations.
type callLogic func(http.ResponseWriter, *http.Request, *callHandler) *JSONResp

// callHandler is an API call handler.
// Contains certain properties and characteristics of API calls.
type callHandler struct {
	op        string    // Operation (name of the API call)
	expMethod string    // Expected HTTP method for the api call
	logic     callLogic // Call handling logic
}

// ServeHTTP implements http.Handler.
// Contains common logic for all api calls, and invokes the logic handler.
// Common logic includes checking expected HTTP method, calling the logic,
// marshaling JSON response. Also future authentication can be added here.
func (ch *callHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Allow JavaScript to access API calls:
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT")
	if r.Method == http.MethodOptions {
		return
	}

	// Disable caching:
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // For HTTP 1.1
	w.Header().Set("Pragma", "no-cache")                                   // For HTTP 1.0
	w.Header().Set("Expires", "0")                                         // For proxies

	// If authentication is required, it can be checked here.
	if r.Method != ch.expMethod {
		http.Error(w, "Method not allowed, use "+ch.expMethod, http.StatusMethodNotAllowed)
		return
	}

	if jsonResp := ch.logic(w, r, ch); jsonResp != nil {
		// Send JSON response
		jsonResp.Op = ch.op
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(jsonResp); err != nil {
			log.Printf("Failed to send JSON response: %v", err)
		}
	}
}

// init registers the HTTP handlers.
func init() {
	http.Handle("/"+opCreate, &callHandler{op: opCreate, expMethod: http.MethodPost, logic: createUpdateLogic})
	http.Handle("/"+opList, &callHandler{op: opList, expMethod: http.MethodGet, logic: listLogic})
	http.Handle("/"+opDetails+"/", &callHandler{op: opDetails, expMethod: http.MethodGet, logic: detailsLogic})
	http.Handle("/"+opUpdate, &callHandler{op: opUpdate, expMethod: http.MethodPut, logic: createUpdateLogic})
	http.Handle("/"+opSetPrices, &callHandler{op: opSetPrices, expMethod: http.MethodPut, logic: setPricesLogic})
}
