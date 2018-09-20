package main

import "strings"

// Capitalize returns input with the first character converted to upper case.
func Capitalize(input string) (result string) {

	// split the first character and the rest of the string
	first, rest := input[0:1], input[1:]

	// convert the first character to upper case
	first = strings.ToUpper(first)

	return first + rest
}

func Capitalize2(input string) (result string) {
	if len(input) == 0 {
		return input
	}
	first := strings.ToUpper(string(input[0]))
	return first + input[1:]
}

func Capitalize3(input string) (result string) {
	first := true
	for _, r := range input {
		if first {
			result += strings.ToUpper(string(r))
			first = false
			continue
		}

		result += string(r)
	}

	return result
}
