package security

import (
	"bufio"
	"os"
	"strings"
)

func LoadBlacklistFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var rules []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		rules = append(rules, line)
	}

	return rules, scanner.Err()
}
