package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ChrisToast89/slate-for-windows/setup/internal/logx"
	"github.com/ChrisToast89/slate-for-windows/setup/internal/manifest"
	"github.com/ChrisToast89/slate-for-windows/setup/internal/paths"
)

type CheckResult struct {
	OK              bool   `json:"ok"`
	Installed       bool   `json:"installed"`
	InstallDir      string `json:"installDir"`
	InstalledTag    string `json:"installedTag"`
	LatestTag       string `json:"latestTag"`
	LatestURL       string `json:"latestUrl"`
	UpdateAvailable bool   `json:"updateAvailable"`
	Message         string `json:"message"`
	Error           string `json:"error,omitempty"`
	CheckedAt       string `json:"checkedAt"`
}

type ghRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Name    string `json:"name"`
}

// Check compares installed version to GitHub latest release for this repo.
func Check() CheckResult {
	res := CheckResult{CheckedAt: time.Now().UTC().Format(time.RFC3339)}
	dir, m, installed := manifest.Discover()
	res.Installed = installed
	res.InstallDir = dir
	if m != nil {
		res.InstalledTag = m.ReleaseTag
		if res.InstalledTag == "" {
			res.InstalledTag = m.AppVersion
		}
	}
	if !installed {
		res.InstalledTag = "(not installed)"
	}

	latest, err := fetchLatest()
	if err != nil {
		logx.Log("update check failed: %v", err)
		res.OK = false
		res.Error = err.Error()
		res.Message = "Could not reach GitHub to check for updates. You can still install the embedded app from this Setup package."
		// Bundled version is always available offline
		res.LatestTag = paths.AppVersion
		return res
	}
	res.OK = true
	res.LatestTag = latest.TagName
	res.LatestURL = latest.HTMLURL
	if res.LatestTag == "" {
		res.LatestTag = paths.AppVersion
	}

	bundled := paths.AppVersion
	// Prefer comparing against latest GitHub tag; also note if Setup bundle is behind.
	if installed && res.InstalledTag != "" && res.InstalledTag != "(not installed)" {
		res.UpdateAvailable = manifest.VersionLooksNewer(res.LatestTag, res.InstalledTag) ||
			manifest.VersionLooksNewer(bundled, res.InstalledTag)
	} else {
		res.UpdateAvailable = false
	}

	if !installed {
		res.Message = fmt.Sprintf("No install found. Latest release on GitHub: %s. This Setup bundles %s.", res.LatestTag, bundled)
	} else if res.UpdateAvailable {
		res.Message = fmt.Sprintf("Update available: installed %s → latest %s (this Setup also carries %s).", res.InstalledTag, res.LatestTag, bundled)
	} else {
		res.Message = fmt.Sprintf("Up to date (installed %s; GitHub latest %s).", res.InstalledTag, res.LatestTag)
	}
	logx.Log("update: installed=%s latest=%s available=%v", res.InstalledTag, res.LatestTag, res.UpdateAvailable)
	return res
}

func fetchLatest() (ghRelease, error) {
	var rel ghRelease
	client := &http.Client{Timeout: 12 * time.Second}
	req, err := http.NewRequest(http.MethodGet, paths.GitHubReleasesAPI, nil)
	if err != nil {
		return rel, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "SlateForWindows-Setup/"+paths.SetupVersion)
	resp, err := client.Do(req)
	if err != nil {
		return rel, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode == 404 {
		// No releases yet — treat package version as latest
		rel.TagName = "v" + paths.AppVersion
		rel.HTMLURL = paths.GitHubRepoURL + "/releases"
		return rel, nil
	}
	if resp.StatusCode != 200 {
		return rel, fmt.Errorf("GitHub API %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	if err := json.Unmarshal(body, &rel); err != nil {
		return rel, err
	}
	return rel, nil
}
