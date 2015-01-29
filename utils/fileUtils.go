package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadLines read lines from file and returns a string slice
func ReadLines(name string) ([]string, error) {
	inputFile, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("filename %s could not be opened", name)
	}

	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)
	if err != nil {
		return nil, err
	}
	// scanner.Scan() advances to the next token returning false if an error
	// was encountered
	l := newList()
	for scanner.Scan() {
		l.add(strings.Trim(scanner.Text(), "\n"))
	}

	// When finished scanning if any error other than io.EOF occured
	// it will be returned by scanner.Err().
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return l.toStringArray(), nil
}

// ReadAll reads the file and returns all in a hopefully
// newline stripped, single string.
func ReadAll(name, separator string) (string, error) {
	strSlice, err := ReadLines(name)
	if err != nil {
		return "", err
	}

	return strings.Join(strSlice, separator), nil
}
