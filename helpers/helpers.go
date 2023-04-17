package helpers

import "strings"

const (
	ShortQueryLen = 64
)

func Min(a, b int) int {
	if a <= b {
		return a
	}

	return b
}

func TruncateQuery(orig string) string {
	return orig[0:Min(ShortQueryLen, len(orig))]
}

func HostDropPort(host string) string {
	return strings.Split(host, ":")[0]
}
