package kovaaks

import (
	"bufio"
	"fmt"
	"github.com/gosuri/uiprogress"
	"os"
	"sort"
	"statsViewer/internal"
	"statsViewer/pkg/charts"
	"statsViewer/pkg/extractor"
	"statsViewer/pkg/models"
	"strconv"
	"strings"
	"sync"
)

func (c *Client) ParseStats(playlistName string) (*models.Stats, error) {
	fmt.Println("Parsing files... This may take a few minutes!")
	playlistInfo, err := c.getPlaylistInfo(playlistName)
	if err != nil {
		return nil, nil
	}
	fmt.Printf("%s info parsed, retrieving relevant scenario files...\n", playlistName)

	files, err := c.getScenarioFiles(playlistInfo)
	if err != nil {
		return nil, err
	}

	bar := uiprogress.AddBar(len(files)).AppendCompleted().PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Files (%d/%d)", b.Current(), len(files))
	})
	uiprogress.Start()

	stats := models.NewStats(playlistName)

	mux := &sync.Mutex{}
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)

		go func() {
			defer wg.Done()
			fileWorker(stats, file, mux, bar, c.getStatsPath())
		}()
	}
	wg.Wait()
	uiprogress.Stop()

	handleScenarios(stats)

	percentagesToPBByDate(stats)
	stats.WordCloud = charts.WordCloud(&stats.Scenarios)
	stats.PerformanceChart = charts.PerformanceChart(&stats.UniqueDays)
	stats.DaysPlayed = len(stats.UniqueDays)
	stats.SortedTimesPlayed = sortScenarios(stats.SortedTimesPlayed, playlistInfo)

	return stats, nil
}

func fileWorker(stats *models.Stats, file os.FileInfo, mux *sync.Mutex, bar *uiprogress.Bar, statsPath string) {
	// Open file
	f, err := os.Open(fmt.Sprintf("%s/%s", statsPath, file.Name()))
	if err != nil {
		fmt.Printf("Error opening file %s, cause: %s", file.Name(), err)
		return
	}
	defer f.Close()

	// New challenge
	challenge := models.Challenge{}

	s := bufio.NewScanner(f)
	// For each line
	for s.Scan() {
		line := s.Text()
		if s.Err() != nil {
			fmt.Printf("Error scanning file %s, cause: %s", file.Name(), err)
			return
		}

		extract := extractor.Extract{Line: line, FileName: file.Name(), Challenge: &challenge}
		extract.ExtractData()
	}

	if valid := challenge.IsValid(); valid == false {
		return
	}

	mux.Lock()
	stats.TotalPlayed++
	if _, ok := stats.Scenarios[challenge.Name]; ok {
		stats.Scenarios[challenge.Name].TimesPlayed++
		stats.Scenarios[challenge.Name].Challenges = append(stats.Scenarios[challenge.Name].Challenges, challenge)
	} else {
		stats.TotalScens++
		stats.Scenarios[challenge.Name] = &models.Scenario{
			FileName:    file.Name(),
			Name:        challenge.Name,
			TimesPlayed: 1,
			Challenges:  []models.Challenge{challenge},
		}
	}
	mux.Unlock()
	bar.Incr()
}

func handleScenarios(stats *models.Stats) {
	mux := &sync.Mutex{}
	var wg sync.WaitGroup
	for _, scenario := range stats.Scenarios {
		wg.Add(1)
		go func() {
			defer wg.Done()
			scenarioWorker(scenario, &stats.SortedTimesPlayed, &stats.UniqueDays, mux)
		}()
	}
	wg.Wait()
}

