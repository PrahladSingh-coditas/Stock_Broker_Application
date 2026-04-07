package commons

import (
	"fmt"
	"net/http"
	genericModels "stock_broker_application/src/models"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(key, errMessage, err, action string, ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
		Message: genericModels.ErrorMessage{
			Key:          key,
			ErrorMessage: errMessage,
		},
		Error: fmt.Sprintf(err, action),
	})
}
