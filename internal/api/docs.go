// Package api provides REST API documentation using Swagger/OpenAPI
package api

// @title REST API Documentation
// @version 1.0
// @description This is a REST API for managing Teachers, Students and Executives
// @host localhost:443
// @BasePath /
// @schemes https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT Bearer token

// SecurityScheme defines the type of security scheme in use.
type SecurityScheme struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	In          string `json:"in"`
	Name        string `json:"name"`
	Scheme      string `json:"scheme"`
}
