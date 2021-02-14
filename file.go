package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ledongthuc/pdf"
)

func parseConfig(path string) (Config, error) {
	configuration := Config{}
	configFile, err := os.Open(path)
	if err != nil {
		return configuration, fmt.Errorf("Cannot open configuration file: %w", err)
	}
	defer configFile.Close()
	dec := json.NewDecoder(configFile)
	if err = dec.Decode(&configuration); errors.Is(err, io.EOF) {
		//do nothing
	} else if err != nil {
		return configuration, fmt.Errorf("Cannot load configuration file: %w", err)
	}
	return configuration, nil
}

func readTest(filepath, checkDate string) (TestingRow, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return TestingRow{}, fmt.Errorf("open error: %w", err)
	}
	defer file.Close()
	r := csv.NewReader(file)

	// skip first line
	r.Read()

	test := TestingRow{0, 0}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("read error: %v\n", err)
			continue
		}
		date := row[Date]
		if date != checkDate {
			continue
		}
		if row[UniqueTested] == "" {
			continue
		}
		unique, err := strconv.ParseFloat(row[UniqueTested], 64)
		if err != nil {
			fmt.Printf("scan unique (%s) error: %v. Positive: %s\n", row[UniqueTested], err, row[Positive])
			continue
		}
		test.UniqueTested += unique
		positive, err := strconv.ParseFloat(row[Positive], 64)
		if err != nil {
			fmt.Printf("scan positive (%s) error: %v\n", row[Positive], err)
			continue
		}
		test.Positive += positive

		// REVIEW: Maybe consider "For Review" status. Discard?
	}
	return test, nil
}

func readPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return "", fmt.Errorf("pdf: open error: %w", err)
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("pdf: read error: %w", err)
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

func findSubstringInPDF(path, substring string) ([]string, error) {
	text, err := readPDF(path)
	if err != nil {
		return nil, err
	}
	words := strings.Split(text, " ")
	var matches []string
	for _, w := range words {
		if strings.Contains(w, substring) {
			matches = append(matches, w)
		}
	}
	return matches, nil
}
