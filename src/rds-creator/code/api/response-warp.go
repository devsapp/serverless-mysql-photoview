package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	StatusFail           = "FAILED"
	StatusOk             = "SUCCESS"
	FcErrorStatusCode    = 417
	FcErrorStatusCodeStr = "417"
)

func RespondError(c *gin.Context, errCode int, errMsg string) {
	c.Header("x-fc-status", FcErrorStatusCodeStr)
	c.JSON(FcErrorStatusCode, errMsg)
}

func respondInternalError(c *gin.Context, errMsg string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  StatusFail,
		"message": errMsg,
	})
}

func respondOk(c *gin.Context, result interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"result": result,
	})
}
