package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v32/github"
)

func main() {
	client := github.NewClient(nil)

	// Only our repos.
	repos := []string{
		"velero",
		"velero-plugin-for-aws",
		"velero-plugin-for-gcp",
		"velero-plugin-for-microsoft-azure",
	}

	var prs []*github.PullRequest

	// Get them all.
	for _, r := range repos {
		p, err := getPRs(client, r)
		if err != nil {
			log.Fatalf("Error getting PRs for repo %s: %s", r, err.Error())
		}
		prs = append(prs, p...)

	}

	// Filter out any PRs by us.
	prs = filterPRsByAuthor(prs)

	// Only show PRs merged in the last week.
	prs = filterMergedPRs(prs)

	// Print them.
	for _, pr := range prs {
		printShoutout(pr)
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
			// Use 10 days as a buffer for long weekends
			if days <= 10 {
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
