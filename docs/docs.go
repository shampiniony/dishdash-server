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
        "/cards": {
            "get": {
                "description": "Get a list of cards from the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cards"
                ],
                "summary": "Get cards",
                "responses": {
                    "200": {
                        "description": "List of cards",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/card.cardOutput"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "post": {
                "description": "Create a new card in the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cards"
                ],
                "summary": "Create a card",
                "parameters": [
                    {
                        "description": "Card data",
                        "name": "card",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/card.cardInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Saved card",
                        "schema": {
                            "$ref": "#/definitions/card.cardOutput"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/cards/tags": {
            "post": {
                "description": "Create a new tag in the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cards"
                ],
                "summary": "Create a tag",
                "parameters": [
                    {
                        "description": "Tag data",
                        "name": "tag",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/card.tagInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Saved tag",
                        "schema": {
                            "$ref": "#/definitions/card.tagOutput"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "card.cardInput": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "image": {
                    "type": "string"
                },
                "location": {
                    "$ref": "#/definitions/domain.Coordinate"
                },
                "price": {
                    "type": "integer"
                },
                "shortDescription": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "card.cardOutput": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "image": {
                    "type": "string"
                },
                "location": {
                    "$ref": "#/definitions/domain.Coordinate"
                },
                "price": {
                    "type": "integer"
                },
                "shortDescription": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/card.tagOutput"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "card.tagInput": {
            "type": "object",
            "properties": {
                "icon": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "card.tagOutput": {
            "type": "object",
            "properties": {
                "icon": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "domain.Coordinate": {
            "type": "object",
            "properties": {
                "lat": {
                    "type": "number"
                },
                "lon": {
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
