package data

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

// Returns all the content of a file
func ReadAllFileContent(path string) (content []byte, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ReadAllFileContent: %w", err)
	}
	defer f.Close()
	content, err = io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("ReadAllFileContent: %w", err)
	}
	return content, nil
}

// Create a new file with a content written in it
// Note that an append write feature is not implemented
func WriteFile(path string, content []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("WriteFile: %w", err)
	}
	defer f.Close()
	f.Write(content)
	return nil
}

// examine the content of a file in the format like this:
// -----------------
// key  value...
// tree 01234example
// blob testtest
// ref  sample
// .
// .
// -----------------
// If there exists a given key in the file, then returns the corresponding value.
func ReadValueFromFile(path string, key []byte) (value []byte, err error) {
	c, err := ReadAllFileContent(path)
	if err != nil {
		return nil, fmt.Errorf("ReadValueFromFile: %w", err)
	}
	if len(c) == 0 {
		return nil, nil
	}
	sc := bufio.NewScanner(bytes.NewReader(c))
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		l := sc.Bytes()
		sep := bytes.Split(l, []byte(" "))
		if len(sep) < 2 {
			return nil, fmt.Errorf("ReadValueFromFile: invalid data: got=%s", sep)
		}
		if bytes.Equal(sep[0], key) {
			return sep[1], nil
		}
	}
	return nil, nil
}
