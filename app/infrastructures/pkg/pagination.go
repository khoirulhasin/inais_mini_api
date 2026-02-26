package pkg

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// JoinedTable defines metadata for a joined table's searchable columns
type JoinedTable struct {
	TableName string
	Alias     string
	Columns   map[string]string // Column name to database type (e.g., "name": "varchar")
}

// ColumnCache for caching column metadata to avoid repeated database queries
var columnCache = make(map[string]map[string]string)

// Paginate creates a reusable GORM query for pagination with support for JOINs and filters
func Paginate(value interface{}, pagination *models.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)
	pagination.TotalRows = totalRows
	limit, offset := setPaginationDefaults(pagination)
	totalPages := int64(math.Ceil(float64(totalRows) / float64(limit)))
	pagination.TotalPages = totalPages

	// Get table name for the primary model
	tableName, err := getTableName(value, db)
	if err != nil {
		log.Printf("Warning: Failed to get table name: %v, using fallback", err)
		tableName = reflectType(value).Name()
	}
	columns := pickColumns(value, db, tableName)

	allColumns := make(map[string]string)
	for colName, colType := range columns {
		allColumns[colName] = colType
	}

	// Track preloaded tables and their aliases
	joinedTables := make(map[string]int) // tableName -> count of occurrences
	aliases := make(map[string]string)   // preload -> alias or tableName

	// Prepare sort field
	sortField := ""
	sortOrder := ""
	if pagination.SortField != nil && *pagination.SortField != "" {
		sortField = *pagination.SortField
	}
	if pagination.SortOrder != nil && *pagination.SortOrder != "" {
		sortOrder = *pagination.SortOrder
	}
	sort := GetSort(tableName, sortField, sortOrder)
	pagination.Sort = &sort

	return func(db *gorm.DB) *gorm.DB {
		query := db.Offset(offset).Limit(limit)

		// First pass: Count preloaded table occurrences
		if db.Statement != nil && len(db.Statement.Preloads) > 0 {
			modelType := reflectType(value)
			for preload := range db.Statement.Preloads {
				fieldType, ok := getRelatedFieldType(modelType, preload)
				if !ok {
					log.Printf("Warning: Field %s not found in model", preload)
					continue
				}

				relatedTable, err := getTableNameFromStruct(fieldType, db)
				if err != nil {
					log.Printf("Failed to get table name for preload %s: %v", preload, err)
					continue
				}
				joinedTables[relatedTable]++
			}

			// Second pass: Process preloads and apply joins
			aliasCounter := make(map[string]int) // tableName -> alias count
			for preload := range db.Statement.Preloads {
				fieldType, ok := getRelatedFieldType(modelType, preload)
				if !ok {
					continue
				}

				relatedTable, err := getTableNameFromStruct(fieldType, db)
				if err != nil {
					log.Printf("Failed to get table name for preload %s: %v", preload, err)
					continue
				}

				// Determine if alias is needed
				joinName := relatedTable
				if joinedTables[relatedTable] > 1 {
					aliasCounter[relatedTable]++
					joinName = fmt.Sprintf("%s_%s%d", relatedTable, strings.ToLower(preload), aliasCounter[relatedTable])
				}
				aliases[preload] = joinName

				// Extract columns from related table
				relModel := reflect.New(fieldType).Interface()
				relatedColumns, err := extractColumnsFromStruct(relModel)
				if err != nil {
					log.Printf("Failed to extract columns for preload %s: %v", preload, err)
					continue
				}

				// Update allColumns with qualified column names
				allColumns = getRelatedColumns(relatedColumns, joinName, allColumns)

				// Assume foreign key: preload in snake_case + "_id"
				foreignKey := toSnakeCase(preload) + "_id"
				if !db.Migrator().HasColumn(tableName, foreignKey) {
					log.Printf("Warning: Foreign key column %s.%s does not exist, skipping JOIN for %s", tableName, foreignKey, preload)
					continue
				}

				// Perform LEFT JOIN
				query = query.Joins(fmt.Sprintf(
					"LEFT JOIN %s AS %s ON %s.id = %s.%s",
					relatedTable, joinName, joinName, tableName, foreignKey,
				))
			}
		}

		// Apply sorting
		query = query.Order(*pagination.Sort)

		// Apply search if provided
		if pagination.Search != nil && *pagination.Search != "" {
			i := 0
			for colName, colType := range allColumns {
				if colType == "" {
					continue
				}

				isNumeric := colType == "bigint" || colType == "int8" || colType == "int4" || colType == "int"
				isString := colType == "varchar" || colType == "text" || colType == "char"

				if isNumeric {
					if num, err := strconv.ParseInt(*pagination.Search, 10, 64); err == nil {
						if i == 0 {
							query = query.Where(colName+" = ?", num)
						} else {
							query = query.Or(colName+" = ?", num)
						}
						i++
					}
				} else if isString {
					if i == 0 {
						query = query.Where("LOWER("+colName+") LIKE LOWER(?)", "%"+*pagination.Search+"%")
					} else {
						query = query.Or("LOWER("+colName+") LIKE LOWER(?)", "%"+*pagination.Search+"%")
					}
					i++
				}
			}
		}
		log.Print(pagination.Filters, "FILTERS")
		// Apply filters if provided
		if len(pagination.Filters) > 0 {

			query = applyFilters(query, pagination.Filters, allColumns, aliases, tableName)
		}

		return query
	}
}

