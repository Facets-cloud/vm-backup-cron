package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/gorhill/cronexpr"
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

	var cronSchedule string
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

	for {
		newCronSchedule := os.Getenv("CRON_SCHEDULE")
		if newCronSchedule == "" {
			log.Fatal("CRON_SCHEDULE environment variable is not set")
		}

		if newCronSchedule != cronSchedule {
			cronSchedule = newCronSchedule
			c.Stop()
			c = cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

			_, err := c.AddFunc(cronSchedule, backup)
			if err != nil {
				log.Fatalf("Failed to parse CRON_SCHEDULE: %v", err)
			}

			// Parse the cron expression to get the next run time
			expr, err := cronexpr.Parse(cronSchedule)
			if err != nil {
				log.Fatalf("Failed to parse CRON_SCHEDULE: %v", err)
			}
			nextTime := expr.Next(time.Now())
			log.Printf("Cron schedule updated to: %s (next run at: %s)", cronSchedule, nextTime.Format(time.RFC1123))
			c.Start()
		}
		
		time.Sleep(10 * time.Second)
	}
}
