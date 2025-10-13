package utils

import (
	"bufio"
	"fmt"
	"strings"

	"os"
)

type BruteForceType struct {
	WorldList string
	BruteForce bool
	Recurse bool
}

type AdminFindeType struct {
	WorldList string
	Exclude []string
}

func ReadFile(path string) []string {
	line, err := readLines(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return line	
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