func scenarioWorker(scenario *models.Scenario, sortedTimesPlayed *[]*models.Scenario, uniqueDays *map[string]int, mux *sync.Mutex) {
	scenario.Lowscore = scenario.Challenges[0].Score
	ByDate := map[string][]models.Challenge{}

	for _, challenge := range scenario.Challenges {
		if challenge.Score > scenario.Highscore {
			scenario.Highscore = challenge.Score
		}
		if challenge.Score < scenario.Lowscore {
			scenario.Lowscore = challenge.Score
		}

		ByDate[challenge.Date] = append(ByDate[challenge.Date], challenge)
	}

	// Group challenges per date
	groupedMax, groupedAvg, groupedScore := group(ByDate, scenario.Highscore)

	// Maps into a slice so we can sort them
	for date, challenge := range groupedMax {
		scenario.ByDateMax = append(scenario.ByDateMax, map[string]models.Challenge{date: challenge})
	}
	scenario.LowestAvgScore = scenario.Highscore
	for date, dateAvg := range groupedAvg {
		scenario.ByDateAvg = append(scenario.ByDateAvg, map[string]models.DateAvg{date: dateAvg})
		if dateAvg.Score < scenario.LowestAvgScore {
			scenario.LowestAvgScore = dateAvg.Score
		}
	}
	for date, scores := range groupedScore {
		scenario.ByDateScores = append(scenario.ByDateScores, map[string][]float64{date: scores})
	}

	// Actually sort by date (descending)
	sort.SliceStable(scenario.ByDateMax, func(i, j int) bool {
		var iDate int
		for k := range scenario.ByDateMax[i] {
			iDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		var jDate int
		for k := range scenario.ByDateMax[j] {
			jDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		return iDate < jDate
	})
	sort.SliceStable(scenario.ByDateAvg, func(i, j int) bool {
		var iDate int
		for k := range scenario.ByDateAvg[i] {
			iDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		var jDate int
		for k := range scenario.ByDateAvg[j] {
			jDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		return iDate < jDate
	})
	sort.SliceStable(scenario.ByDateScores, func(i, j int) bool {
		var iDate int
		for k := range scenario.ByDateScores[i] {
			iDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		var jDate int
		for k := range scenario.ByDateScores[j] {
			jDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		return iDate < jDate
	})

	charts.CalculateWMA(scenario)

	mux.Lock()
	defer mux.Unlock()
	*sortedTimesPlayed = append(*sortedTimesPlayed, scenario)
	// Less than 2 datapoints or 3 challenges => skip chart
	if scenario.TimesPlayed <= 2 || len(ByDate) <= 1 {
		return
	}
	scenario.ChartByDate = charts.ScenarioLineChart(scenario)
}

func group(challengesByDate map[string][]models.Challenge, scenarioHighScore float64) (map[string]models.Challenge, map[string]models.DateAvg, map[string][]float64) {
	ByDateMax := map[string]models.Challenge{}
	ByDateAvg := map[string]models.DateAvg{}
	ByDateAll := map[string][]float64{}

	for date, challenges := range challengesByDate {
		challengesAmount := len(challenges)
		var maxScore float64
		var maxChallange models.Challenge
		var avgScore float64
		var sum float64
		var allScores []float64

		for i, challenge := range challenges {
			if i == 0 || challenge.Score > maxScore {
				maxScore = challenge.Score
				maxChallange = challenge
			}
			sum += challenge.Score
			allScores = append(allScores, challenge.Score)
		}

		avgScore = sum / float64(challengesAmount)

		ByDateMax[date] = maxChallange
		ByDateAvg[date] = models.DateAvg{
			Score:        float64(int(avgScore*10)) / 10,
			Grouped:      challengesAmount,
			PercentagePB: int((avgScore / scenarioHighScore) * 100),
		}
		ByDateAll[date] = allScores
	}

	return ByDateMax, ByDateAvg, ByDateAll
}

func percentagesToPBByDate(s *models.Stats) {
	// Each date will have a list of scenario averages
	dates := make(map[string][]int)

	for _, scenario := range s.Scenarios {
		for _, dateAndDateAvg := range scenario.ByDateAvg {
			for date, DateAvg := range dateAndDateAvg {
				dates[date] = append(dates[date], DateAvg.PercentagePB)
			}
		}
	}

	for date, percentages := range dates {
		amount := len(percentages)
		if amount <= 5 {
			continue
		}

		var sum int
		for _, v := range percentages {
			sum += v
		}
		avgPercentage := sum / amount

		s.UniqueDays[date] = avgPercentage
	}
}

func sortScenarios(timesPlayed []*models.Scenario, playlistInfo *models.Playlist) []*models.Scenario {
	sortedScenarios := timesPlayed

	if playlistInfo.PlaylistName == internal.RecentScenarios {
		sort.SliceStable(sortedScenarios, func(i, j int) bool {
			return sortedScenarios[i].TimesPlayed > sortedScenarios[j].TimesPlayed
		})
	} else {
		sortedScenarios = make([]*models.Scenario, 0)
		for _, scenario := range playlistInfo.GetScenarioList() {
			for _, timesPlayedScenario := range timesPlayed {
				if scenario == timesPlayedScenario.Name {
					sortedScenarios = append(sortedScenarios, timesPlayedScenario)
				}
			}
		}

	}

	return sortedScenarios
}
