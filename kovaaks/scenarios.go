package kovaaks

import (
	"fmt"
	"os"
	"slices"
	"statsViewer/internal"
	"statsViewer/pkg/models"
	"strings"
)

func (c *Client) getScenarioFiles(playlistInfo *models.Playlist) ([]os.FileInfo, error) {
	var (
		err        error
		dirEntries = make([]os.DirEntry, 0)
	)

	dirEntries, err = c.getScenarios()
	if err != nil {
		return nil, err
	}

	files := make([]os.FileInfo, 0)
	for _, entry := range dirEntries {
		if entry.IsDir() == true {
			continue
		}

		// if recent scenarios was selected, skip filtering scenarios out
		if playlistInfo.PlaylistName != internal.RecentScenarios && !isValidScenario(playlistInfo.GetScenarioList(), entry.Name()) {
			continue
		}

		file, fileErr := entry.Info()
		if fileErr != nil {
			fmt.Printf("Error getting file info for %s, cause: %s", entry.Name(), fileErr)
			continue
		}

		files = append(files, file)
	}

	return files, nil
}

func (c *Client) getScenarios() ([]os.DirEntry, error) {
	files, err := os.ReadDir(c.getStatsPath())
	if err != nil {
		return nil, err
	}

	return files, nil
}

func isValidScenario(validScenarios []string, scenarioName string) bool {
	return slices.Contains(validScenarios, trimScenarioName(scenarioName))
}

func trimScenarioName(scenarioName string) string {
	index := strings.Index(scenarioName, internal.ScenariosPostfix)
	if index == -1 {
		return scenarioName
	}

	return scenarioName[:index]
}
