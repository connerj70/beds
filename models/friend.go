package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Friend is used by pop to map your friends database table to your go code.
type Friend struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Receiver  uuid.UUID `json:"receiver_id" db:"receiver_id"`
	Requester uuid.UUID `json:"requester_id" db:"requester_id"`
	Accepted  bool      `json:"accepted" db:"accepted"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (f Friend) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Friends is not required by pop and may be deleted
type Friends []Friend

// String is not required by pop and may be deleted
func (f Friends) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (f *Friend) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (f *Friend) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (f *Friend) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
