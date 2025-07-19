package utils

import (
	"bufio"
	"fmt"
	"os"
)

func LeerText() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("entrada vacía o EOF")
	}

	text := scanner.Text()
	return text, nil
}
