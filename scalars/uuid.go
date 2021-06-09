package scalars

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gobuffalo/uuid"
)

type UUID struct {
	uuid.UUID
}

func MarshalUUID(id UUID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, "\""+id.String()+"\"")
	})
}

func UnmarshalUUID(v interface{}) (UUID, error) {
	id := UUID{}
	s, ok := v.(string)
	if !ok {
		return id, fmt.Errorf("points must be strings")
	}

	uid, err := uuid.FromString(s)

	if err != nil {
		return id, fmt.Errorf("invalid uuid")
	}

	if err := id.UnmarshalBinary(uid.Bytes()); err != nil {
		return id, fmt.Errorf("invalid uuid unmarshaling")
	}
	return id, nil
}

func UUIDNil() UUID {
	id := UUID{}
	id.UnmarshalBinary(uuid.Nil.Bytes())
	return id
}

func ModelsUUID2GhqlUUID(id uuid.UUID) UUID {
	return UUID{
		id,
	}
}

func GhqlUUID2ModelsUUID(id UUID) uuid.UUID {
	return uuid.FromBytesOrNil(id.Bytes())
}
