package kovaaks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"statsViewer/internal"
)

// populateKovaaksPath returns the kovaaks path from the config file, or the default path if config file is not provided
func (c *Client) populateKovaaksPath() (err error) {
	config, err := os.ReadFile(internal.DefaultConfigPath)
	if err != nil {
		return err
	}

	var parsedConfig map[string]interface{}
	if err = json.Unmarshal(config, &parsedConfig); err != nil {
		return err
	}

	kovaaksPath := filepath.Clean(parsedConfig[internal.ConfigKeyKovaaksPath].(string))
	if kovaaksPath == "" {
		kovaaksPath = internal.DefaultKovaaksPath
	}

	c.kovaaksPath = kovaaksPath

	return nil
}

// getStatsPath returns the path to the stats folder
func (c *Client) getStatsPath() (statsPath string) {
	return filepath.Clean(c.kovaaksPath + "/stats/")
}

// getPlaylistsPath returns the path to the playlists folder
func (c *Client) getPlaylistsPath() (playlistsPath string) {
	return filepath.Clean(c.kovaaksPath + "/Saved/SaveGames/Playlists/")
}
