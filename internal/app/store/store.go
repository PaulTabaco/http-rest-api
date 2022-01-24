package store

// Store ...
type Store interface {
	UserRep() UserRepository
}
