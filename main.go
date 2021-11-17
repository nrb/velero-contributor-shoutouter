package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/v32/github"
)

type args struct {
	debug                    bool
	config, token, team, org string
}

var a *args

func main() {

	c, err := makeConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Validate(); err != nil {
		log.Fatal(err)
	}

	p := NewProject(c)

	err = p.GetDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range p.Devs {
		if d != nil {
			fmt.Printf("Login: %s, Name: %s\n", d.GetLogin(), d.GetName())
		}
	}

	err = p.GetReposFromTeam()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("---- REPOS ----")
	for _, r := range p.Repos {
		fmt.Printf("%s\n", *r.Name)
	}
}

func init() {
	a = &args{}
	flag.StringVar(&a.config, "config", "", "--config file.yaml")
	flag.StringVar(&a.token, "token", "", "--token ${GITHUB_API_TOKEN}")
	flag.StringVar(&a.org, "org", "", "--org ${GITHUB_ORG_NAME}")
	flag.StringVar(&a.team, "team", "", "--team ${GITHUB_TEAM_NAME}")
	flag.BoolVar(&a.debug, "debug", false, "--debug for detailed logging")

	flag.Parse()

	if a.debug {
		log.SetLevel(log.DebugLevel)
	}
}

func makeConfig() (*ProjectConfig, error) {
	var err error
	c := &ProjectConfig{}

	if a != nil && a.config != "" {
		c, err = NewConfigFromFile(a.config)
	}
	if err != nil {
		return nil, err
	}

	if a != nil && a.org != "" {
		c.OrgName = a.org
	}

	if a.team != "" {
		c.TeamNames = append(c.TeamNames, a.team)
	}

	if a.token != "" {
		c.Token = a.token
	}

	return c, nil
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
