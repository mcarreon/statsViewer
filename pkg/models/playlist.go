package models

// Playlist is a struct for playlist data in kovaaks
type Playlist struct {
	PlaylistName        string             `json:"playlistName"`
	PlaylistId          int                `json:"playlistId"`
	AuthorSteamId       string             `json:"authorSteamId"`
	AuthorName          string             `json:"authorName"`
	ScenarioList        []PlaylistScenario `json:"scenarioList"`
	Description         string             `json:"description"`
	HasOfflineScenarios bool               `json:"hasOfflineScenarios"`
	HasEdited           bool               `json:"hasEdited"`
	ShareCode           string             `json:"shareCode"`
	Version             int                `json:"version"`
	Updated             int64              `json:"updated"`
	IsPrivate           bool               `json:"isPrivate"`
}

// PlaylistScenario is a struct for a scenario in a playlist
type PlaylistScenario struct {
	ScenarioName string `json:"scenario_name"`
	PlayCount    int    `json:"play_Count"`
}

func (p *Playlist) GetScenarioList() []string {
	scenarios := make([]string, 0)
	for _, scenario := range p.ScenarioList {
		scenarios = append(scenarios, scenario.ScenarioName)
	}

	return scenarios
}
