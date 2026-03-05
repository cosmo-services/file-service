package bootstrap

import (
	"main/internal/application/api/v2"
	"main/internal/application/jobs"
	"main/internal/config"
	"main/pkg"

	auth_infrastructure "main/internal/infrastructure/auth"
	file_infrastructure "main/internal/infrastructure/file"

	health_api "main/internal/application/api/v2/health"
	swagger_api "main/internal/application/api/v2/swagger"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	config.Module,
	pkg.Module,

	auth_infrastructure.Module,
	file_infrastructure.Module,

	api.Module,
	jobs.Module,
	health_api.Module,
	swagger_api.Module,
)
