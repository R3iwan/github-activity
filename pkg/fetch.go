package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Event struct {
	Actor   Actor   `json:"actor"`
	Repo    Repo    `json:"repo"`
	Payload Payload `json:"payload"`
}

type Actor struct {
	Login string `json:"login"`
	URL   string `json:"url"`
}

type Repo struct {
	Name string `json:"name"`
}

type Payload struct {
	Commits []Commit `json:"commits"`
}

type Commit struct {
	Message string `json:"message"`
}

func FetchJSON() {
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error fetching a username")
	}
	username = strings.TrimSpace(username)

	link := fmt.Sprintf("https://api.github.com/users/%s/events", username)
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println("Error fetching a link")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error receiving a StatusOK")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading a body")
	}

	var events []Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Println("Error parsing a JSON")
		return
	}

	Repos := make(map[string]bool)

	fmt.Println("--Output--")
	for _, event := range events {
		if !Repos[event.Repo.Name] {
			fmt.Printf("Login: %s\n", event.Actor.Login)
			fmt.Printf("Repo Name: %s\n", event.Repo.Name)
			fmt.Printf("URL: %s\n\n", event.Actor.URL)
			Repos[event.Repo.Name] = true
		}

		if len(event.Payload.Commits) > 0 {
			fmt.Println("Commits:")
			for _, commit := range event.Payload.Commits {
				fmt.Printf("  - %s\n", commit.Message)
			}
			fmt.Println()
		}
	}
}
