package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"statsViewer/internal"
	"statsViewer/kovaaks"
	"text/template"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error:", err)
			EnterToExit()
		}
	}()

	start := time.Now()

	client, err := kovaaks.New()
	if err != nil {
		log.Printf("failed to initialize kovaaks client: %v", err)
		EnterToExit()
	}

	playlists, err := client.GetPlaylists()
	if err != nil {
		log.Printf("failed to get playlists: %v", err)
		EnterToExit()
	}

	playlist := getPlaylistFromUser(playlists)
	fmt.Println("Playlist selected:", playlist)

	stats, err := client.ParseStats(playlist)
	if err != nil {
		log.Printf("failed to parse stats: %v", err)
		EnterToExit()
	}
	fmt.Println("Files parsed. Creating HTML file...")

	// Output HTML
	t, err := template.ParseFiles("static/statsViewerTpl.html")
	Check(err)
	f, err := os.Create(internal.StatsViewerHTML)
	Check(err)
	err = t.Execute(f, stats)
	Check(err)
	f.Close()

	fmt.Println("Success!")
	fmt.Println(time.Now().Sub(start))
	exec.Command("cmd", "/C", "start", internal.StatsViewerHTML).Run()
}

// EnterToExit ...
func EnterToExit() {
	log.Println("Press \"enter\" key to exit.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(1)
}
