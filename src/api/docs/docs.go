// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/health": {
            "get": {
                "security": [
                    {
                        "None": []
                    }
                ],
                "description": "Check if the API is up and running",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.Message"
                        }
                    }
                }
            }
        },
        "/api/v1/auth": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Check if the API is up and running",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Authenticated Health check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.Message"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.Message": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "prospector.ie",
	BasePath:         "/api",
	Schemes:          []string{"https"},
	Title:            "Prospector API",
	Description:      "API backend for Prospector",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
