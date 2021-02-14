package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type TestingRow struct {
	UniqueTested float64
	Positive     float64
}

type TestingColumns int

const (
	Date         TestingColumns = 1
	UniqueTested                = 4
	Positive                    = 5
	Status                      = 17
)

type Config struct {
	TelegramConfig TelegramConfig `json:"telegram"`
	GDriveConfig   GDriveConfig   `json:"gdrive"`
}

func main() {
	configPath := flag.String("c", "config.json", "Config file")
	date := flag.String("d", "", "Date to check")
	link := flag.String("l", "", "Link to source")
	flag.Parse()
	if *date == "" || *link == "" || *configPath == "" {
		flag.PrintDefaults()
		return
	}
	config, err := parseConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	csvPath := fmt.Sprintf("testaggregates/%s.csv", *date)
	err = downloadCSV(config.GDriveConfig, *link, csvPath)
	if err != nil {
		log.Fatal(err)
	}
	test, err := readTest(csvPath, *date)
	if err != nil {
		log.Fatal(err)
	}
	if test.UniqueTested == 0 {
		log.Fatalf("No unique tests for date: %s\n", *date)
	}
	// TODO: Consider moving message format to telegram config
	message := fmt.Sprintf("[%s](%s): *%.2f%%* positivity (%.0f/%.0f)", *date, *link, test.Positive*100/test.UniqueTested, test.Positive, test.UniqueTested)
	fmt.Println(message)

	sendMessage(config.TelegramConfig, message)
}

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
			fmt.Printf("blank unique with %s positive\n", row[Positive])
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

// TODO: (Stretch) Download revisions and edit sent messages?
