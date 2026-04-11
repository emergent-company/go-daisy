package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/emergent-company/go-daisy/cmd/gallery/internal/gallery"
	"github.com/emergent-company/go-daisy/galleryruntime"
)

func main() {
	// Wire GitHub App client if all required env vars are present.
	var ghCfg *galleryruntime.GitHubConfig
	appIDStr := os.Getenv("GITHUB_APP_ID")
	installIDStr := os.Getenv("GITHUB_APP_INSTALLATION_ID")
	keyFile := os.Getenv("GITHUB_APP_PRIVATE_KEY_FILE")
	repo := os.Getenv("GITHUB_REPO")
	if appIDStr != "" && installIDStr != "" && keyFile != "" && repo != "" {
		appID, errA := strconv.ParseInt(appIDStr, 10, 64)
		installID, errI := strconv.ParseInt(installIDStr, 10, 64)
		keyPEM, errK := os.ReadFile(keyFile)
		if errA != nil || errI != nil || errK != nil {
			log.Printf("warning: invalid GitHub App config: appID=%v installID=%v keyErr=%v — GitHub integration disabled", errA, errI, errK)
		} else {
			ghCfg = &galleryruntime.GitHubConfig{
				AppID:          appID,
				InstallationID: installID,
				PrivateKeyPEM:  string(keyPEM),
				Repo:           repo,
			}
		}
	} else {
		log.Printf("GitHub App integration disabled (set GITHUB_APP_ID, GITHUB_APP_INSTALLATION_ID, GITHUB_APP_PRIVATE_KEY_FILE, GITHUB_REPO to enable)")
	}

	dbPath := filepath.Join(os.TempDir(), "go-daisy-gallery.db")

	if err := galleryruntime.Serve(galleryruntime.Options{
		Title:      "go-daisy",
		Components: gallery.AllComponents(),
		Port:       11000,
		StorePath:  dbPath,
		GitHubCfg:  ghCfg,
		DevMode:    true,
	}); err != nil {
		log.Fatal(err)
	}
}
