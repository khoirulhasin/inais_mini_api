package error_handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khoirulhasin/globe_tracker_api/app/models"
	"github.com/khoirulhasin/globe_tracker_api/helpers"
)

func PanicHandler(c *gin.Context) (*models.Response, error) {
	if err := recover(); err != nil {
		str := fmt.Sprint(err)
		strArr := strings.Split(str, ":")

		key := strArr[0]
		msg := strings.Trim(strArr[1], " ")

		switch key {
		case DataNotFound.GetResponseStatus():
			c.Abort()
			return helpers.FailedResponseFormat(http.StatusBadRequest, msg)
		case Unauthorized.GetResponseStatus():
			c.Abort()
			return helpers.FailedResponseFormat(http.StatusUnauthorized, msg)
		default:
			c.Abort()
			return helpers.FailedResponseFormat(http.StatusInternalServerError, msg)
		}
	}

	return nil, nil
}

func PanicException(responseKey ResponseStatus) (*models.Response, error) {
	err := errors.New(responseKey.GetResponseMessage())
	err = fmt.Errorf("%s: %w", responseKey.GetResponseStatus(), err)
	if err != nil {
		panic(err)
	}

	return helpers.FailedResponseFormat(responseKey.GetResponseCode(), responseKey.GetResponseMessage())
}

func ListPanicException(responseKey ResponseStatus) (*models.ListResponse, error) {
	err := errors.New(responseKey.GetResponseMessage())
	err = fmt.Errorf("%s: %w", responseKey.GetResponseStatus(), err)
	if err != nil {
		panic(err)
	}

	return helpers.ListFailedResponseFormat(responseKey.GetResponseCode(), responseKey.GetResponseMessage())
}
