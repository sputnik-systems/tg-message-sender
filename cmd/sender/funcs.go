package main

import (
	"strings"
	"text/template"
)

var (
	tmplFuncs = template.FuncMap{
		"split":   split,
		"replace": replace,
	}
)

func split(input, sep string) []string {
	return strings.Split(input, sep)
}

func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}
