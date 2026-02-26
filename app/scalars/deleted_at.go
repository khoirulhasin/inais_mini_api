package scalars

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"gorm.io/plugin/soft_delete"
)

// MarshalDeletedAt mendefinisikan cara serialisasi soft_delete.DeletedAt ke GraphQL
func MarshalDeletedAt(d soft_delete.DeletedAt) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		// Konversi soft_delete.DeletedAt ke int64 dan kemudian ke string
		_, _ = io.WriteString(w, strconv.FormatInt(int64(d), 10))
	})
}

// UnmarshalDeletedAt mendefinisikan cara deserialisasi dari GraphQL ke soft_delete.DeletedAt
func UnmarshalDeletedAt(v interface{}) (soft_delete.DeletedAt, error) {
	switch v := v.(type) {
	case int:
		return soft_delete.DeletedAt(v), nil
	case int64:
		return soft_delete.DeletedAt(v), nil
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid DeletedAt value: %v", err)
		}
		return soft_delete.DeletedAt(i), nil
	default:
		return 0, fmt.Errorf("invalid type for DeletedAt: %T", v)
	}
}
