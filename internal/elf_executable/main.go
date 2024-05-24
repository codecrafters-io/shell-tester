package elf_executable

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/tester-utils/random"
)

func CreateELFexecutable(randomString string, outputFile string) error {
	var concatenatedData []byte

	// part1 contains the elf_header, program_header and program_code
	// part2 contains the string_table and section_headers
	part1, part2 := getELFHexData()
	part1Bytes, err := parseData(part1)
	if err != nil {
		return fmt.Errorf("CodeCrafters internal error. Failed to parse ELF data: %v", err)
	}
	part2Bytes, err := parseData(part2)
	if err != nil {
		return fmt.Errorf("CodeCrafters internal error. Failed to parse ELF data: %v", err)
	}

	concatenatedData = append(concatenatedData, part1Bytes...)
	concatenatedData = append(concatenatedData, []byte(randomString)...)
	concatenatedData = append(concatenatedData, part2Bytes...)

	err = os.WriteFile(outputFile, concatenatedData, 0755)
	if err != nil {
		return fmt.Errorf("CodeCrafters internal error. Failed to write output ELF file: %v", err)
	}

	return nil
}

func GetRandomString() string {
	// We will use a random numeric string of length = 12
	var result string
	for i := 0; i < 6; i++ {
		result += fmt.Sprintf("%d", random.RandomInt(10, 99))
	}

	return result
}

func parseData(inputData string) ([]byte, error) {
	lines := strings.Split(inputData, "\n")
	hexStr := ""

	for _, line := range lines {
		line = removeWhitespace(line)
		if !strings.HasPrefix(line, "#") {
			hexStr += line
		}
	}

	// Decode the hex string to binary
	binaryData, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return binaryData, nil
}

func removeWhitespace(str string) string {
	result := make([]rune, 0, len(str))
	for _, r := range str {
		if !isWhitespace(r) {
			result = append(result, r)
		}
	}
	return string(result)
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t' || r == '\r'
}
