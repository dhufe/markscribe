package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/shurcooL/githubv4"
)

func TestGistFromQL(t *testing.T) {
	ts := time.Date(2024, 7, 10, 12, 0, 0, 0, time.UTC)
	in := qlGist{
		Name:        githubv4.String("example.md"),
		Description: githubv4.String("an example gist"),
		URL:         githubv4.String("https://gist.github.com/1"),
		CreatedAt:   githubv4.DateTime{Time: ts},
	}

	got := gistFromQL(in)
	want := Gist{
		Name:        "example.md",
		Description: "an example gist",
		URL:         "https://gist.github.com/1",
		CreatedAt:   ts,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("gistFromQL mismatch:\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRepoFromQL(t *testing.T) {
	in := qlRepository{
		NameWithOwner: githubv4.String("octo/repo"),
		URL:           githubv4.String("https://github.com/octo/repo"),
		Description:   githubv4.String("demo repo"),
		IsPrivate:     githubv4.Boolean(false),
	}
	in.Stargazers.TotalCount = githubv4.Int(42)

	got := repoFromQL(in)
	want := Repo{
		Name:        "octo/repo",
		URL:         "https://github.com/octo/repo",
		Description: "demo repo",
		Stargazers:  42,
		IsPrivate:   false,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("repoFromQL mismatch:\nwant: %#v\n got: %#v", want, got)
	}
}

func TestUserFromQL(t *testing.T) {
	in := qlUser{
		Login:     githubv4.String("octocat"),
		Name:      githubv4.String("The Octocat"),
		AvatarURL: githubv4.String("https://avatars.example/1.png"),
		URL:       githubv4.String("https://github.com/octocat"),
	}

	got := userFromQL(in)
	want := User{
		Login:     "octocat",
		Name:      "The Octocat",
		AvatarURL: "https://avatars.example/1.png",
		URL:       "https://github.com/octocat",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("userFromQL mismatch:\nwant: %#v\n got: %#v", want, got)
	}
}

func TestPullRequestFromQL(t *testing.T) {
	t0 := time.Date(2023, 3, 1, 8, 30, 0, 0, time.UTC)
	repo := qlRepository{
		NameWithOwner: githubv4.String("octo/repo"),
		URL:           githubv4.String("https://github.com/octo/repo"),
		Description:   githubv4.String("demo"),
		IsPrivate:     githubv4.Boolean(false),
	}
	repo.Stargazers.TotalCount = githubv4.Int(7)
	in := qlPullRequest{
		URL:        githubv4.String("https://github.com/octo/repo/pull/1"),
		Title:      githubv4.String("Add feature"),
		State:      githubv4.PullRequestState("OPEN"),
		CreatedAt:  githubv4.DateTime{Time: t0},
		Repository: repo,
	}

	got := pullRequestFromQL(in)
	want := PullRequest{
		Title:     "Add feature",
		URL:       "https://github.com/octo/repo/pull/1",
		State:     "OPEN",
		CreatedAt: t0,
		Repo: Repo{
			Name:        "octo/repo",
			URL:         "https://github.com/octo/repo",
			Description: "demo",
			Stargazers:  7,
			IsPrivate:   false,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("pullRequestFromQL mismatch:\nwant: %#v\n got: %#v", want, got)
	}
}

func TestReleaseFromQL(t *testing.T) {
	t0 := time.Date(2022, 1, 2, 3, 4, 5, 0, time.UTC)
	in := qlRelease{
		Nodes: []struct {
			Name         githubv4.String
			TagName      githubv4.String
			PublishedAt  githubv4.DateTime
			URL          githubv4.String
			IsPrerelease githubv4.Boolean
			IsDraft      githubv4.Boolean
		}{
			{
				Name:        githubv4.String("v1.0.0"),
				TagName:     githubv4.String("v1.0.0"),
				PublishedAt: githubv4.DateTime{Time: t0},
				URL:         githubv4.String("https://github.com/octo/repo/releases/tag/v1.0.0"),
			},
		},
	}

	got := releaseFromQL(in)
	want := Release{
		Name:        "v1.0.0",
		TagName:     "v1.0.0",
		PublishedAt: t0,
		URL:         "https://github.com/octo/repo/releases/tag/v1.0.0",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("releaseFromQL mismatch:\nwant: %#v\n got: %#v", want, got)
	}
}
