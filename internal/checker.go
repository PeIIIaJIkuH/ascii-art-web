package internal

import (
	"strings"
)

func isPrintableChar(inputStr string) bool {
	inputStr = strings.ToLower(inputStr)
	flag := true
	for _, value := range inputStr {
		flag = flag && ((int(value) > 32 && int(value) < 132) || int(value) == 134 || value == '\n' || value == ' ' || value == '\r')
	}
	return flag
}

func checkFs(filename string) bool {
	if !(filename == "standard.txt" || filename == "shadow.txt" || filename == "thinkertoy.txt") {
		return false
	}
	return true
}
