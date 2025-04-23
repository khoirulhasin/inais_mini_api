package error_handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khoirulhasin/globe_tracker_api/app/models"
)

func PanicHandler(c *gin.Context) (any, error) {
	if err := recover(); err != nil {
		str := fmt.Sprint(err)
		strArr := strings.Split(str, ":")

		key := strArr[0]
		msg := strings.Trim(strArr[1], " ")

		switch key {
		case DataNotFound.GetResponseStatus():
			c.Abort()
			return &models.Response{
				Message: msg,
				Status:  http.StatusBadRequest,
			}, nil
		case Unauthorized.GetResponseStatus():
			c.Abort()
			return &models.Response{
				Message: msg,
				Status:  http.StatusUnauthorized,
			}, nil
		default:
			c.Abort()
			return &models.Response{
				Message: msg,
				Status:  http.StatusInternalServerError,
			}, nil
		}
	}

	return nil, nil
}
