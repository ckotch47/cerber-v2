package utils

import (
	"bufio"
	"errors"
	"strings"

	"os"
)

type BruteForceType struct {
	WorldList   string
	Recurse     bool
	MaxDepth    int
	Concurrency int
}

type AdminFindeType struct {
	WorldList      string
	Exclude        []string
	Timeout        int
	RequestTimeout int
}

func ReadFile(path string) ([]string, error) {
	line, err := readLines(path)
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, errors.New("файл пустой")
	}
	return line, nil
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") { // пропускаем пустые и комментарии
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