// Paginate creates a reusable GORM query for pagination with support for JOINs
// func Paginate(value interface{}, pagination *models.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
// 	var totalRows int64
// 	db.Model(value).Count(&totalRows)
// 	pagination.TotalRows = totalRows
// 	limit, offset := setPaginationDefaults(pagination)
// 	totalPages := int64(math.Ceil(float64(totalRows) / float64(limit)))
// 	pagination.TotalPages = totalPages

// 	// Get table name for the primary model
// 	tableName, err := getTableName(value, db)
// 	if err != nil {
// 		log.Printf("Warning: Failed to get table name: %v, using fallback", err)
// 		tableName = reflectType(value).Name()
// 	}
// 	columns := pickColumns(value, db, tableName)

// 	allColumns := make(map[string]string)
// 	for colName, colType := range columns {
// 		allColumns[colName] = colType
// 	}

// 	// Track preloaded tables and their aliases (if needed)
// 	joinedTables := make(map[string]int) // tableName -> count of occurrences
// 	aliases := make(map[string]string)   // preload -> alias or tableName

// 	// Prepare sort field
// 	sortField := ""
// 	sortOrder := ""
// 	if pagination.SortField != nil && *pagination.SortField != "" {
// 		sortField = *pagination.SortField
// 	}
// 	if pagination.SortOrder != nil && *pagination.SortOrder != "" {
// 		sortOrder = *pagination.SortOrder
// 	}
// 	sort := GetSort(tableName, sortField, sortOrder)
// 	pagination.Sort = &sort

// 	return func(db *gorm.DB) *gorm.DB {
// 		query := db.Offset(offset).Limit(limit)

// 		// First pass: Count preloaded table occurrences
// 		if db.Statement != nil && len(db.Statement.Preloads) > 0 {
// 			modelType := reflectType(value)
// 			for preload := range db.Statement.Preloads {
// 				fieldType, ok := getRelatedFieldType(modelType, preload)
// 				if !ok {
// 					log.Printf("Warning: Field %s not found in model", preload)
// 					continue
// 				}

// 				relatedTable, err := getTableNameFromStruct(fieldType, db)
// 				if err != nil {
// 					log.Printf("Failed to get table name for preload %s: %v", preload, err)
// 					continue
// 				}
// 				joinedTables[relatedTable]++
// 			}

// 			// Second pass: Process preloads and apply joins
// 			aliasCounter := make(map[string]int) // tableName -> alias count
// 			for preload := range db.Statement.Preloads {
// 				fieldType, ok := getRelatedFieldType(modelType, preload)
// 				if !ok {
// 					continue
// 				}

