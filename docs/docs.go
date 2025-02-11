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
        "/currency/add": {
            "post": {
                "description": "Adds a currency to the tracking list.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "currency"
                ],
                "summary": "Add currency to tracking",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Currency symbol (e.g., BTCUSDT)",
                        "name": "symbol",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Currency added successfully"
                    },
                    "400": {
                        "description": "Bad Request (e.g., invalid symbol)",
                        "schema": {
                            "$ref": "#/definitions/delivery.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.Response"
                        }
                    }
                }
            }
        },
        "/currency/price": {
            "get": {
                "description": "Retrieves the price of a specific currency at a given timestamp.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "currency"
                ],
                "summary": "Get price of a currency at a specific timestamp",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Currency symbol (e.g., BTCUSDT)",
                        "name": "symbol",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Currency symbol (e.g., BTCUSDT)",
                        "name": "timestamp",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Price retrieved successfully",
                        "schema": {
                            "$ref": "#/definitions/delivery.Response"
                        }
                    },
                    "400": {
                        "description": "Response \"Bad Request (e.g., invalid symbol)",
                        "schema": {
                            "$ref": "#/definitions/delivery.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.Response"
                        }
                    }
                }
            }
        },
        "/currency/remove": {
            "delete": {
                "description": "Removes a currency from the tracking list.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "currency"
                ],
                "summary": "Remove currency from tracking",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Currency symbol (e.g., BTCUSDT)",
                        "name": "symbol",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Currency added successfully"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "delivery.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "$ref": "#/definitions/delivery.errorResponse"
                },
                "response": {}
            }
        },
        "delivery.errorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "text": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Crypto service API",
	Description:      "REST API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
