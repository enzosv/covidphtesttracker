package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
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
	// date defaults to yesterday
	date := flag.String("d", time.Unix(time.Now().Unix()-86400, 0).Format("2006-01-02"), "Date to check")
	flag.Parse()
	if *date == "" || *configPath == "" {
		flag.PrintDefaults()
		return
	}
	config, err := parseConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = ProcessTesting(config, *date)
	if err != nil {
		log.Fatal(err)
	}
}

func ProcessTesting(config Config, date string) error {
	//pdf
	links, err := config.GDriveConfig.GetReadmeLinks(date)
	if err != nil {
		return err
	}

	// csv

	csvPath := fmt.Sprintf("%s.csv", date)
	defer os.Remove(csvPath)
	folder, _, err := config.GDriveConfig.GetTestFolderLink(links, csvPath)
	if err != nil {
		return err
	}

	// test
	test, err := readTest(csvPath, date)
	if err != nil {
		return err
	}
	if test.UniqueTested == 0 {
		return fmt.Errorf("No unique tests for date: %s\n", date)
	}
	// telegram
	// TODO: Consider moving message format to telegram config
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("Invalid date (%s) %w", date, err)
	}
	message := fmt.Sprintf("[%s](%s)\n*%.2f%%* positivity `(%.0f/%.0f)`", t.Format("Jan 2, 2006"), folder, test.Positive*100/test.UniqueTested, test.Positive, test.UniqueTested)
	log.Printf("\n\t[INFO] Sending Message:\n\t%s", message)

	sendMessage(config.TelegramConfig, message)
	return nil
}

// TODO: (Stretch) Download revisions and edit sent messages?
