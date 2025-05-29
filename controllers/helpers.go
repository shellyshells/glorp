package controllers

import "html/template"

var TemplateFuncMap template.FuncMap

func init() {
	TemplateFuncMap = template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
	}
}