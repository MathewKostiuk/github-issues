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
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const baseURL = "https://api.github.com"

var command, owner, repo, title, body string
var num int

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
	Title  string `json:"title"`
	Body   string `json:"body"` // in Markdown format
	Number int    `json:"id"`
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
	flag.IntVar(&num, "issue", 0, "The issue number (needed with read, update, and lock actions)")
	flag.Parse()
}

func main() {
	switch command {
	case "read":
		url := baseURL + "/repos/" + owner + "/" + repo + "/issues"
		read(url)
	case "create":
		url := baseURL + "/repos/" + owner + "/" + repo + "/issues"
		issue := NewIssue{Title: title, Body: body}
		create(url, issue)
	case "update":
		i := strconv.Itoa(num)
		url := baseURL + "/repos/" + owner + "/" + repo + "/issues/" + i
		var issue = NewIssue{Title: title, Body: body, Number: num}
		update(url, issue)
	case "lock":
		i := strconv.Itoa(num)
		url := baseURL + "/repos/" + owner + "/" + repo + "/issues/" + i + "/lock"
		lock(url)
	}
}

func create(url string, issue NewIssue) {
	github := auth()
	client := &http.Client{}
	var result *Issue
	data, err := json.Marshal(issue)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	resp := initRequest(req, github, client)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Request failed: %s\n", resp.Status)
	}

	issues := parseJSON(resp, result)
	printResponse(issues)
}

func read(url string) {
	github := auth()
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp := initRequest(req, github, client)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed: %s\n", resp.Status)
	}

	var issuesResult IssuesResult
	if err := json.NewDecoder(resp.Body).Decode(&issuesResult); err != nil {
		log.Fatal(err)
	}

	printResponse(issuesResult)
}

func update(url string, issue NewIssue) {
	github := auth()
	var result *Issue
	client := &http.Client{}

	data, err := json.Marshal(issue)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	resp := initRequest(req, github, client)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed: %s\n", resp.Status)
	}

	issues := parseJSON(resp, result)
	printResponse(issues)
}

func lock(url string) {
	github := auth()
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp := initRequest(req, github, client)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		fmt.Printf("Request failed: %s\n", resp.Status)
	}

	fmt.Printf("%s\n", resp.Status)
}

func printResponse(result IssuesResult) {
	for _, issue := range result {
		fmt.Printf("#%-5d %9s %.55s\n",
			issue.Number, issue.User.Login, issue.Title)
	}
}

func auth() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	githubKey := os.Getenv("GITHUB_TOKEN")

	return githubKey
}

func initRequest(req *http.Request, github string, client *http.Client) *http.Response {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer "+github)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func parseJSON(resp *http.Response, result *Issue) IssuesResult {
	var issuesResult IssuesResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}
	issuesResult = append(issuesResult, result)
	return issuesResult
}
