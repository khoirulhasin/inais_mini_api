package validates

import (
	"regexp"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)

func ValidateUsername(val any, fieldName string) (interface{}, error) {
	if !usernameRegex.MatchString(val.(string)) {
		return nil, gqlerror.Errorf("Field '%s' harus berupa username yang valid (3-20 karakter, hanya huruf, angka, '_', atau '-')", fieldName)
	}
	return val, nil
}
