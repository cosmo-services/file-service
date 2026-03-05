package bootstrap

import (
	"main/internal/application/api/v2"
	"main/internal/application/jobs"
	"main/internal/application/nats"
	"main/internal/config"
	"main/pkg"

	auth_infrastructure "main/internal/infrastructure/auth"
	file_infrastructure "main/internal/infrastructure/file"

	auth_api "main/internal/application/api/v2/auth"
	file_api "main/internal/application/api/v2/file"
	health_api "main/internal/application/api/v2/health"
	swagger_api "main/internal/application/api/v2/swagger"

	file_domain "main/internal/domain/file"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	config.Module,
	pkg.Module,

	file_domain.Module,

	auth_infrastructure.Module,
	file_infrastructure.Module,

	api.Module,
	jobs.Module,
	nats.Module,
	auth_api.Module,
	health_api.Module,
	swagger_api.Module,
	file_api.Module,
)
