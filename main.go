package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GithubActivity struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-activity <username>")
		return
	}
	username := os.Args[1]
	if username == "" {
		fmt.Printf("Place username as an argument")
	}
	activity, err := api(username)
	if err == nil {

	}
	displayActivity(activity)
}

func api(username string) ([]GithubActivity, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events/public", username)
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	if response.StatusCode == 404 {
		return nil, fmt.Errorf("User not found.")
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var activity []GithubActivity
	if err := json.Unmarshal(body, &activity); err != nil {
		return nil, err
	}
	return activity, nil
}

func displayActivity(events []GithubActivity) {
	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			fmt.Printf("Pushed to %s\n", event.Repo.Name)
		case "IssuesEvent":
			fmt.Printf("Opened issue in %s\n", event.Repo.Name)
		case "WatchEvent":
			fmt.Printf("Starred %s\n", event.Repo.Name)
		case "PullRequestEvent":
			fmt.Printf("Opened PR in %s\n", event.Repo.Name)
		case "IssueCommentEvent":
			fmt.Printf("Commented in %s\n", event.Repo.Name)
		case "PullRequestReviewCommentEvent":
			fmt.Printf("Commented on PR in %s\n", event.Repo.Name)
		case "PullRequestReviewEvent":
			fmt.Printf("Reviewed PR in %s\n", event.Repo.Name)
		case "CreateEvent":
			fmt.Printf("Created a repository %s\n", event.Repo.Name)
		default:
			fmt.Printf("Unknown event: %s\n", event.Type)
		}
	}
}
