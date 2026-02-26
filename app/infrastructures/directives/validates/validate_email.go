package validates

import (
	"regexp"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(val any, fieldName string) (interface{}, error) {
	if !emailRegex.MatchString(val.(string)) {
		return nil, gqlerror.Errorf("Field '%s' harus berupa alamat email yang valid", fieldName)
	}
	return val, nil
}
