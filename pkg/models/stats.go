package models

import (
	"html/template"
)

// Stats ...
type Stats struct {
	PlaylistName      string
	Scenarios         map[string]*Scenario
	SortedTimesPlayed []*Scenario
	UniqueDays        map[string]int
	DaysPlayed        int
	TotalScens        int
	TotalPlayed       int
	PerformanceChart  template.HTML
	WordCloud         template.HTML
}

func NewStats(playlistName string) *Stats {
	return &Stats{
		PlaylistName: playlistName,
		Scenarios:    make(map[string]*Scenario),
		UniqueDays:   make(map[string]int),
	}
}
