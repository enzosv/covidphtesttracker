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
	// date defaults to 2 days ago
	date := flag.String("d", time.Unix(time.Now().Unix()-172800, 0).Format("2006-01-02"), "Date to check")
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
	link, err := config.GDriveConfig.GetTestFolderLink(links, csvPath)
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
	message := fmt.Sprintf("[%s](%s)\n*%.2f%%* positivity (%.0f/%.0f)", date, link, test.Positive*100/test.UniqueTested, test.Positive, test.UniqueTested)
	fmt.Println(message)

	sendMessage(config.TelegramConfig, message)
	return nil
}

// TODO: (Stretch) Download revisions and edit sent messages?
