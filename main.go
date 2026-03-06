package main

import (
	"main/cmd"
	_ "main/docs"
)

// @title Cosmo File API
// @version 2.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @security BearerAuth

// @BasePath /api/v2/file
func main() {
	cmd.StartApp()
}
