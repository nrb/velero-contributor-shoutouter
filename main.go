package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/v32/github"
)

func main() {
	token := flag.String("token", "", "--token {GITHUB_TOKEN}")
	debug := flag.Bool("debug", false, "--debug to print debug lines")

	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if *token == "" {
		fmt.Println("Please provide a GitHub API token with --token.")
		os.Exit(1)
	}

	// TODO: ensure we have an org name
	tce := NewProject(
		&ProjectConfig{
			Name:    "tce",
			OrgName: "vmware-tanzu",
			TeamNames: []string{
				"tce-owners",
			},
			Token: *token,
		},
	)

	err := tce.GetDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range tce.Devs {
		if d != nil {
			fmt.Printf("Login: %s, Name: %s\n", d.GetLogin(), d.GetName())
		}
	}

	err = tce.GetReposFromTeam()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("---- REPOS ----")
	for _, r := range tce.Repos {
		fmt.Printf("%s\n", *r.Name)
	}
}

func getPRs(client *github.Client, repo string) ([]*github.PullRequest, error) {
	orgName := "vmware-tanzu"
	prOpts := &github.PullRequestListOptions{State: "closed"}
	prs, _, err := client.PullRequests.List(context.Background(), orgName, repo, prOpts)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

// filterPRsByAuthor returns PRs by authors who were not part of the core maintainers team.
func filterPRsByAuthor(prs []*github.PullRequest) []*github.PullRequest {
	authors := []string{
		"ashish-amarnath",
		"carlisia",
		"jonasrosland",
		"michmike",
		"nrb",
		"zubron",
		"dsu-igeek",
	}

	var filtered []*github.PullRequest

	for _, pr := range prs {
		found := false
		for _, a := range authors {
			if *pr.User.Login == a {
				found = true
			}
		}

		if !found {
			filtered = append(filtered, pr)
		}
	}
	return filtered
}

// filterMergedPRs returns only PRs that have been merged within the last week,
// not PRs that have been closed without merging.
func filterMergedPRs(prs []*github.PullRequest) []*github.PullRequest {
	var filtered []*github.PullRequest

	for _, pr := range prs {
		if pr.MergedAt != nil {
			diff := time.Now().Sub(*pr.MergedAt)
			days := diff.Hours() / 24
			if days <= 7 {
				filtered = append(filtered, pr)
			}
		}

	}
	return filtered
}

// Formats the shoutouts in a consistent format.
func printShoutout(pr *github.PullRequest) {
	fmt.Printf("- [@%s](%s): [%s](%s)\n", *pr.User.Login, *pr.User.HTMLURL, *pr.Title, *pr.HTMLURL)
}