// 				relatedTable, err := getTableNameFromStruct(fieldType, db)
// 				if err != nil {
// 					log.Printf("Failed to get table name for preload %s: %v", preload, err)
// 					continue
// 				}

// 				// Determine if alias is needed
// 				joinName := relatedTable
// 				if joinedTables[relatedTable] > 1 {
// 					// Generate unique alias if multiple preloads reference the same table
// 					aliasCounter[relatedTable]++
// 					joinName = fmt.Sprintf("%s_%s%d", relatedTable, strings.ToLower(preload), aliasCounter[relatedTable])
// 				}
// 				aliases[preload] = joinName

// 				// Extract columns from related table
// 				relModel := reflect.New(fieldType).Interface()
// 				relatedColumns, err := extractColumnsFromStruct(relModel)
// 				if err != nil {
// 					log.Printf("Failed to extract columns for preload %s: %v", preload, err)
// 					continue
// 				}

// 				// Update allColumns with qualified column names
// 				allColumns = getRelatedColumns(relatedColumns, joinName, allColumns)

// 				// Assume foreign key: preload in snake_case + "_id"
// 				foreignKey := toSnakeCase(preload) + "_id"
// 				if !db.Migrator().HasColumn(tableName, foreignKey) {
// 					log.Printf("Warning: Foreign key column %s.%s does not exist, skipping JOIN for %s", tableName, foreignKey, preload)
// 					continue
// 				}

// 				// Perform LEFT JOIN
// 				query = query.Joins(fmt.Sprintf(
// 					"LEFT JOIN %s AS %s ON %s.id = %s.%s",
// 					relatedTable, joinName, joinName, tableName, foreignKey,
// 				))
// 			}
// 		}

// 		// Apply sorting
// 		query = query.Order(*pagination.Sort)

// 		// Apply search if provided
// 		if pagination.Search != nil && *pagination.Search != "" {
// 			i := 0
// 			for colName, colType := range allColumns {
// 				if colType == "" {
// 					continue
// 				}

// 				isNumeric := colType == "bigint" || colType == "int8" || colType == "int4" || colType == "int"
// 				isString := colType == "varchar" || colType == "text" || colType == "char"

// 				if isNumeric {
// 					if num, err := strconv.ParseInt(*pagination.Search, 10, 64); err == nil {
// 						if i == 0 {
// 							query = query.Where(colName+" = ?", num)
// 						} else {
// 							query = query.Or(colName+" = ?", num)
// 						}
// 						i++
// 					}
// 				} else if isString {
// 					if i == 0 {
// 						query = query.Where("LOWER("+colName+") LIKE LOWER(?)", "%"+*pagination.Search+"%")
// 					} else {
// 						query = query.Or("LOWER("+colName+") LIKE LOWER(?)", "%"+*pagination.Search+"%")
// 					}
// 					i++
// 				}
// 			}
// 		}

// 		return query
// 	}
// }

// getRelatedColumns updates column names with table alias
func getRelatedColumns(relatedColumns map[string]string, alias string, allColumns map[string]string) map[string]string {
	columnsToDelete := []string{"uuid", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by", "id"}
	for colName, colType := range relatedColumns {
		if colName == "" {
			continue
		}

		// Skip unwanted columns
		skip := false
		for _, delCol := range columnsToDelete {
			if strings.EqualFold(colName, delCol) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		// Use alias instead of table name for column qualification
		qualified := fmt.Sprintf("%s.%s", alias, colName)
		allColumns[qualified] = colType
	}
	return allColumns
}

// Custom snake_case naming strategy
var snakeCaseNamer = schema.NamingStrategy{
	SingularTable: true,                        // Optional: ensures table names are singular
	NameReplacer:  strings.NewReplacer("", ""), // Not strictly needed, but kept for compatibility
}

// toSnakeCase converts a string from CamelCase to snake_case
func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(toLower(r))
	}
	return strings.ToLower(result.String())
}

