package http

import (
	"strings"
	"testing"
)

func TestPersonaTemplateRequestValidation(t *testing.T) {
	persona := func(name string) personaTemplatePersonaRequestBody {
		return personaTemplatePersonaRequestBody{Name: name}
	}

	tooMany := make([]personaTemplatePersonaRequestBody, 0, 33)
	for i := range 33 {
		tooMany = append(tooMany, persona(strings.Repeat("a", i+1)))
	}

	tests := []struct {
		name     string
		template personaTemplateRequestBody
		valid    bool
	}{
		{"valid", personaTemplateRequestBody{Name: "SaaS", Personas: []personaTemplatePersonaRequestBody{persona("Admin")}}, true},
		{"blank name", personaTemplateRequestBody{Name: "  ", Personas: []personaTemplatePersonaRequestBody{persona("Admin")}}, false},
		{"no personas", personaTemplateRequestBody{Name: "SaaS"}, false},
		{"blank persona name", personaTemplateRequestBody{Name: "SaaS", Personas: []personaTemplatePersonaRequestBody{persona(" ")}}, false},
		{"long persona name", personaTemplateRequestBody{Name: "SaaS", Personas: []personaTemplatePersonaRequestBody{persona(strings.Repeat("a", 257))}}, false},
		{"duplicate persona names", personaTemplateRequestBody{Name: "SaaS", Personas: []personaTemplatePersonaRequestBody{persona("Admin"), persona(" Admin ")}}, false},
		{"too many personas", personaTemplateRequestBody{Name: "SaaS", Personas: tooMany}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.template.normalize()
			err := validate.Struct(tt.template)
			if (err == nil) != tt.valid {
				t.Errorf("expected valid %v, got error %v", tt.valid, err)
			}
		})
	}
}
