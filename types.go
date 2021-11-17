package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

// ProjectConfig is a representation of an umbrella project to thank contributors on.
// A project can span one or more repostories within one GitHub organization.
// Values within a Project are strings used to query the GitHub API.
type ProjectConfig struct {
	// Name of the project to find contributors.
	Name string
	// Name of the GitHub organization which the repositories are in.
	OrgName string
	// List of Repositories that should be considered within the project.
	RepoNames []string
	// List of core developer accounts names that will be excluded. Will be merged with teams.
	DevNames []string
	// GitHub Teams to extract members from that will be exluded. Members will be merged with developers.
	TeamNames []string
	// Token is a GitHub Token that is able to read the GitHub repositories.
	Token string
}

type Project struct {
	Config *ProjectConfig
	// PullRequests are the merged pull requests for the project.
	PullRequests []*github.PullRequest
	// Devs are GitHub Users that will be excluded from consideration.
	Devs []*github.User
	// Teams are GitHub Teams that will be inspected for Users to exclude from consideration.
	Teams []*github.Team
	// Repos is a list of GitHub Repositories associated with the team(s)
	Repos []*github.Repository
	// Client is a GitHub client used to query the REST API
	client *github.Client
}

func NewProject(c *ProjectConfig) *Project {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Project{
		Config: c,
		client: github.NewClient(tc),
	}
}

// GetPullRequests will get all pull requests across all of a project's repositories.
func (p *Project) GetPullRequests() error {
	opts := &github.PullRequestListOptions{State: "closed"}
	for _, r := range p.Config.RepoNames {
		prs, _, err := p.client.PullRequests.List(context.Background(), p.Config.OrgName, r, opts)
		if err != nil {
			return err
		}
		p.PullRequests = append(p.PullRequests, prs...)
	}
	return nil
}

// GetDevs will get all GitHub Users based on the configuration, as well Users from the specified Teams.
func (p *Project) GetDevs() error {
	// devsSeen holds a list of devs we've already seen based on a team, so we don't get duplicates.
	var devsSeen []string

	// Get all the team members first
	teamOpts := &github.TeamListTeamMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for _, t := range p.Config.TeamNames {
		log.Debugf("Trying to get members for team %s\n", t)
		devs, response, err := p.client.Teams.ListTeamMembersBySlug(context.Background(), p.Config.OrgName, t, teamOpts)

		log.Debugf("Response: %v\n", response)
		if err != nil {
			return err
		}
		for _, d := range devs {
			log.Debugf("Found user: %v\n", d)
			p.Devs = append(p.Devs, d)
			devsSeen = append(devsSeen, *d.Login)
		}
	}

	// Get all devs, if they're not already in the list.
	for _, n := range p.Config.DevNames {
		// Skip over any devs we've already seen.
		for _, ds := range devsSeen {
			if ds == n {
				log.Debugf("Skipping user %s, already present", n)
				continue
			}
		}

		log.Debugf("Trying to get user %s", n)
		dev, response, err := p.client.Users.Get(context.Background(), n)
		log.Debugf("Response: %v\n", response)
		if err != nil {
			return err
		}
		p.Devs = append(p.Devs, dev)
	}

	return nil
}

// GetReposFromTeam will retrieve the repositories that all the teams are responsible for.
func (p *Project) GetReposFromTeam() error {
	opts := &github.ListOptions{}
	for _, t := range p.Config.TeamNames {
		log.Debugf("Trying to get repos for team %s\n", t)
		repos, response, err := p.client.Teams.ListTeamReposBySlug(context.Background(), p.Config.OrgName, t, opts)
		log.Debugf("Response: %v\n", response)

		if err != nil {
			return err
		}
		p.Repos = append(p.Repos, repos...)
	}
	return nil
}