// toLower converts a rune to lowercase
func toLower(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + 32
	}
	return r
}
func extractColumnsFromStruct(model interface{}) (map[string]string, error) {
	columns := map[string]string{}

	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Slice && v.Len() > 0 {
		model = v.Index(0).Interface()
	}

	s, err := schema.Parse(model, &sync.Map{}, snakeCaseNamer)
	if err != nil {
		return columns, err
	}

	for _, field := range s.Fields {
		if !field.IgnoreMigration {
			colType := dataTypeToString(field.DataType) // seperti varchar, bigint, dll (bisa disesuaikan)
			if colType == "" {
				colType = "varchar" // fallback
			}
			columns[toSnakeCase(field.DBName)] = colType
		}
	}

	return columns, nil
}

func dataTypeToString(dt schema.DataType) string {
	switch dt {
	case schema.Bool:
		return "bool"
	case schema.Int:
		return "int"
	case schema.Uint:
		return "uint"
	case schema.Float:
		return "float"
	case schema.String:
		return "varchar"
	case schema.Time:
		return "timestamp"
	case schema.Bytes:
		return "bytes"
	default:
		return ""
	}
}

// pickColumns extracts column data with table prefix to avoid ambiguity in JOINs
func pickColumns(value interface{}, db *gorm.DB, tableName string) map[string]string {
	columnsToDelete := []string{"uuid", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by", "id"}
	columns := make(map[string]string)

	// Get column types from GORM
	columnTypes, err := db.Migrator().ColumnTypes(value)
	if err != nil {
		log.Printf("Error getting column types for table %s: %v", tableName, err)
		return columns
	}

	// Map columns to their database types with table prefix
	for _, col := range columnTypes {
		columnName := col.Name()
		columnType := col.DatabaseTypeName()
		// Use tableName.columnName to fully qualify the column
		qualifiedColumnName := tableName + "." + columnName
		if _, exists := columns[qualifiedColumnName]; !exists {
			columns[qualifiedColumnName] = strings.ToLower(columnType)
		}
	}

	// Remove unwanted columns
	for _, column := range columnsToDelete {
		delete(columns, tableName+"."+strings.ToLower(column))
	}

	return columns
}

// GetSort constructs the ORDER BY clause with table-qualified column names
func GetSort(tableName, sortField, sortOrder string) string {
	// Ensure tableName is not empty
	if tableName == "" {
		tableName = "unknown_table" // Fallback table name to avoid empty prefix
	}

	if sortField == "" && sortOrder == "" {
		return tableName + ".id desc"
	}

	if sortField == "" {
		return tableName + ".id " + sortOrder
	}

	if sortOrder == "" {
		return tableName + "." + sortField + " desc"
	}

	if strings.Contains(sortField, ".") {
		return sortField + " " + sortOrder
	}

	return tableName + "." + sortField + " " + sortOrder
}

// setPaginationDefaults sets default values for pagination parameters
func setPaginationDefaults(pagination *models.Pagination) (limit, offset int) {
	defaultLimit := 10
	defaultOffset := 0

	if pagination.Limit != nil {
		limit = *pagination.Limit
	} else {
		limit = defaultLimit
		pagination.Limit = &defaultLimit
	}

	if pagination.Offset != nil {
		offset = *pagination.Offset
	} else {
		offset = defaultOffset
		pagination.Offset = &defaultOffset
	}

	return limit, offset
}

// getTableName retrieves the table name for the given model
func getTableName(value interface{}, db *gorm.DB) (string, error) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct && v.Kind() != reflect.Slice {
		return "", fmt.Errorf("invalid model type: %v", v.Kind())
	}

	t := reflectType(value)
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	tableName := db.NamingStrategy.TableName(t.Name())
	if tableName == "" {
		tableName = t.Name()
	}

	if tableName == "" {
		return "", fmt.Errorf("unable to determine table name for model: %v", t.Name())
	}

	return strings.ToLower(tableName), nil
}

