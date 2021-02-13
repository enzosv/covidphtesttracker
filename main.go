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
}

func main() {
	telegramPath := flag.String("tc", "telegram.json", "Telegram config file")
	testPath := flag.String("ta", "", "Testing Aggregates csv file")
	date := flag.String("d", "", "Date to check")
	link := flag.String("l", "", "Link to source") // TODO: Download csv file straight from this link
	flag.Parse()
	if *date == "" || *link == "" || *testPath == "" {
		flag.PrintDefaults()
		return
	}
	test, err := readTest(*testPath, *date)
	if err != nil {
		log.Fatal(err)
	}
	if test.UniqueTested == 0 {
		log.Fatalf("No unique tests for date: %s\n", *date)
	}
	// TODO: Consider moving message format to telegram config
	message := fmt.Sprintf("[%s](%s): `%.2f%%` positivity `(%.0f/%.0f)`\n", *date, *link, test.Positive*100/test.UniqueTested, test.Positive, test.UniqueTested)
	fmt.Println(message)
	telegramConfig := parseConfig(*telegramPath)
	sendMessage(telegramConfig, message)
}

func parseConfig(path string) TelegramConfig {
	telegramConfiguration := TelegramConfig{}
	configFile, err := os.Open(path)
	if err != nil {
		log.Fatal("Cannot open telegram configuration file: ", err)
	}
	defer configFile.Close()
	dec := json.NewDecoder(configFile)
	if err = dec.Decode(&telegramConfiguration); errors.Is(err, io.EOF) {
		//do nothing
	} else if err != nil {
		log.Fatal("Cannot load telegram configuration file: ", err)
	}
	return telegramConfiguration
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
			log.Printf("read error: %v\n", err)
			continue
		}
		date := row[Date]
		if date != checkDate {
			continue
		}
		unique, err := strconv.ParseFloat(row[UniqueTested], 64)
		if err != nil {
			log.Printf("scan unique (%s) error: %v\n", row[UniqueTested], err)
			continue
		}
		test.UniqueTested += unique
		positive, err := strconv.ParseFloat(row[Positive], 64)
		if err != nil {
			log.Printf("scan positive (%s) error: %v\n", row[Positive], err)
		}
		test.Positive += positive

		// REVIEW: Maybe consider "For Review" status. Discard?
	}
	return test, nil
}

// TODO: (Stretch) Download revisions and edit sent messages?
