package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	// Function to execute the backup
	backup := func() {
		cmd := exec.Command("/vmbackup-prod",
			"-storageDataPath", os.Getenv("STORAGE_DATA_PATH"),
			"-dst", os.Getenv("DESTINATION"),
			"-snapshot.createURL", os.Getenv("SNAPSHOT_CREATE_URL"))

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Backup attempt failed: %v\nOutput: %s", err, string(output))
		} else {
			log.Println("Victoria metrics backup process completed successfully.")
		}
	}

	var ticker *time.Ticker
	var cronSchedule string

	for {
		newCronSchedule := os.Getenv("CRON_SCHEDULE")
		if newCronSchedule == "" {
			log.Fatal("CRON_SCHEDULE environment variable is not set")
		}

		if newCronSchedule != cronSchedule {
			cronSchedule = newCronSchedule
			if ticker != nil {
				ticker.Stop()
			}

			duration, err := time.ParseDuration(cronSchedule)
			if err != nil {
				log.Fatalf("Failed to parse CRON_SCHEDULE: %v", err)
			}

			ticker = time.NewTicker(duration)
			log.Printf("Cron schedule updated to: %s", cronSchedule)
			log.Printf("Next backup will occur in: %s", duration.String())
		}

		select {
		case <-ticker.C:
			backup()
		case <-time.After(10 * time.Second):
			// Check for updated cron schedule every 10 seconds
		}
	}
}
