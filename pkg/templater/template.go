package templater

import (
	"bytes"
	"log"
	"regexp"
	"strings"
	"text/template"
)

func ExecToString(baseTemplate string, vuln any) (string, error) {
	compiled, err := template.New("").Parse(baseTemplate)
	if err != nil {
		log.Printf("Error compiling template: %s", err)
		return "", err
	}

	var buff bytes.Buffer
	if err := compiled.Execute(&buff, &vuln); err != nil {
		log.Printf("Error executing template: %s", err)
		return "", err
	}

	note := regexp.MustCompile(`[\n{}]`).ReplaceAllString(buff.String(), "")
	note = strings.ReplaceAll(note, "\"", "`")

	return note, nil
}
