package swaggerdocs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "version": "{{.Version}}"
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "BETA",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Thunderdome API",
	Description:      "Thunderdome Planning Poker API for both Internal and External use.\nWARNING: Currently not considered stable and is subject to change until 1.0 is released.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
}
