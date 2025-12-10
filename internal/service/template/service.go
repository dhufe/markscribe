package template

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	texttmpl "text/template"
	"time"

	kbgoodreads "github.com/KyleBanks/goodreads"
	"github.com/KyleBanks/goodreads/responses"
	"github.com/dustin/go-humanize"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	githubadapter "hufschlaeger.net/markscribe/internal/adapters/github"
	goodreadsadapter "hufschlaeger.net/markscribe/internal/adapters/goodreads"
	literaladapter "hufschlaeger.net/markscribe/internal/adapters/literal"
	rssadapter "hufschlaeger.net/markscribe/internal/adapters/rss"
	domain "hufschlaeger.net/markscribe/internal/domain"
	githubsvc "hufschlaeger.net/markscribe/internal/service/github"
	goodreadssvc "hufschlaeger.net/markscribe/internal/service/goodreads"
	literalsvc "hufschlaeger.net/markscribe/internal/service/literal"
	rsssvc "hufschlaeger.net/markscribe/internal/service/rss"
)

// Service composes per-port services and exposes template-facing API.
type Service struct {
	gh  *githubsvc.Service
	gr  *goodreadssvc.Service
	lit *literalsvc.Service
	rss *rsssvc.Service
}

func New(gh *githubsvc.Service, gr *goodreadssvc.Service, lit *literalsvc.Service, rss *rsssvc.Service) *Service {
	return &Service{gh: gh, gr: gr, lit: lit, rss: rss}
}

// NewFromEnv wires all dependencies based on environment variables and returns a ready-to-use Service.
// This consolidates startup logic so callers (like cmd/markscribe) can remain lean.
func NewFromEnv(ctx context.Context) (*Service, error) {
	// Tokens and settings
	gitHubToken := os.Getenv("GITHUB_TOKEN")
	goodReadsToken := os.Getenv("GOODREADS_TOKEN")
	goodReadsID := os.Getenv("GOODREADS_USER_ID")

	// Optional authenticated HTTP client for GitHub
	var httpClient *http.Client
	if len(gitHubToken) > 0 {
		httpClient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: gitHubToken},
		))
	}

	// External clients
	ghClient := githubv4.NewClient(httpClient)
	grClient := kbgoodreads.NewClient(goodReadsToken)

	// Adapters
	ghPort := githubadapter.New(ghClient)
	grPort := goodreadsadapter.New(grClient, goodReadsID)
	litPort := literaladapter.New()
	rssPort := rssadapter.New()

	// Username is only available with a token; non-fatal if missing
	username := ""
	if len(gitHubToken) > 0 {
		var err error
		username, err = ghPort.ViewerLogin(ctx)
		if err != nil {
			return nil, fmt.Errorf("can't retrieve GitHub profile: %w", err)
		}
	}

	// Services
	ghSvc := githubsvc.New(ghPort, username)
	grSvc := goodreadssvc.New(grPort)
	litSvc := literalsvc.New(litPort)
	rssSvc := rsssvc.New(rssPort)

	return New(ghSvc, grSvc, litSvc, rssSvc), nil
}

// GitHub
func (s *Service) RecentRepos(count int) []domain.Repo { return s.gh.RecentRepos(count) }
func (s *Service) RecentForks(count int) []domain.Repo { return s.gh.RecentForks(count) }
func (s *Service) Repo(owner, name string) domain.Repo { return s.gh.Repo(owner, name) }
func (s *Service) Followers(count int) []domain.User   { return s.gh.Followers(count) }
func (s *Service) RecentPullRequests(count int) []domain.PullRequest {
	return s.gh.RecentPullRequests(count)
}
func (s *Service) RecentReleases(count int) []domain.Repo { return s.gh.RecentReleases(count) }
func (s *Service) RecentContributions(count int) []domain.Contribution {
	return s.gh.RecentContributions(count)
}
func (s *Service) Gists(count int) []domain.Gist         { return s.gh.Gists(count) }
func (s *Service) RecentStars(count int) []domain.Star   { return s.gh.RecentStars(count) }
func (s *Service) RecentIssues(count int) []domain.Issue { return s.gh.RecentIssues(count) }
func (s *Service) Sponsors(count int) []domain.Sponsor   { return s.gh.Sponsors(count) }

// GoodReads
func (s *Service) GoodReadsReviews(count int) []responses.Review { return s.gr.Reviews(count) }
func (s *Service) GoodReadsCurrentlyReading(count int) []responses.Review {
	return s.gr.CurrentlyReading(count)
}

// Literal.club
func (s *Service) LiteralCurrentlyReading(count int) []literaladapter.Book {
	return s.lit.CurrentlyReading(count)
}

// LatestRssFeeds RSS
func (s *Service) LatestRssFeeds(url string, count int) []domain.RSSEntry {
	return s.rss.LastFeedEntries(url, count)
}

// Utils (moved from root template.go to declutter main package)
func (s *Service) Humanize(t interface{}) string {
	switch v := t.(type) {
	case time.Time:
		// flatten time to prevent updating README too often
		v = time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location())
		if time.Since(v) <= time.Hour*24 {
			return "today"
		}
		return humanize.Time(v)
	default:
		return fmt.Sprintf("%v", t)
	}
}

func (s *Service) Reverse(slc interface{}) interface{} {
	n := reflect.ValueOf(slc).Len()
	swap := reflect.Swapper(slc)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
	return slc
}

// Funcs returns the template FuncMap with all functions exposed by the service.
func (s *Service) Funcs() texttmpl.FuncMap {
	return texttmpl.FuncMap{
		// GitHub
		"recentContributions": s.RecentContributions,
		"recentPullRequests":  s.RecentPullRequests,
		"recentRepos":         s.RecentRepos,
		"recentForks":         s.RecentForks,
		"recentReleases":      s.RecentReleases,
		"followers":           s.Followers,
		"recentStars":         s.RecentStars,
		"gists":               s.Gists,
		"recentIssues":        s.RecentIssues,
		"sponsors":            s.Sponsors,
		"repo":                s.Repo,
		// RSS
		"rss": s.LatestRssFeeds,
		// GoodReads
		"goodReadsReviews":          s.GoodReadsReviews,
		"goodReadsCurrentlyReading": s.GoodReadsCurrentlyReading,
		// Literal.club
		"literalClubCurrentlyReading": s.LiteralCurrentlyReading,
		// Utils
		"humanize": s.Humanize,
		"reverse":  s.Reverse,
		"now":      time.Now,
		"contains": strings.Contains,
		"toLower":  strings.ToLower,
	}
}
