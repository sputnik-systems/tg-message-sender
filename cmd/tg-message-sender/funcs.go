package main

import (
	"strings"
	"text/template"
)

var (
	tmplFuncs = template.FuncMap{
		"split":   split,
		"replace": replace,
		"empty":   empty,
	}
)

func split(input, sep string) []string {
	return strings.Split(input, sep)
}

func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

func empty(input string) bool {
	return len(input) == 0
}
