package charts

import (
	"fmt"
	charts2 "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"html/template"
	"sort"
	"statsViewer/pkg/models"
	"strconv"
	"strings"
)

// ScenarioLineChart ...
func ScenarioLineChart(scen *models.Scenario) template.HTML {
	line := charts2.NewLine()
	line.Renderer = newSnippetRenderer(line, line.Validate)
	line.SetGlobalOptions(chartGlobalOpts(max, scen.Name, len(scen.ByDateMax), scen.Highscore, scen.LowestAvgScore))

	var maxDates []string
	var maxScores []opts.LineData
	for _, dateScore := range scen.ByDateMax {
		for date, chall := range dateScore {
			maxDates = append(maxDates, simplifyDate(date))
			maxScores = append(maxScores, opts.LineData{
				Name:  fmt.Sprintf("%v: %v. FOV: %v. %v", simplifyDate(date), chall.Score, chall.FOV, chall.SensStr()),
				Value: chall.Score,
			})
		}
	}

	var avgDates []string
	var avgScores []opts.LineData
	for _, dateScore := range scen.ByDateAvg {
		for date, data := range dateScore {
			avgDates = append(avgDates, simplifyDate(date))
			avgScores = append(avgScores, opts.LineData{
				Name:  fmt.Sprintf("%v: %v. Grouped: %v", simplifyDate(date), data.Score, data.Grouped),
				Value: data.Score,
			})
		}
	}

	var wmaDates []string
	var wmaScores []opts.LineData
	for _, dateScore := range scen.ByDateWMA {
		for date, dateWMA := range dateScore {
			wmaDates = append(wmaDates, simplifyDate(date))
			wmaScores = append(wmaScores, opts.LineData{
				Name:  fmt.Sprintf("%v: %v. Grouped: %v", simplifyDate(date), dateWMA.Avg, dateWMA.Grouped),
				Value: dateWMA.Avg,
			})
		}
	}

	line.SetXAxis(maxDates).
		AddSeries("Max scores", maxScores).
		AddSeries("Average scores", avgScores).
		AddSeries(strconv.Itoa(defaultWMAWindow)+"-day moving average", wmaScores).
		SetSeriesOptions(seriesOpts...)

	return renderToHTML(line)
}

