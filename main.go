// githubissues is a tool that lets users create, read, update,
// and close GitHub issues from the command line.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const baseURL = "https://api.github.com"

var command, owner, repo, title, body string
var issue int

type IssuesResult []*Issue

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type NewIssue struct {
	Title string `json:"title"`
	Body  string `json:"body"` // in Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

func init() {
	flag.StringVar(&command, "command", "read", "Used to determine what action to take")
	flag.StringVar(&owner, "owner", "", "Owner of the respository")
	flag.StringVar(&repo, "repo", "", "Repository to search")
	flag.StringVar(&title, "title", "", "The title of the issue")
	flag.StringVar(&body, "body", "", "The body of the issue")
	flag.IntVar(&issue, "issue", 0, "The issue number (needed with read, update, and lock actions)")
	flag.Parse()
}

func main() {
	switch command {
	case "read":
		read()
	case "create":
		create()
	case "update":
		update()
	case "lock":
		lock()
	}
}

func create() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	githubKey := os.Getenv("GITHUB_TOKEN")

	client := &http.Client{}
	var issue = NewIssue{Title: title, Body: body}
	var result *Issue
	var issuesResult IssuesResult

	url := baseURL + "/repos/" + owner + "/" + repo + "/issues"
	data, err := json.Marshal(issue)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer "+githubKey)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Request failed: %s\n", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("json failed: %s\n", err)
		log.Fatal(err)
	}
	issuesResult = append(issuesResult, result)
	printResponse(issuesResult)
}

func read() {
	url := baseURL + "/repos/" + owner + "/" + repo + "/issues"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
	}

	var result IssuesResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		log.Fatal(err)
	}
	resp.Body.Close()
	printResponse(result)
}

func update() {
	fmt.Printf("%s\n", flag.Args()[0:])
}

func lock() {
	fmt.Printf("%s\n", flag.Args()[0:])
}

func printResponse(result IssuesResult) {
	for _, issue := range result {
		fmt.Printf("#%-5d %9s %.55s\n",
			issue.Number, issue.User.Login, issue.Title)
	}
}
