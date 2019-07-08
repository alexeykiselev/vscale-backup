package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	vscale "github.com/vscale/go-vscale"
)

const (
	appName = "vscale-backup"
)

var version = "v0.0.0"

func main() {
	log.Printf("Starting %s %s", appName, version)
	var (
		token      = flag.String("token", "", "Vscale API token")
		expiration = flag.String("expiration", "336h", "Backups expiration time, older backups will be removed")
	)
	flag.Parse()

	interval, err := time.ParseDuration(*expiration)
	if err != nil {
		log.Fatalf("Invalid backup expiration '%s': %v", *expiration, err)
	}
	backupValidTill := time.Now().Add(-interval)

	c := vscale.NewClient(*token)
	servers, _, err := c.Scalet.List()
	if err != nil {
		log.Fatalf("Failed to retrieve servlets: %v", err)
	}

	backupDate := time.Now().Format("2006-01-02")
	for _, s := range *servers {
		backupName := fmt.Sprintf("%s_%s", s.Name, backupDate)
		_, _, err := c.Scalet.Backup(s.CTID, backupName)
		if err != nil {
			log.Printf("Error creating backup of servlet '%s' named '%s': %v", s.Name, backupName, err)
			continue
		}
		log.Printf("Backup '%s' of servlet '%s' successfully created", backupName, s.Name)
	}
	backups, _, err := c.Backup.List()
	if err != nil {
		log.Fatalf("Failed to retrieve backups: %v", err)
	}
	for _, b := range *backups {
		if b.Status != "finished" {
			log.Printf("Skipping backup '%s': not finished", b.Name)
			continue
		}
		createdAt, err := time.ParseInLocation("02.01.2006 15:04:05", b.Created, time.Local)
		if err != nil {
			log.Printf("Skipping backup '%s': invalid creation time: %v", b.Name, err)
			continue
		}
		if createdAt.After(backupValidTill) {
			log.Printf("Skipping backup '%s': not expired", b.Name)
			continue
		}
		_, _, err = c.Backup.Remove(b.ID)
		if err != nil {
			log.Printf("Error removing backup '%s': %v", b.Name, err)
			continue
		}
		log.Printf("Backup '%s' successfully deleted", b.Name)
	}
}
