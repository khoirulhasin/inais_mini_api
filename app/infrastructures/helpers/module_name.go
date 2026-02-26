package helpers

import (
	"bufio"
	"os"
	"strings"
)

func GetModuleRootName() string {
	file, err := os.Open("go.mod")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}

	return ""
}
