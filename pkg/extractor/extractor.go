package extractor

import (
	"statsViewer/pkg/models"
	"strconv"
	"strings"
)

// Extract ...
type Extract struct {
	Line      string
	FileName  string
	Challenge *models.Challenge
}

var separator = ","

func (e *Extract) ExtractData() (err string) {
	// Malformed file names will not panic whole process
	defer func() {
		if fileError := recover(); fileError != nil {
			err = "Malformed file name"
		}
	}()

	e.name()
	e.score()
	e.sensScale()
	e.hsens()
	e.vsens()
	e.fov()
	e.DateAndTime()

	return err
}

func (e *Extract) name() {
	if strings.Contains(e.Line, "Scenario:,") {
		e.Challenge.Name = strings.Split(e.Line, separator)[1]
	}
}

func (e *Extract) sensScale() {
	if strings.Contains(e.Line, "Sens Scale:,") {
		e.Challenge.SensScale = strings.Split(e.Line, separator)[1]
	}
}

func (e *Extract) score() {
	if strings.Contains(e.Line, "Score:,") {
		scoreStr := strings.Split(e.Line, separator)[1]
		scoreFloat, _ := strconv.ParseFloat(scoreStr, 1)
		e.Challenge.Score = float64(int(scoreFloat*10)) / 10
	}
}

func (e *Extract) hsens() {
	if strings.Contains(e.Line, "Horiz Sens:,") {
		hsensStr := strings.Split(e.Line, separator)[1]
		hsensFloat, _ := strconv.ParseFloat(hsensStr, 1)
		e.Challenge.HSens = float64(int(hsensFloat*10)) / 10
	}
}

func (e *Extract) vsens() {
	if strings.Contains(e.Line, "Vert Sens:,") {
		vsensStr := strings.Split(e.Line, separator)[1]
		vsensFloat, _ := strconv.ParseFloat(vsensStr, 1)
		e.Challenge.VSens = float64(int(vsensFloat*10)) / 10
	}
}

func (e *Extract) fov() {
	if strings.Contains(e.Line, "FOV:,") {
		fovStr := strings.Split(e.Line, separator)[1]
		fovFloat, _ := strconv.ParseFloat(fovStr, 1)
		e.Challenge.FOV = float64(int(fovFloat*10)) / 10
	}
}

// DateAndTime ...
func (e *Extract) DateAndTime() {
	datetimeAndExtension := (strings.Split(e.FileName, " - Challenge - "))[1]
	datetime := strings.Split(datetimeAndExtension, " ")[0]
	dateAndTime := strings.Split(datetime, "-")

	// e.Challenge.Datetime = datetime
	e.Challenge.Date = dateAndTime[0]
	// e.Challenge.Time = dateAndTime[1]
}
