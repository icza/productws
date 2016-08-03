/*

This is the main package of the product web service demo application.

It starts the web service with an in-memory store implementation.

*/
package main

import (
	"flag"
	"github.com/icza/productws"
	"github.com/icza/productws/inmemstore"
	"log"
	"net/http"
)

var (
	addr     string // Address to start the server on (host:port)
	testData bool   // Tells if test data should be inserted on startup
)

func init() {
	flag.StringVar(&addr, "addr", ":8081", "address to start server on (host:port)")
	flag.BoolVar(&testData, "testdata", true, "tells if test data should be inserted on startup")
}

func main() {
	flag.Parse()

	store := inmemstore.NewInmemStore()
	productws.SetStore(store)

	if testData {
		// Insert test products:
		store.Save(&productws.Product{Name: "small-prod", Desc: "short-desc",
			Prices: map[string]productws.Price{"USD": {Value: 1, Multiplier: 1}}})
		store.Save(&productws.Product{Name: "Full-prod", Desc: "long description is entered here",
			Tags: []string{"Big", "Full", "Giant"},
			Prices: map[string]productws.Price{
				"USD": {Value: 100, Multiplier: 1},
				"GBP": {Value: 7528, Multiplier: 100},
				"HUF": {Value: 27725, Multiplier: 1},
			}})
	}

	log.Printf(`Starting server on %q...`, addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
