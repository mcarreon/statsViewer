package kovaaks

import (
	"encoding/json"
	"fmt"
	"os"
	"statsViewer/internal"
	"statsViewer/pkg/models"
)

// GetPlaylists returns a list of local playlists exported from kovaaks
func (c *Client) GetPlaylists() ([]string, error) {
	playlistsPath := c.getPlaylistsPath()
	files, err := os.ReadDir(playlistsPath)
	if err != nil {
		return nil, err
	}

	playlists := make([]string, 0)
	for _, file := range files {
		playlists = append(playlists, stripJSONFileExtension(file.Name()))
	}

	playlists = append(playlists, internal.RecentScenarios)

	return playlists, nil
}

func (c *Client) getPlaylistInfo(playlistName string) (*models.Playlist, error) {
	if playlistName == internal.RecentScenarios {
		return nil, nil
	}

	data, err := os.ReadFile(fmt.Sprintf("%s/%s.json", c.getPlaylistsPath(), playlistName))
	if err != nil {
		return nil, err
	}

	var playlistInfo models.Playlist
	if err = json.Unmarshal(data, &playlistInfo); err != nil {
		return nil, err
	}

	return &playlistInfo, nil
}
