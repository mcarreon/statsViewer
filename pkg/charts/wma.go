package charts

import "statsViewer/pkg/models"

// TODO: Make this configurable
const defaultWMAWindow = 7

type WeightedMovingAverage struct {
	N int
	// Values in each group receive the same weight
	groups [][]float64
}

func NewWMA(n int) *WeightedMovingAverage {
	return &WeightedMovingAverage{
		N: n,
	}
}

func (wma *WeightedMovingAverage) Add(values ...float64) {
	n := wma.N
	if len(wma.groups) >= n {
		wma.groups = wma.groups[1:n]
	}

	wma.groups = append(wma.groups, values)
}

func (wma *WeightedMovingAverage) Average() (float64, int) {
	groups := wma.groups
	if groups == nil {
		return 0, 0
	}

	weight := wma.N - len(groups) + 1
	weightTotal := 0
	sum := float64(0)
	count := 0
	for _, group := range groups {
		for _, value := range group {
			sum += value * float64(weight)
			weightTotal += weight
			count++
		}
		weight++
	}

	avg := sum / float64(weightTotal)
	return avg, count
}

func CalculateWMA(scenario *models.Scenario) {
	wma := NewWMA(defaultWMAWindow)
	for _, dateScores := range scenario.ByDateScores {
		for date, scores := range dateScores {
			wma.Add(scores...)
			avg, count := wma.Average()
			dateWMA := models.DateWMA{Avg: float64(int(avg*10)) / 10, Grouped: count}
			scenario.ByDateWMA = append(scenario.ByDateWMA, map[string]models.DateWMA{date: dateWMA})
		}
	}
}
