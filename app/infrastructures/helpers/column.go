package helpers

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GetColumns(value interface{}, db *gorm.DB) map[string]string {
	// get columns from database
	s, err := schema.Parse(&value, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		panic("failed to parse schema")
	}

	fmt.Print(s.Table)

	m := make(map[string]string)
	for _, field := range s.Fields {
		dbName := field.DBName
		modelName := field.Name
		if dbName != "" {
			m[modelName] = s.Table + "." + dbName
		}

	}

	return m
}
