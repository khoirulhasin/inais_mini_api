package geofences

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/plugin/soft_delete"
)

// GeoJSONCoords — custom type untuk JSONB, tidak bisa di-generate gqlgen
type GeoJSONCoords []interface{}

func (g GeoJSONCoords) Value() (driver.Value, error) {
	b, err := json.Marshal(g)
	return string(b), err
}

func (g *GeoJSONCoords) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, g)
	case string:
		return json.Unmarshal([]byte(v), g)
	default:
		return fmt.Errorf("GeoJSONCoords: unsupported type %T", value)
	}
}

// pastikan soft_delete diimport
var _ soft_delete.DeletedAt = soft_delete.DeletedAt(0)
