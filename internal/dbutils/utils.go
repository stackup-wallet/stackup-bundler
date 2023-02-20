package dbutils

import "strings"

const separator = ":"

func JoinValues(values ...string) string {
	return strings.Join(values, separator)
}

func SplitValues(value string) []string {
	return strings.Split(value, separator)
}
