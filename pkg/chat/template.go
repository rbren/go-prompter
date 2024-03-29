package chat

import (
	"fmt"
	"strings"
	"text/template"
)

// fillTemplate fills a specified template with data and returns the resultant string.
func (s *Session) fillTemplate(fileName string, data map[string]interface{}) (string, error) {
	// Parse the template
	tmpl, err := template.New("tmpl").Funcs(s.templateFuncMap).ParseFS(s.templateFS, "prompts/*.md")
	if err != nil {
		return "", err
	}
	tmpl = tmpl.Lookup(fileName + ".md")
	if tmpl == nil {
		return "", fmt.Errorf("template %s not found", fileName)
	}

	// Fill out the template with the provided data
	var filledTemplate strings.Builder
	err = tmpl.Execute(&filledTemplate, data)
	if err != nil {
		return "", err
	}
	return filledTemplate.String(), nil
}
