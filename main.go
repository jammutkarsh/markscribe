package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/KyleBanks/goodreads"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var (
	gitHubClient    *githubv4.Client
	goodReadsClient *goodreads.Client
	goodReadsID     string
	username        string

	write = flag.String("write", "", "write output to")
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: markscribe [template]")
		os.Exit(1)
	}

	tplIn, err := os.ReadFile(flag.Args()[0])
	if err != nil {
		fmt.Println("Can't read file:", err)
		os.Exit(1)
	}

	tpl, err := template.New("tpl").Funcs(template.FuncMap{
		/* GitHub */
		"recentContributions": recentContributions,
		"recentPullRequests":  recentPullRequests,
		"recentRepos":         recentRepos,
		"recentForks":         recentForks,
		"recentReleases":      recentReleases,
		"followers":           recentFollowers,
		"recentStars":         recentStars,
		"gists":               gists,
		"sponsors":            sponsors,
		"repo":                repo,
		/* RSS */
		"rss": rssFeed,
		/* GoodReads */
		"goodReadsReviews":          goodReadsReviews,
		"goodReadsCurrentlyReading": goodReadsCurrentlyReading,
		/* Literal.club */
		"literalClubCurrentlyReading": literalClubCurrentlyReading,
		/* Utils */
		"humanize": humanized,
		"reverse":  reverse,
		"now":      time.Now,
		"contains": strings.Contains,
		"toLower":  strings.ToLower,
	}).Parse(string(tplIn))
	if err != nil {
		fmt.Println("Can't parse template:", err)
		os.Exit(1)
	}

	var httpClient *http.Client
	gitHubToken := os.Getenv("GITHUB_TOKEN")
	goodReadsToken := os.Getenv("GOODREADS_TOKEN")
	goodReadsID = os.Getenv("GOODREADS_USER_ID")
	if len(gitHubToken) > 0 {
		httpClient = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: gitHubToken},
		))
	}

	gitHubClient = githubv4.NewClient(httpClient)
	goodReadsClient = goodreads.NewClient(goodReadsToken)

	if len(gitHubToken) > 0 {
		username, err = getUsername()
		if err != nil {
			fmt.Println("Can't retrieve GitHub profile:", err)
			os.Exit(1)
		}
	}

	// render fully first, so a failing template never truncates the output file
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, nil); err != nil {
		fmt.Println("Can't render template:", err)
		os.Exit(1)
	}

	if len(*write) == 0 {
		_, _ = buf.WriteTo(os.Stdout)
		return
	}
	if err := os.WriteFile(*write, buf.Bytes(), 0o644); err != nil {
		fmt.Println("Can't write:", err)
		os.Exit(1)
	}
}
