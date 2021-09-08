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
	"time"
)

const baseURL = "https://api.github.com"

var command, owner, repo, title, body string
var issue int

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
	Title string
	Body  string // in Markdown format
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
	var issue = NewIssue{Title: title, Body: body}
	url := baseURL + "/repos/" + owner + "/" + repo + "/issues"
	fmt.Printf("%s\n", url)
	data, err := json.Marshal(issue)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		resp.Body.Close()
		fmt.Printf("search query failed: %s", resp.Status)
	}
}

func read() {
	fmt.Printf("%s\n", flag.Args()[0:])
}

func update() {
	fmt.Printf("%s\n", flag.Args()[0:])
}

func lock() {
	fmt.Printf("%s\n", flag.Args()[0:])
}
