package internal

import (
	"os/exec"
)

func getOutputText(inputText, fontFile string) (string, error) {
	outputText, err := exec.Command("./internal/ascii-art-fs/ascii-art-fs", inputText, fontFile).Output()
	if err != nil {
		return "", err
	}
	return string(outputText), nil
}
