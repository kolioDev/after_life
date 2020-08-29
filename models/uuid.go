package models

import (
	"fmt"
	"github.com/gobuffalo/uuid"
	"io"
)

type UUID struct {
	uuid.UUID
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (id *UUID) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("points must be strings")
	}

	uid, err := uuid.FromString(s)

	if err != nil {
		return fmt.Errorf("invalid uuid")
	}

	if err := id.UnmarshalBinary(uid.Bytes()); err != nil {
		return fmt.Errorf("invalid uuid unmarshaling")
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (id UUID) MarshalGQL(w io.Writer) {
	w.Write(id.Bytes())
}

func UUIDNil() UUID {
	id := UUID{}
	id.UnmarshalBinary(uuid.Nil.Bytes())
	return id
}
