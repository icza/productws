/*

Package main is the main package of the product web service demo application.

It starts the web service with an in-memory store implementation.

Test products are inserted by default (can be disabled with -testdata=false).

Also imports html-tester, so the tester page will be self-contained and made available under
    /tester.html

*/
package main

import (
	"flag"
	"github.com/icza/productws"
	_ "github.com/icza/productws/html-tester"
	"github.com/icza/productws/inmemstore"
	"log"
	"net/http"
)

// Command line flags
var (
	addr     = flag.String("addr", ":8081", "address to start server on (host:port)")
	testData = flag.Bool("testdata", true, "tells if test data should be inserted on startup")
)

func main() {
	flag.Parse()

	store := inmemstore.NewInmemStore()
	productws.SetStore(store)

	if *testData {
		insertTestData(store)
	}

	log.Printf("Starting server on %q...", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// insertTestData inserts test products into the store.
func insertTestData(store productws.Store) {
	ps := []*productws.Product{
		&productws.Product{Name: "small-prod", Desc: "short-desc",
			Prices: map[string]productws.Price{"USD": {Value: 1, Multiplier: 1}}},
		&productws.Product{Name: "Full-prod", Desc: "long description is entered here",
			Tags: []string{"Big", "Full", "Giant"},
			Prices: map[string]productws.Price{
				"USD": {Value: 100, Multiplier: 1},
				"GBP": {Value: 7528, Multiplier: 100},
				"HUF": {Value: 27725, Multiplier: 1},
			}},
	}

	for _, p := range ps {
		if err := store.Save(p); err != nil {
			log.Printf("Failed to insert test product ID=%d: %v", p.ID, err)
		} else {
			log.Printf("Test product inserted (ID=%d)", p.ID)
		}
	}
}
