package handler

import (
	"io/ioutil"
	"log"
	"net/http"
)

var swaggerJSON string

// InitSwagger loads the swagger specification from file
func InitSwagger() {
	swaggerPath := "docs/swagger.json"

	// Try to read the swagger.json file
	content, err := ioutil.ReadFile(swaggerPath)
	if err != nil {
		log.Printf("Warning: Could not load swagger.json from %s: %v\n", swaggerPath, err)
		// Fallback to empty spec if file not found
		swaggerJSON = `{
  "swagger": "2.0",
  "info": {
    "title": "REST API Documentation",
    "version": "1.0",
    "description": "This is a REST API for managing Teachers, Students and Executives"
  },
  "host": "localhost:443",
  "basePath": "/",
  "schemes": ["https"],
  "paths": {}
}`
		return
	}

	swaggerJSON = string(content)
}

// SwaggerHandler serves the Swagger UI
// @Summary Serve Swagger UI
// @Description Serves the Swagger API documentation
// @Tags documentation
// @Produce html
// @Success 200 {object} string "Swagger UI HTML"
// @Router /swagger [get]
func SwaggerHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>REST API Swagger UI</title>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@3/swagger-ui.css">
		<style>
			html {
				box-sizing: border-box;
				overflow: -moz-scrollbars-vertical;
				overflow-y: scroll;
			}
			*, *:before, *:after {
				box-sizing: inherit;
			}
			body {
				margin: 0;
				padding: 0;
			}
		</style>
	</head>
	<body>
		<div id="swagger-ui"></div>
		<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@3/swagger-ui-bundle.js"> </script>
		<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@3/swagger-ui-standalone-preset.js"> </script>
		<script>
			const ui = SwaggerUIBundle({
				url: "/swagger.json",
				dom_id: '#swagger-ui',
				presets: [
					SwaggerUIBundle.presets.apis,
					SwaggerUIStandalonePreset
				],
				layout: "StandaloneLayout"
			})
		</script>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// SwaggerJSONHandler serves the Swagger JSON specification
// @Summary Serve Swagger specification
// @Description Serves the Swagger specification in JSON format
// @Tags documentation
// @Produce json
// @Success 200 {object} object "Swagger JSON specification"
// @Router /swagger.json [get]
func SwaggerJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(swaggerJSON))
}

// swaggerSpec will contain the generated Swagger specification
var swaggerSpec = `{
  "swagger": "2.0",
  "info": {
    "title": "REST API Documentation",
    "version": "1.0",
    "description": "This is a REST API for managing Teachers, Students and Executives"
  },
  "host": "localhost:443",
  "basePath": "/",
  "schemes": ["https"],
  "paths": {}
}`
