package validates

import (
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func ValidateRequired(val any, fieldName string) (interface{}, error) {

	switch val := val.(type) {
	case string:
		if val == "" {
			return nil, gqlerror.Errorf("Field '%s' wajib diisi (tidak boleh kosong)", fieldName)
		}
	case *string:
		if val == nil || *val == "" {
			return nil, gqlerror.Errorf("Field '%s' wajib diisi (tidak boleh kosong)", fieldName)
		}
	case int:
		if val == 0 {
			return nil, gqlerror.Errorf("Field '%s' wajib diisi (tidak boleh 0)", fieldName)
		}
	case *int:
		if val == nil || *val == 0 {
			return nil, gqlerror.Errorf("Field '%s' wajib diisi (tidak boleh 0)", fieldName)
		}
	case int32:
		if val == 0 {
			return nil, gqlerror.Errorf("Field '%s' wajib diisi (tidak boleh 0)", fieldName)
		}
	case *int32:
		if val == nil || *val == 0 {
			return nil, gqlerror.Errorf("Field '%s' wajib diisi (tidak boleh 0)", fieldName)
		}
	case int64:
		if val == 0 {
			return nil, gqlerror.Errorf("Field '%s' wajib diisi (tidak boleh 0)", fieldName)
		}
	case *int64:
		if val == nil || *val == 0 {
			return nil, gqlerror.Errorf("Field '%s' wajib diisi (tidak boleh 0)", fieldName)
		}
	case float64:
		if val == 0.0 {
			return nil, gqlerror.Errorf("Field '%s' wajib diisi (tidak boleh 0.0)", fieldName)
		}
	case *float64:
		if val == nil || *val == 0.0 {
			return nil, gqlerror.Errorf("Field '%s' wajib diisi (tidak boleh 0.0)", fieldName)
		}
	case []interface{}:
		if len(val) == 0 {
			return nil, gqlerror.Errorf("Field '%s' wajib diisi (daftar tidak boleh kosong)", fieldName)
		}
	case *[]interface{}:
		if val == nil || len(*val) == 0 {
			return nil, gqlerror.Errorf("Field '%s' wajib diisi (daftar tidak boleh kosong)", fieldName)
		}
	default:
		if val == nil {
			return nil, gqlerror.Errorf("Field '%s' wajib diisi (tidak boleh null)", fieldName)
		}
	}

	return val, nil
}
