package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <username>\n", os.Args[0])
		os.Exit(1)
	}
	username := os.Args[1]

	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Fprintf(os.Stderr, "User not found\n")
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Error: %s\n", resp.Status)
		os.Exit(1)
	}

	var events []map[string]any
	err = json.NewDecoder(resp.Body).Decode(&events)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	
	PrintEvents(events)
}

func PrintEvents(events []map[string]any) {
	for _, event := range events {
		handler, ok := EventHanlders[event["type"].(string)]
		if ok {
			handler(event)
		}
	}
}

var EventHanlders = map[string]func(map[string]any){}

func init() {
	EventHanlders["PushEvent"] = PrintPushEvent
	EventHanlders["WatchEvent"] = PrintWatchEvent
	EventHanlders["PublicEvent"] = PrintPublicEvent
	EventHanlders["PullRequestEvent"] = PrintPullRequestEvent
}

func PrintWatchEvent(event map[string]any) {
	repo := event["repo"].(map[string]any)
	repoName := repo["name"].(string)
	payload := event["payload"].(map[string]any)
	action := payload["action"].(string)

	switch action {
	case "started":
		fmt.Printf("- Starred %s\n", repoName)
	case "stopped":
		fmt.Printf("- Unstarred %s\n", repoName)
	}
}

func PrintPublicEvent(event map[string]any) {
	repo := event["repo"].(map[string]any)
	repoName := repo["name"].(string)

	fmt.Printf("- Public repo %s\n", repoName)
}

func PrintPushEvent(event map[string]any) {
	repo := event["repo"].(map[string]any)
	repoName := repo["name"].(string)
	payload := event["payload"].(map[string]any)
	commits := payload["commits"].([]any)

	fmt.Printf("- Pushed %d commits to %s\n", len(commits), repoName)
}

func PrintPullRequestEvent(event map[string]any) {
	repo := event["repo"].(map[string]any)
	repoName := repo["name"].(string)
	payload := event["payload"].(map[string]any)
	action := payload["action"].(string)

	switch action {
	case "opened":
		fmt.Printf("- Pull request opened in %s\n", repoName)
	case "closed":
		fmt.Printf("- Pull request closed in %s\n", repoName)
	}
}