// PerformanceChart ...
func PerformanceChart(uniqueDays *map[string]int) template.HTML {
	progress := charts2.NewLine()
	progress.Renderer = newSnippetRenderer(progress, progress.Validate)

	// Order map by date
	var orderedDates []map[string]int
	for k, v := range *uniqueDays {
		orderedDates = append(orderedDates, map[string]int{k: v})
	}
	sort.SliceStable(orderedDates, func(i, j int) bool {
		var iDate int
		for k := range orderedDates[i] {
			iDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		var jDate int
		for k := range orderedDates[j] {
			jDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		return iDate < jDate
	})

	lowestAvgPb := 99
	var dates []string
	var avgPercentagePBs []opts.LineData
	for _, dateAndAvgPercentagePB := range orderedDates {
		for date, avgPercentagePB := range dateAndAvgPercentagePB {
			if avgPercentagePB <= 0 {
				continue
			}
			if lowestAvgPb > avgPercentagePB {
				lowestAvgPb = avgPercentagePB
			}
			dates = append(dates, simplifyDate(date))
			avgPercentagePBs = append(avgPercentagePBs, opts.LineData{
				Name:  simplifyDate(date) + " " + strconv.Itoa(avgPercentagePB) + "%",
				Value: avgPercentagePB,
			})
		}
	}

	progress.SetGlobalOptions(
		charts2.WithTitleOpts(opts.Title{
			Title:    "Experimental performance tracker",
			Subtitle: "Data points are average scores for every scenario played that day, converted into a percentage of your current highscore.",
		}),
		charts2.WithYAxisOpts(opts.YAxis{
			Type:      "value",
			Max:       100,
			Min:       10 * (lowestAvgPb / 10),
			AxisLabel: &yAxisLabelFormatter,
		}),
		charts2.WithTooltipOpts(opts.Tooltip{
			Trigger:   "axis",
			TriggerOn: "mousemove|click",
			Show:      true,
			Formatter: "{b}",
		}),
		ToolBoxOpts("performance"),
		xAxisOpts,
		initOpts,
	)

	progress.SetXAxis(dates).
		AddSeries("PB %", avgPercentagePBs).
		SetSeriesOptions(seriesOpts...)

	return renderToHTML(progress)
}

var max = "max"
var avg = "avg"

var seriesOpts = []charts2.SeriesOpts{
	charts2.WithLabelOpts(opts.Label{Show: true, Color: "black"}),
	charts2.WithLineChartOpts(opts.LineChart{Smooth: true}),
}

func titleOpts(scenName string, length int) charts2.GlobalOpts {
	return charts2.WithTitleOpts(opts.Title{
		Title:    scenName,
		Subtitle: fmt.Sprintf("Grouped by day, %v datapoints.", length),
	})
}

var dataZoom = opts.ToolBoxFeatureDataZoom{
	Show:  true,
	Title: map[string]string{"zoom": "Zoom", "back": "Restore"},
}

func saveAsImage(fileName string) *opts.ToolBoxFeatureSaveAsImage {
	return &opts.ToolBoxFeatureSaveAsImage{
		Show:  true,
		Name:  fileName,
		Title: "Download",
	}
}

func toolBoxFeatures(fileName string) *opts.ToolBoxFeature {
	return &opts.ToolBoxFeature{
		SaveAsImage: saveAsImage(fileName),
		DataZoom:    &dataZoom,
	}
}

// ToolBoxOpts ...
func ToolBoxOpts(fileName string) charts2.GlobalOpts {
	return charts2.WithToolboxOpts(opts.Toolbox{
		Show:    true,
		Feature: toolBoxFeatures(fileName),
	})
}

var tooltipOpts = charts2.WithTooltipOpts(opts.Tooltip{
	Trigger:   "item",
	TriggerOn: "mousemove|click",
	Show:      true,
	Formatter: "{b}",
})

var xAxisLabelFormatter = opts.AxisLabel{
	Rotate: 45,
}

var xAxisOpts = charts2.WithXAxisOpts(opts.XAxis{
	AxisLabel: &xAxisLabelFormatter,
})

var yAxisLabelFormatter = opts.AxisLabel{}

func yAxisOpts(highscore, lowestAvg float64) charts2.GlobalOpts {
	return charts2.WithYAxisOpts(opts.YAxis{
		Type:      "value",
		Max:       10 * ((int(highscore*1.05) + 9) / 10),
		Min:       10 * (int(lowestAvg*0.95) / 10),
		AxisLabel: &yAxisLabelFormatter,
	})
}

var legendOpts = charts2.WithLegendOpts(opts.Legend{
	Show: true,
})

var initOpts = charts2.WithInitializationOpts(opts.Initialization{
	AssetsHost: "static/",
})

func chartGlobalOpts(groupedBy string, scenName string, length int, hs float64, ls float64) (charts2.GlobalOpts, charts2.GlobalOpts, charts2.GlobalOpts, charts2.GlobalOpts, charts2.GlobalOpts, charts2.GlobalOpts, charts2.GlobalOpts) {
	return titleOpts(scenName, length), ToolBoxOpts(scenName), tooltipOpts, xAxisOpts, yAxisOpts(hs, ls), legendOpts, initOpts
}

// simplifyDate ...
func simplifyDate(d string) string {
	sep := "/"
	d = strings.ReplaceAll(d, ".", sep)
	d = reorderDate(d, sep)
	d = strings.ReplaceAll(d, "2018", "18")
	d = strings.ReplaceAll(d, "2019", "19")
	d = strings.ReplaceAll(d, "2020", "20")
	d = strings.ReplaceAll(d, "/2021", "")

	return d
}

func reorderDate(d, sep string) string {
	dateUnits := strings.Split(d, sep)
	dateUnits[0], dateUnits[1], dateUnits[2] = dateUnits[1], dateUnits[2], dateUnits[0]

	return strings.Join(dateUnits, sep)
}
