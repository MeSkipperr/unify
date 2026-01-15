package adb

import "strings"

func RenderTemplate(template string, values map[string]string) string {
	result := template

	for k, v := range values {
		result = strings.ReplaceAll(result, "{"+k+"}", v)
	}

	return result
}
