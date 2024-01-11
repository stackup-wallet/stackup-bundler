package filter

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	userOpHashStrLen = 64
	userOpHashRegex  = regexp.MustCompile(fmt.Sprintf("(?i)0x[0-9a-f]{%d}", userOpHashStrLen))
)

func IsValidUserOpHash(userOpHash string) bool {
	return len(strings.TrimPrefix(strings.ToLower(userOpHash), "0x")) == userOpHashStrLen &&
		userOpHashRegex.MatchString(userOpHash)
}