func getTableNameFromStruct(structType reflect.Type, db *gorm.DB) (string, error) {
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	s, err := schema.Parse(reflect.New(structType).Interface(), &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return "", err
	}
	return s.Table, nil
}

// pagePayloadIsNil handles pagination input defaults
func pagePayloadIsNil(pageInput models.PageInput) (int, int, string, string, string, []*models.FilterInput) {
	limit := 10
	offset := 0
	sortField := "id"
	sortOrder := "desc"
	search := ""
	filters := []*models.FilterInput{} // Initialize as empty slice

	if pageInput.Limit != nil {
		limit = int(*pageInput.Limit)
	}

	if pageInput.Offset != nil {
		offset = int(*pageInput.Offset)
	}

	if pageInput.SortField != nil {
		sortField = *pageInput.SortField
	}

	if pageInput.SortOrder != nil {
		sortOrder = *pageInput.SortOrder
	}

	if pageInput.Search != nil {
		search = *pageInput.Search
	}

	if len(pageInput.Filters) > 0 {
		filters = pageInput.Filters
	}

	return limit, offset, sortField, sortOrder, search, filters
}

// PageInputIsNil handles nil page input and returns default pagination parameters
func PageInputIsNil(pageInput *models.PageInput) (int, int, string, string, string, []*models.FilterInput) {
	limit := 10
	offset := 0
	sortField := ""
	sortOrder := ""
	search := ""
	filters := []*models.FilterInput{} // Initialize as empty slice

	if pageInput != nil {
		limit, offset, sortField, sortOrder, search, filters = pagePayloadIsNil(*pageInput)
	}

	return limit, offset, sortField, sortOrder, search, filters
}

