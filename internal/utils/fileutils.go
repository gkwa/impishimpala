package utils

import (
	"bufio"
	"bytes"
	"os"
	"unicode/utf8"
)

// IsTextFile determines if a file contains text content
func IsTextFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read first 512 bytes to determine file type
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && n == 0 {
		return false, err
	}

	// Check if content is valid UTF-8 and doesn't contain null bytes
	buffer = buffer[:n]

	// Check for null bytes (common in binary files)
	if bytes.Contains(buffer, []byte{0}) {
		return false, nil
	}

	// Check if content is valid UTF-8
	return utf8.Valid(buffer), nil
}

// ReadFileLines reads a file and returns its lines
func ReadFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
