package controllers

import (
	"loan-tracker/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogController struct {
	logUsecase domain.LogUsecase
}

func NewLogController(logUsecase domain.LogUsecase) *LogController {
	return &LogController{
		logUsecase: logUsecase,
	}
}

func (c *LogController) ViewSystemLogs(ctx *gin.Context) {
	logs, err := c.logUsecase.GetSystemLogs(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, logs)
}
