package error_handlers

type ResponseStatus int
type Headers int
type General int

// Constant Api
// iota to get the order number start by 0
const (
	Success ResponseStatus = iota + 1
	DataNotFound
	UnknownError
	InvalidRequest
	Unauthorized
	InternalError
)

func (r ResponseStatus) GetResponseStatus() string {
	return [...]string{"SUCCESS", "DATA_NOT_FOUND", "UNKNOWN_ERROR", "INVALID_REQUEST", "UNAUTHORIZED", "INTERNAL_ERROR"}[r-1]
}
func (r ResponseStatus) GetResponseCode() int32 {
	return [...]int32{200, 404, 500, 400, 401, 500}[r-1]
}

func (r ResponseStatus) GetResponseMessage() string {
	return [...]string{"Success", "Data Not Found", "Unknown Error", "Invalid Request", "Unauthorized", "InternalError"}[r-1]
}
