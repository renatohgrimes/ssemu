package utils

import (
	"runtime"
	"strings"

	"zappem.net/pub/debug/xxd"
)

func StringIsAlphanumeric(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, c := range str {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

func EnvNewLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func BufferToString(buffer []byte) string {
	var sb strings.Builder
	lines := xxd.Dump(0, buffer)
	for _, line := range lines {
		sb.WriteString(line)
		sb.WriteString(EnvNewLine())
	}
	return sb.String()
}
