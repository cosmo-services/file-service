package health_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	msg string
}

func NewHealthController() HealthController {
	return HealthController{
		msg: "healthy",
	}
}

// Health godoc
//
//	@Summary		Health endpoint
//	@Description	Check if the service is alive
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Successful response"
//	@Router			/health [get]
func (c HealthController) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": c.msg,
	})
}
