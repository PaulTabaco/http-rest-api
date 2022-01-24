package sqlstore

import (
	"database/sql"
	"paulTabaco/http-rest-api/internal/app/model"
	"paulTabaco/http-rest-api/internal/app/store"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Create ...
func (r *UserRepository) Create(u *model.User) error {

	if err := u.Validate(); err != nil {
		return err
	}

	// Deal with hashing password
	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id", u.Email, u.EncryptedPassword,
	).Scan(&u.ID) // Fill ID to passed arg u
}

// FindById ...
func (r *UserRepository) FindById(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE id = $1", id,
	).Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil { // Scan to fill returned u with id, email, EncryptedPassword
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE email = $1", email,
	).Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil { // Scan to fill returned u with id, email, EncryptedPassword
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}
