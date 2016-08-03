/*

Package inmemstore contains an in-memory Store implementation, safe for concurrent use.

*/
package inmemstore

import (
	"github.com/icza/productws"
	"sync"
)

// inmemStore is an in-memory store implementation.
type inmemStore struct {
	// Map storing products, mapped from their ID
	m map[productws.ID]*productws.Product

	// Mutex to protect concurrent access to the store
	mux sync.RWMutex

	// Id counter to generate new ids
	idCounter productws.ID
}

// NewInmemStore returns a new in-memory Store implementation.
// Safe for concurrent use.
// Also safe against modifying saved or returned products:
// implementation makes necessarying cloning to "detach" saved/returned products
// from the ones in the store.
func NewInmemStore() productws.Store {
	return &inmemStore{m: make(map[productws.ID]*productws.Product), mux: sync.RWMutex{}}
}

// AllIDs implements Store.AllIDs().
// This implementation never returns an error.
func (s *inmemStore) AllIDs() ([]productws.ID, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	ids, count := make([]productws.ID, len(s.m)), 0
	for k := range s.m {
		ids[count] = k
		count++
	}

	return ids, nil
}

// Save implements Store.Save().
// If p.Id is 0, a new Id will be generated and set.
// productws.ErrInvalidId is returned if p.Id is not 0 but no product exists with that ID.
func (s *inmemStore) Save(p *productws.Product) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if p.ID == 0 {
		// Generate id for new product
		s.idCounter++
		p.ID = s.idCounter
	} else {
		// Check if product exists
		if s.m[p.ID] == nil {
			return productws.ErrInvalidId
		}
	}

	s.m[p.ID] = p.Clone() // Clone to be safe!
	return nil
}

// Load implements Store.Load().
// productws.ErrInvalidId is returned if no product exists with the specified ID.
func (s *inmemStore) Load(id productws.ID) (*productws.Product, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	p := s.m[id]
	if p == nil {
		return nil, productws.ErrInvalidId
	}

	return p.Clone(), nil // Clone to be safe!
}
