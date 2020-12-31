package main

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

var client *hcloud.Client
var keepAmount int

func main() {
	// Get executable dir
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("error getting executable path: %s", err)
	}
	exPath := filepath.Dir(ex)

	// Load config and initialize
	err = godotenv.Load(".env", filepath.Join(exPath, ".env"))
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("error loading .env file(s): %s", reflect.TypeOf(err))
	}

	keepAmount, err = strconv.Atoi(os.Getenv("KEEP_AMOUNT"))
	if err != nil {
		log.Fatalf("environment variable KEEP_AMOUNT can't be converted to an integer: %s", err)
	}
	client = hcloud.NewClient(hcloud.WithToken(os.Getenv("HCLOUD_TOKEN")))

	// Loop servers and run backup for each
	servers, err := client.Server.All(context.Background())
	if err != nil {
		log.Fatalf("error retrieving servers: %s\n", err)
	}

	for _, server := range servers {
		if server.BackupWindow != "" {
			log.Printf("%s has backups enabled, skipping\n", server.Name)
			continue
		}

		// TODO: create backups async
		createBackup(server)
		pruneBackups(server)
	}
}

func createBackup(server *hcloud.Server) {
	log.Printf("creating backup for %s\n", server.Name)

	image, _, err := client.Server.CreateImage(context.Background(), server, &hcloud.ServerCreateImageOpts{
		Type:        hcloud.ImageTypeSnapshot,
		Description: hcloud.String(fmt.Sprintf("Backup %s %s", server.Name, time.Now().Format("Jan 02 2006 15:04 MST"))),
		Labels:      map[string]string{"autobackup": ""},
	})

	if err != nil {
		log.Fatalf("failed to create backup image: %s\n", err)
	}

	waitForAction(image.Action)
	log.Printf("successfully created image: %s\n", image.Image.Description)
}

func pruneBackups(server *hcloud.Server) {
	// Get all images with `autobackup` label, sorted by creation date (most recent first)
	images, err := client.Image.AllWithOpts(context.Background(), hcloud.ImageListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: "autobackup",
		},
		Sort: []string{"created:desc"},
		Type: []hcloud.ImageType{hcloud.ImageTypeSnapshot},
	})

	if err != nil {
		log.Fatalf("error retrieving images: %s\n", err)
	}

	backupCount := 0
	for _, image := range images {

		// Filter backups for current server only
		if image.CreatedFrom.ID == server.ID {
			backupCount++

			if backupCount > keepAmount {
				log.Printf("deleting backup image: %s", image.Description)
				_, err := client.Image.Delete(context.Background(), image)

				if err != nil {
					log.Fatalf("error deleting backup image (%s): %s", image.Description, err)
				}
			}
		}
	}
}

func waitForAction(action *hcloud.Action) {
	_, errors := client.Action.WatchProgress(context.Background(), action)

	err := <-errors // This blocks until finished/errored
	if err != nil {
		log.Fatalf("action %s failed: %s\n", action.Command, err)
	}
}