// Helper function to get the reflect type of the value
func reflectType(value interface{}) reflect.Type {
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func getRelatedFieldType(parentType reflect.Type, fieldName string) (reflect.Type, bool) {
	field, ok := parentType.FieldByName(fieldName)
	if !ok {
		return nil, false
	}
	return field.Type, true
}

// applyFilters applies filters to the GORM query
func applyFilters(query *gorm.DB, filters []*models.Filter, allColumns map[string]string, aliases map[string]string, tableName string) *gorm.DB {

	log.Print("OKOK", filters[0].Value)
	supportedOperators := map[string]bool{
		"=":        true,
		"!=":       true,
		">":        true,
		"<":        true,
		">=":       true,
		"<=":       true,
		"LIKE":     true,
		"NOT LIKE": true,
		"IN":       true, // Added support for IN operator
		"NOT IN":   true,
	}

	for _, filter := range filters {
		if filter.Key == "" || filter.Operator == "" || filter.Value == nil {
			log.Printf("Warning: Invalid filter, skipping: %+v", filter)
			continue
		}
		// Validate operator
		if !supportedOperators[filter.Operator] {
			log.Printf("Warning: Unsupported operator %s for filter key %s, skipping", filter.Operator, filter.Key)
			continue
		}

		// Resolve column name (handle nested fields like vehicle.make)
		colName := filter.Key
		if strings.Contains(colName, ".") {
			parts := strings.Split(colName, ".")
			if len(parts) > 1 {
				if parts[0] == tableName {
					colName = fmt.Sprintf("%s.%s", tableName, toSnakeCase(parts[1]))
				} else {
					preloadPath := strings.Join(parts[:len(parts)-1], ".")
					lastField := parts[len(parts)-1]
					if alias, ok := aliases[preloadPath]; ok {
						colName = fmt.Sprintf("%s.%s", alias, toSnakeCase(lastField))
					} else {
						log.Printf("Warning: No alias found for preload path %s, skipping filter", preloadPath)
						continue
					}
				}
			}
		} else {
			colName = fmt.Sprintf("%s.%s", tableName, toSnakeCase(colName))
		}

		// Check if column exists
		colType, exists := allColumns[colName]
		if !exists {
			log.Printf("Warning: Column %s not found, skipping filter", colName)
			continue
		}

		// Handle different column types and operators
		isNumeric := colType == "bigint" || colType == "int8" || colType == "int4" || colType == "int"
		isString := colType == "varchar" || colType == "text" || colType == "char"
		isBool := colType == "bool"

		switch filter.Operator {
		case "IN", "NOT IN":
			if isNumeric {
				values, err := convertToNumericSlice(filter.Value)
				if err != nil {
					log.Printf("Warning: Invalid numeric values for %s filter %s: %v, skipping", filter.Operator, colName, err)
					continue
				}
				query = query.Where(fmt.Sprintf("%s %s ?", colName, filter.Operator), values)
			} else if isString {
				values, err := convertToStringSlice(filter.Value)
				if err != nil {
					log.Printf("Warning: Invalid string values for %s filter %s: %v, skipping", filter.Operator, colName, err)
					continue
				}
				query = query.Where(fmt.Sprintf("%s %s ?", colName, filter.Operator), values)
			} else if isBool {
				values, err := convertToBoolSlice(filter.Value)
				if err != nil {
					log.Printf("Warning: Invalid boolean values for %s filter %s: %v, skipping", filter.Operator, colName, err)
					continue
				}
				query = query.Where(fmt.Sprintf("%s %s ?", colName, filter.Operator), values)
			} else {
				log.Printf("Warning: %s operator not supported for column type %s in filter %s, skipping", filter.Operator, colType, colName)
				continue
			}
		case "LIKE", "NOT LIKE":
			if !isString {
				log.Printf("Warning: Operator %s is only supported for string columns, skipping filter for %s", filter.Operator, colName)
				continue
			}
			if val, ok := filter.Value.(string); ok {
				query = query.Where(fmt.Sprintf("LOWER(%s) %s LOWER(?)", colName, filter.Operator), "%"+val+"%")
			} else {
				log.Printf("Warning: Value for filter %s must be a string for operator %s, skipping", colName, filter.Operator)
				continue
			}
		case "=", "!=", ">", "<", ">=", "<=":
			if isNumeric {
				// Handle numeric values
				switch v := filter.Value.(type) {
				case int, int32, int64:
					query = query.Where(fmt.Sprintf("%s %s ?", colName, filter.Operator), v)
				case float32, float64:
					query = query.Where(fmt.Sprintf("%s %s ?", colName, filter.Operator), v)
				case string:
					if num, err := strconv.ParseInt(v, 10, 64); err == nil {
						query = query.Where(fmt.Sprintf("%s %s ?", colName, filter.Operator), num)
					} else if num, err := strconv.ParseFloat(v, 64); err == nil {
						query = query.Where(fmt.Sprintf("%s %s ?", colName, filter.Operator), num)
					} else {
						log.Printf("Warning: Invalid numeric value %v for filter %s, skipping", v, colName)
						continue
					}
				default:
					log.Printf("Warning: Unsupported value type %T for numeric filter %s, skipping", v, colName)
					continue
				}
			} else if isString {
				// Handle string values
				if val, ok := filter.Value.(string); ok {
					query = query.Where(fmt.Sprintf("%s %s ?", colName, filter.Operator), val)
				} else {
					log.Printf("Warning: Value for filter %s must be a string, skipping", colName)
					continue
				}
			} else if isBool {
				// Handle boolean values
				switch v := filter.Value.(type) {
				case bool:
					query = query.Where(fmt.Sprintf("%s %s ?", colName, filter.Operator), v)
				case string:
					if val, err := strconv.ParseBool(v); err == nil {
						query = query.Where(fmt.Sprintf("%s %s ?", colName, filter.Operator), val)
					} else {
						log.Printf("Warning: Invalid boolean value %v for filter %s, skipping", v, colName)
						continue
					}
				default:
					log.Printf("Warning: Unsupported value type %T for boolean filter %s, skipping", v, colName)
					continue
				}
			} else {
				log.Printf("Warning: Unsupported column type %s for filter %s, skipping", colType, colName)
				continue
			}
		default:
			log.Printf("Warning: Unsupported operator %s for filter %s, skipping", filter.Operator, colName)
			continue
		}
	}

	return query
}

// Helper functions to convert filter.Value to appropriate slice types for IN operator
// convertToNumericSlice converts a slice to a numeric slice for IN/NOT IN operators
func convertToNumericSlice(value interface{}) (interface{}, error) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("value must be a slice for IN/NOT IN operator")
	}

	length := v.Len()
	if length == 0 {
		return nil, fmt.Errorf("slice is empty for IN/NOT IN operator")
	}

	// Try to determine the type of the first element
	firstElem := v.Index(0).Interface()
	switch firstElem.(type) {
	case int, int32, int64:
		result := make([]int64, length)
		for i := 0; i < length; i++ {
			elem := v.Index(i).Interface()
			switch val := elem.(type) {
			case int:
				result[i] = int64(val)
			case int32:
				result[i] = int64(val)
			case int64:
				result[i] = val
			case float32:
				result[i] = int64(val)
			case float64:
				result[i] = int64(val)
			case string:
				num, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid numeric string value at index %d: %s", i, val)
				}
				result[i] = num
			default:
				return nil, fmt.Errorf("unsupported element type %T at index %d", elem, i)
			}
		}
		return result, nil
	case float32, float64:
		result := make([]float64, length)
		for i := 0; i < length; i++ {
			elem := v.Index(i).Interface()
			switch val := elem.(type) {
			case int:
				result[i] = float64(val)
			case int32:
				result[i] = float64(val)
			case int64:
				result[i] = float64(val)
			case float32:
				result[i] = float64(val)
			case float64:
				result[i] = val
			case string:
				num, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid float string value at index %d: %s", i, val)
				}
				result[i] = num
			default:
				return nil, fmt.Errorf("unsupported element type %T at index %d", elem, i)
			}
		}
		return result, nil
	default:
		// Handle []interface{} (common with JSON unmarshaling)
		if _, ok := firstElem.(interface{}); ok {
			// Try to convert to int64 slice first
			intResult := make([]int64, length)
			for i := 0; i < length; i++ {
				elem := v.Index(i).Interface()
				switch val := elem.(type) {
				case float64: // JSON numbers are often float64
					intVal := int64(val)
					if float64(intVal) != val { // Check if it's a true integer
						return nil, fmt.Errorf("non-integer float value at index %d: %v", i, val)
					}
					intResult[i] = intVal
				case int:
					intResult[i] = int64(val)
				case int64:
					intResult[i] = val
				case string:
					num, err := strconv.ParseInt(val, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("invalid numeric string value at index %d: %s", i, val)
					}
					intResult[i] = num
				default:
					return nil, fmt.Errorf("unsupported element type %T at index %d", elem, i)
				}
			}
			return intResult, nil
		}
		return nil, fmt.Errorf("unsupported slice element type: %v", reflect.TypeOf(firstElem))
	}
}

func convertToStringSlice(value interface{}) ([]string, error) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("value must be a slice for IN operator")
	}

	result := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		if str, ok := v.Index(i).Interface().(string); ok {
			result[i] = str
		} else {
			return nil, fmt.Errorf("invalid string value at index %d", i)
		}
	}
	return result, nil
}

func convertToBoolSlice(value interface{}) ([]bool, error) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("value must be a slice for IN operator")
	}

	result := make([]bool, v.Len())
	for i := 0; i < v.Len(); i++ {
		switch val := v.Index(i).Interface().(type) {
		case bool:
			result[i] = val
		case string:
			if b, err := strconv.ParseBool(val); err == nil {
				result[i] = b
			} else {
				return nil, fmt.Errorf("invalid boolean string value at index %d: %s", i, val)
			}
		default:
			return nil, fmt.Errorf("invalid boolean value type at index %d", i)
		}
	}
	return result, nil
}
