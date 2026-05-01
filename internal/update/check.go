package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/dimer47/action1-cli/internal/config"
)

const (
	owner         = "dimer47"
	repo          = "action1-cli"
	binary        = "action1"
	cacheDuration = 24 * time.Hour
	httpTimeout   = 3 * time.Second
	cacheFile     = "update-check.json"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

type updateCache struct {
	LatestVersion string    `json:"latest_version"`
	CheckedAt     time.Time `json:"checked_at"`
}

// Result holds the update check result.
type Result struct {
	CurrentVersion string
	LatestVersion  string
	HasUpdate      bool
}

// Check verifies if a newer version is available. Safe to call from a goroutine.
func Check(currentVersion string) *Result {
	if currentVersion == "" || currentVersion == "dev" {
		return nil
	}

	latest, err := getLatestVersion(currentVersion)
	if err != nil || latest == "" {
		return nil
	}

	if isNewer(latest, currentVersion) {
		return &Result{
			CurrentVersion: currentVersion,
			LatestVersion:  latest,
			HasUpdate:      true,
		}
	}

	return nil
}

// FormatMessage returns the update message for stderr.
func (r *Result) FormatMessage() string {
	if r == nil || !r.HasUpdate {
		return ""
	}

	goos := runtime.GOOS
	goarch := runtime.GOARCH

	var cmd string
	switch goos {
	case "darwin", "linux":
		ext := "tar.gz"
		cmd = fmt.Sprintf(
			"curl -sL https://github.com/%s/%s/releases/latest/download/%s_%s_%s.%s | tar xz && sudo mv %s /usr/local/bin/",
			owner, repo, repo, goos, goarch, ext, binary,
		)
	case "windows":
		cmd = fmt.Sprintf(
			"https://github.com/%s/%s/releases/latest/download/%s_%s_%s.zip",
			owner, repo, repo, goos, goarch,
		)
	default:
		cmd = fmt.Sprintf("https://github.com/%s/%s/releases/latest", owner, repo)
	}

	return fmt.Sprintf(
		"\n  Une nouvelle version de %s est disponible : %s → %s\n  Mise à jour :  %s\n",
		binary, r.CurrentVersion, r.LatestVersion, cmd,
	)
}

func getLatestVersion(currentVersion string) (string, error) {
	cachePath := filepath.Join(config.Dir(), cacheFile)

	// Try cache first
	if cached, err := loadCache(cachePath); err == nil {
		if time.Since(cached.CheckedAt) < cacheDuration {
			return cached.LatestVersion, nil
		}
	}

	// Fetch from GitHub
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github API returned %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	latest := strings.TrimPrefix(release.TagName, "v")

	// Save cache
	_ = saveCache(cachePath, updateCache{
		LatestVersion: latest,
		CheckedAt:     time.Now(),
	})

	return latest, nil
}

func loadCache(path string) (*updateCache, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c updateCache
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func saveCache(path string, c updateCache) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// isNewer returns true if latest > current (semver comparison).
func isNewer(latest, current string) bool {
	latest = strings.TrimPrefix(latest, "v")
	current = strings.TrimPrefix(current, "v")

	lParts := parseSemver(latest)
	cParts := parseSemver(current)

	for i := 0; i < 3; i++ {
		if lParts[i] > cParts[i] {
			return true
		}
		if lParts[i] < cParts[i] {
			return false
		}
	}
	return false
}

func parseSemver(v string) [3]int {
	var parts [3]int
	segments := strings.SplitN(v, ".", 3)
	for i, s := range segments {
		if i >= 3 {
			break
		}
		// Strip any pre-release suffix (e.g. "1-beta")
		s = strings.SplitN(s, "-", 2)[0]
		parts[i], _ = strconv.Atoi(s)
	}
	return parts
}
