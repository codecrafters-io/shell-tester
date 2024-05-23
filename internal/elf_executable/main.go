package elf_executable

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/tester-utils/random"
)

var inputFiles = []string{
	"./internal/elf_executable/0_elf_header.hex",
	"./internal/elf_executable/1_program_header.hex",
	"./internal/elf_executable/2_program_code.hex",
	"./internal/elf_executable/3_string_table.hex",
	"./internal/elf_executable/4_section_header_0.hex",
	"./internal/elf_executable/5_section_header_1.hex",
	"./internal/elf_executable/6_section_header_2.hex",
}

func CreateELFexecutable(randomString string, outputFile string) error {
	var concatenatedData []byte
	for _, inputFile := range inputFiles {
		binaryData, err := getBinaryDataFromHexFile(inputFile, randomString)
		if err != nil {
			return fmt.Errorf("CodeCrafters internal error. Unable to read from ELF constituent files: %v", err)
		}
		concatenatedData = append(concatenatedData, binaryData...)
	}

	err := os.WriteFile(outputFile, concatenatedData, 0755)
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

func readFile(inputFile string) ([]byte, error) {
	hexFile, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}

	// Remove any whitespace characters (e.g., newlines) from the hex data
	hexStr := ""
	scanner := bufio.NewScanner(hexFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		line = removeWhitespace(line)
		hexStr += line
	}

	// Decode the hex string to binary data
	binaryData, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return binaryData, nil
}

func getBinaryDataFromHexFile(inputFile string, randomString string) ([]byte, error) {
	hexData, err := readFile(inputFile)
	if err != nil {
		return nil, err
	}

	// Here, we add our random output to the ELF file's program code section
	if strings.Contains(inputFile, "program_code") {
		randomString += "\n"
		hexRandomString := hex.EncodeToString([]byte(randomString))
		hexData = append(hexData, []byte(hexRandomString)...)
	}

	return hexData, nil
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
