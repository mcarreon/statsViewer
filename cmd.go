package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// getPlaylistFromUser returns the selected playlist from the user from the list of their local playlists
// as default functionality, if no playlists exist, recent scenarios will be displayed instead
func getPlaylistFromUser(playlists []string) string {
	if len(playlists) == 1 {
		return playlists[0]
	}

	var target string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		playlistStr := strings.Join(playlists, "\n")
		fmt.Printf("Enter the playlist to display stats for...valid choices are: \n%s\n", playlistStr)
		scanner.Scan()
		target = scanner.Text()

		if !strings.Contains(playlistStr, target) {
			fmt.Printf("\n%s is an invalid value. please select from an option above.", target)
		} else {
			break
		}
	}

	return target
}
