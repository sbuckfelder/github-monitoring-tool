/*
   Copyright awslabs Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package proxy

import (
	"context"
	"fmt"

	"github.com/google/go-github/v48/github"
)

var PRS_PER_PAGE int = 100

func (p *GithubProxy) GetPullRequests(org, repo string) {
	ctx := context.Background()
	pullRequests, err := p.getAllOpenPullRequests(ctx, org, repo)
	if err != nil {
		fmt.Printf("Failed to get pull requests: %v\n", err)
		return
	}
	var output = []string{}
	titleLine := "Title,URL,Created,Updated,PR Author,LastCommentDate,CommentAuthor"
	output = append(output, titleLine)
	for _, PR := range pullRequests {
		comment := p.getLastComment(ctx, org, repo, *PR.Number)
		commentCsv := ""
		if comment != nil {
			commentCsv = fmt.Sprintf("%v,%s",
				comment.CreatedAt,
				*comment.User.Login)
		}
		csvLine := fmt.Sprintf("%s,%s,%v,%v,%s,%s",
			*PR.Title,
			*PR.HTMLURL,
			PR.CreatedAt,
			PR.UpdatedAt,
			*PR.User.Login,
			commentCsv)
		output = append(output, csvLine)
	}
	for _, line := range output {
		fmt.Println(line)
	}
}

func (p *GithubProxy) getAllOpenPullRequests(ctx context.Context, org, repo string) ([]*github.PullRequest, error) {
	openPRs := []*github.PullRequest{}
	morePRs := true
	pageNum := 1
	for morePRs {
		prOpts := &github.PullRequestListOptions{
			State:     "Open",
			Sort:      "Created",
			Direction: "asc",
			ListOptions: github.ListOptions{
				Page:    pageNum,
				PerPage: PRS_PER_PAGE,
			},
		}
		pullRequests, _, err := p.client.PullRequests.List(ctx, org, repo, prOpts)
		if err != nil {
			fmt.Printf("Error Listing Pull Requests: %v\n", err)

		}
		if len(pullRequests) == 0 {
			morePRs = false
		} else {
			openPRs = append(openPRs, pullRequests...)
			pageNum++
		}
	}
	return openPRs, nil
}

func (p *GithubProxy) getLastComment(ctx context.Context, org, repo string, prNum int) *github.PullRequestComment {
	commentOpts := &github.PullRequestListCommentsOptions{
		Sort:      "Created",
		Direction: "desc",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 1,
		},
	}
	comments, _, _ := p.client.PullRequests.ListComments(ctx, org, repo, prNum, commentOpts)
	if len(comments) == 0 {
		return nil
	}
	return comments[0]
}
