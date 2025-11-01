package root

import (
	"bsnack/lib/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(context *gin.Context) {
	response.Json(context, http.StatusOK, nil)
}
