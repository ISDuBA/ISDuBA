// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "ISDuBA team",
            "url": "https://github.com/ISDuBA/ISDuBA",
            "email": "info@intevation.de"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/documents": {
            "post": {
                "description": "Upload endpoint for CSAF documents.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Imports a CSAF document.",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/results.ID"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/results.Error"
                        }
                    },
                    "403": {
                        "description": "False TLP or publisher",
                        "schema": {
                            "$ref": "#/definitions/results.Error"
                        }
                    },
                    "409": {
                        "description": "Already in database",
                        "schema": {
                            "$ref": "#/definitions/results.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/results.Error"
                        }
                    }
                }
            }
        },
        "/documents/{id}": {
            "delete": {
                "description": "Delete endpoint for CSAF documents.",
                "produces": [
                    "application/json"
                ],
                "summary": "Deletes a CSAF document.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Document ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/results.ID"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/results.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/results.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/results.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "results.Error": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "error": {
                    "type": "string"
                }
            }
        },
        "results.ID": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "ISDuBA API",
	Description:      "This is the ISDuBA API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
