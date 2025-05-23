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
        "/v1/restaurants/getRestaurants": {
            "post": {
                "description": "Returns all restaurants",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "restuarants"
                ],
                "summary": "Get list of restuarants",
                "parameters": [
                    {
                        "description": "Location and Filter Parameters",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.GetRestauranstByLocationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.GetRestauranstByLocationResponse"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.GeoJSON": {
            "type": "object",
            "properties": {
                "coordinates": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "models.GetRestauranstByLocation": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string"
                },
                "address": {
                    "type": "string"
                },
                "closingTime": {},
                "distanceInKms": {
                    "type": "number"
                },
                "isOpen": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "openingTime": {}
            }
        },
        "models.GetRestauranstByLocationRequest": {
            "type": "object",
            "required": [
                "userLocation"
            ],
            "properties": {
                "distance": {
                    "type": "number"
                },
                "isOpen": {
                    "type": "boolean"
                },
                "orderBy": {
                    "type": "integer"
                },
                "sortBy": {
                    "type": "string"
                },
                "userLocation": {
                    "$ref": "#/definitions/models.GeoJSON"
                }
            }
        },
        "models.GetRestauranstByLocationResponse": {
            "type": "object",
            "properties": {
                "restaurants": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.GetRestauranstByLocation"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Food Delivery API",
	Description:      "API Server for Food Delivery App",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
