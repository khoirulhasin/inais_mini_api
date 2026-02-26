package validates

import (
	"reflect"
	"slices"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

func ValidateInteger(val any, fieldName string) (interface{}, error) {
	typeName := reflect.TypeOf(val).String()
	dataTypes := []string{"int", "int32", "int64"}
	isNumber := slices.Contains(dataTypes, typeName)
	if !isNumber {
		return nil, gqlerror.Errorf("Field '%s' harus berupa integer untuk validasi integer", fieldName)
	}

	return val, nil
}
