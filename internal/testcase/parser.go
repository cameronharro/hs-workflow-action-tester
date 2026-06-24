package testcase

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type TestCase struct {
	ActionUID              string
	InputFields            map[string]any
	ObjectID               int
	ObjectType             string
	PortalID               int
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
		c := TestCase{
			InputFields: map[string]any{},
		}
		for j, col := range record {
			switch headers[j] {
			case "actionUID":
				c.ActionUID = col
			case "expectedExecutionLabel":
				c.ExpectedExecutionLabel = col
			case "objectId":
				n, err := strconv.Atoi(col)
				if err != nil {
					return nil, fmt.Errorf("Invalid objectId: %s", col)
				}

				c.ObjectID = n
			case "portalId":
				n, err := strconv.Atoi(col)
				if err != nil {
					return nil, fmt.Errorf("Invalid portalId: %s", col)
				}

				c.PortalID = n
			default:
				c.InputFields[headers[j]] = col
			}
		}
		if c.ActionUID == "" {
			return nil, errors.New("No actionUID provided")
		}
		result[i] = c
	}
	return result, nil
}
