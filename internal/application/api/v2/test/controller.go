package test_api

import (
	"main/internal/domain/test"
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestController struct {
	repo   test.TestRepo
	logger pkg.Logger
}

func NewTestController(repo test.TestRepo, logger pkg.Logger) *TestController {
	return &TestController{
		repo:   repo,
		logger: logger,
	}
}

// Test godoc
//
//	@Summary		Test endpoint
//	@Description	Test service
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Successful response"
//	@Router			/test [get]
func (controller *TestController) GetOne(ctx *gin.Context) {
	result, err := controller.repo.GetOne()
	if err != nil {
		controller.logger.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, result)
}
