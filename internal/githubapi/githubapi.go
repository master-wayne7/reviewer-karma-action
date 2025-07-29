package githubapi

import (
	"context"

	"github.com/google/go-github/v62/github"
)

// FetchAllPullRequests fetches all pull requests from a repository
func FetchAllPullRequests(ctx context.Context, client *github.Client, owner, repo string) ([]*github.PullRequest, error) {
	var allPRs []*github.PullRequest
	opts := &github.PullRequestListOptions{
		State: "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		prs, resp, err := client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}
		allPRs = append(allPRs, prs...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allPRs, nil
}

// FetchPullRequestReviews fetches all reviews for a specific pull request
func FetchPullRequestReviews(ctx context.Context, client *github.Client, owner, repo string, prNumber int) ([]*github.PullRequestReview, error) {
	var allReviews []*github.PullRequestReview
	opts := &github.ListOptions{
		PerPage: 100,
	}

	for {
		reviews, resp, err := client.PullRequests.ListReviews(ctx, owner, repo, prNumber, opts)
		if err != nil {
			return nil, err
		}
		allReviews = append(allReviews, reviews...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allReviews, nil
}

// FetchPullRequestComments fetches all comments for a specific pull request
func FetchPullRequestComments(ctx context.Context, client *github.Client, owner, repo string, prNumber int) ([]*github.PullRequestComment, error) {
	var allComments []*github.PullRequestComment
	opts := &github.PullRequestListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		comments, resp, err := client.PullRequests.ListComments(ctx, owner, repo, prNumber, opts)
		if err != nil {
			return nil, err
		}
		allComments = append(allComments, comments...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allComments, nil
}
