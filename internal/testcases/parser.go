package testcases

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
)

type InputField struct {
	name  string
	value string
}
type TestCase struct {
	ActionUID              string
	InputFields            []InputField
	ExpectedExecutionLabel string
}

func Parse(filePath string) ([]TestCase, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parse(file)
}

func parse(data io.Reader) ([]TestCase, error) {
	reader := csv.NewReader(data)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, errors.New("No test cases provided")
	}

	result := make([]TestCase, len(records))
	for i, record := range records {
		c := TestCase{}
		for j, col := range record {
			switch headers[j] {
			case "actionUID":
				c.ActionUID = col
			case "expectedExecutionLabel":
				c.ExpectedExecutionLabel = col
			default:
				field := InputField{
					name:  headers[j],
					value: col,
				}
				c.InputFields = append(c.InputFields, field)
			}
		}
		if c.ActionUID == "" {
			return nil, errors.New("No actionUID provided")
		}
		result[i] = c
	}
	return result, nil
}
