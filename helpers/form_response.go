package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/khoirulhasin/globe_tracker_api/app/models"
)

func SuccessResponseFormat() (*models.Response, error) {
	return &models.Response{
		Message: "Successfully",
		Status:  http.StatusOK,
	}, nil
}

func FailedResponseFormat(code int32, message string) (*models.Response, error) {
	return &models.Response{
		Message: message,
		Status:  code,
	}, nil
}

func ListFailedResponseFormat(code int32, message string) (*models.ListResponse, error) {
	return &models.ListResponse{
		Message: message,
		Status:  code,
	}, nil
}

func ResponseDataFormat(data interface{}, model string) (*models.Response, error) {
	modelName, _ := convertToModel(data, model)
	return &models.Response{
		Message: "Successfully",
		Status:  http.StatusOK,
		Data:    modelName,
	}, nil
}

func ListResponseDataFormat(data interface{}, model string) (*models.ListResponse, error) {

	modelName, _ := convertToModelSlice(data, model)
	return &models.ListResponse{
		Message: "Successfully",
		Status:  http.StatusOK,
		Data:    modelName,
	}, nil
}

var modelRegistry = map[string]reflect.Type{
	"Answer":         reflect.TypeOf(models.Answer{}),
	"Question":       reflect.TypeOf(models.Question{}),
	"QuestionOption": reflect.TypeOf(models.QuestionOption{}),
}

func convertToModel(data interface{}, modelName string) (models.Data, error) {
	modelType, found := modelRegistry[modelName]
	if !found {
		return nil, fmt.Errorf("model type '%s' not registered", modelName)
	}

	// Buat instance baru dari model
	newInstance := reflect.New(modelType).Interface()

	// Decode dari data asli ke struct instance
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}

	err = json.Unmarshal(jsonBytes, newInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal to %s: %w", modelName, err)
	}

	// Kembalikan hasilnya sebagai models.Data
	if d, ok := newInstance.(models.Data); ok {
		return d, nil
	}

	// Coba dereference pointer kalau perlu
	if d, ok := reflect.ValueOf(newInstance).Elem().Interface().(models.Data); ok {
		return d, nil
	}

	return nil, fmt.Errorf("converted object does not implement models.Data")
}

func convertToModelSlice(data interface{}, modelName string) ([]models.Data, error) {
	modelType, found := modelRegistry[modelName]
	if !found {
		return nil, fmt.Errorf("model type '%s' not registered", modelName)
	}

	// data (bisa []map[string]any atau []any), marshal dulu
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input data: %w", err)
	}

	// Buat slice dari model (bukan pointer)
	sliceType := reflect.SliceOf(modelType)
	slicePtr := reflect.New(sliceType).Interface()

	// Unmarshal ke slice tersebut
	err = json.Unmarshal(jsonBytes, slicePtr)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal into model slice: %w", err)
	}

	// Convert hasilnya ke []models.Data
	resultSlice := reflect.ValueOf(slicePtr).Elem()
	var finalData []models.Data

	for i := 0; i < resultSlice.Len(); i++ {
		item := resultSlice.Index(i).Interface()

		if d, ok := item.(models.Data); ok {
			finalData = append(finalData, d)
		} else {
			// coba dereference jika struct pointer
			if d, ok := reflect.ValueOf(item).Elem().Interface().(models.Data); ok {
				finalData = append(finalData, d)
			} else {
				return nil, fmt.Errorf("item at index %d does not implement models.Data", i)
			}
		}
	}

	return finalData, nil
}
