package main

import "text/template"

var tmplFuncs = template.FuncMap{
	"match": func(matches [][]string, i int, j int) string {
		if len(matches) <= i {
			return ""
		}

		if len(matches[i]) <= j {
			return ""
		}

		return matches[i][j]
	},
	"item": func(items []string, i int) string {
		if len(items) <= i {
			return ""
		}

		return items[i]
	},
}
