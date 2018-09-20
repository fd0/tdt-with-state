package main

import "strings"

// Capitalize returns input with the first character converted to upper case.
func Capitalize(input string) (result string) {

	first := input[0:1]
	rest := input[1:]

	return strings.ToUpper(first) + rest
}
