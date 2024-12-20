package db

import "database/sql"

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	// add transactional or other extra functions
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
