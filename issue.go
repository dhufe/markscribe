package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/shurcooL/githubv4"
)

var recentIssuesQuery struct {
	User struct {
		Login                   githubv4.String
		ContributionsCollection struct {
			IssueContributionsByRepository []struct {
				Contributions struct {
					Edges []struct {
						Cursor githubv4.String
						Node   struct {
							OccurredAt githubv4.DateTime
							Issue struct {
								Title githubv4.String
							}
						}
					}
				} `graphql:"contributions(first: 1)"`
				Repository qlRepository
			} `graphql:"issueContributionsByRepository(maxRepositories: 100)"`
		}
	} `graphql:"user(login:$username)"`
}

func recentIssues(count int) []Issue {

	var issues []Issue
	variables := map[string]interface{}{
		"username": githubv4.String(username),
	}
	err := gitHubClient.Query(context.Background(), &recentIssuesQuery, variables)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d issues!\n", len(recentIssuesQuery.User.ContributionsCollection.IssueContributionsByRepository))

	for _, v := range recentIssuesQuery.User.ContributionsCollection.IssueContributionsByRepository {
		// ignore meta-repo
		if string(v.Repository.NameWithOwner) == fmt.Sprintf("%s/%s", username, username) {
			continue
		}

		if v.Repository.IsPrivate {
			continue
		}

		c := Issue {
			Repo:       repoFromQL(v.Repository),
			OccurredAt: v.Contributions.Edges[0].Node.OccurredAt.Time,
			Title: 		string(v.Contributions.Edges[0].Node.Issue.Title),
		}

		fmt.Println(c)

		issues = append(issues, c) 
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].OccurredAt.After(issues[j].OccurredAt)
	})

	// fmt.Printf("Found %d issues!\n", len(repos))
	if len(issues) > count {
		return issues[:count]
	}
	return issues
}