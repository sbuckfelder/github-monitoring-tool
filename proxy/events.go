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
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v48/github"
)

var EVENTS_PER_PAGE int = 100
var REPORT_SEPERATOR string = strings.Repeat("*", 20)

func (p *GithubProxy) GetEventsForHours(org string, repos []string, hours int) {
	currentTime := time.Now()
	adjTime := currentTime.Add(time.Hour * time.Duration(-1*hours))
	ctx := context.Background()
	for _, repo := range repos {
		events, _ := p.getEventsSince(ctx, org, repo, adjTime)
		fmt.Printf("https://github.com/%s/%s Events Since %s\n", org, repo, adjTime.Format(time.RFC3339))
		GenerateEventReport(events, repo)
	}
	fmt.Printf("_Based on Events from %s to %s_\n",
		adjTime.Format(time.RFC3339),
		currentTime.Format(time.RFC3339))
}

func (p *GithubProxy) GetEventsSinceRFC3339(org string, repos []string, sinceString string) {
	currentTime := time.Now()
	since, err := time.Parse(time.RFC3339, sinceString)
	if err != nil {
		panicMsg := fmt.Sprintf("Since value:'%s' not in RFC3339 e.g. 2006-01-02T15:04:05Z", sinceString)
		panic(panicMsg)
	}
	ctx := context.Background()
	for _, repo := range repos {
		events, _ := p.getEventsSince(ctx, org, repo, since)
		fmt.Printf("https://github.com/%s/%s Events Since %s\n", org, repo, sinceString)
		GenerateEventReport(events, repo)
	}
	fmt.Printf("_Based on Events from %s to %s_\n",
		since.Format(time.RFC3339),
		currentTime.Format(time.RFC3339))
}

func (p *GithubProxy) GetEventsForDate(org string, repos []string, dateString string) {
	const dateLayout = "2006-01-02"
	targetDate, err := time.Parse(dateLayout, dateString)
	if err != nil {
		panicMsg := fmt.Sprintf("Since value:'%s' not in YYYY-MM-DD e.g. 2006-01-02", dateString)
		panic(panicMsg)
	}
	ctx := context.Background()
	for _, repo := range repos {
		events, _ := p.getEventsSince(ctx, org, repo, targetDate)
		events = filterEvents(events, eventFilterDate(targetDate))
		fmt.Printf("https://github.com/%s/%s Events for %s\n", org, repo, dateString)
		GenerateEventReport(events, repo)
	}
	fmt.Printf("_Based on Events for %s\n", dateString)
}

func (p *GithubProxy) getEventsSince(
	ctx context.Context,
	org, repo string,
	since time.Time) ([]*github.Event, error) {
	pageNumber := 1
	found := false
	var events []*github.Event
	for found == false {
		listOpts := &github.ListOptions{
			PerPage: EVENTS_PER_PAGE,
			Page:    pageNumber}

		newEvents, _, err := p.client.Activity.ListRepositoryEvents(ctx, org, repo, listOpts)
		if err != nil {
			return nil, err
		}
		if len(newEvents) == 0 {
			found = true
		} else {
			lastEvent := newEvents[len(newEvents)-1]
			if lastEvent.CreatedAt.Before(since) {
				found = true
			}
		}
		events = append(events, newEvents...)
		pageNumber++
	}
	events = filterEvents(events, eventFilterSince(since))
	return events, nil
}

func filterEvents(events []*github.Event, filter func(*github.Event) bool) []*github.Event {
	var filteredEvents []*github.Event
	for _, event := range events {
		if filter(event) {
			filteredEvents = append(filteredEvents, event)
		}
	}
	return filteredEvents
}

func eventFilterSince(since time.Time) func(*github.Event) bool {
	return func(event *github.Event) bool {
		return event.CreatedAt.After(since)
	}
}

func eventFilterDate(date time.Time) func(*github.Event) bool {
	return func(event *github.Event) bool {
		return event.CreatedAt.Year() == date.Year() &&
			event.CreatedAt.Month() == date.Month() &&
			event.CreatedAt.Day() == date.Day()

	}
}

func eventFilterType() func(*github.Event) bool {
	return func(event *github.Event) bool {
		switch *event.Type {
		case "IssuesEvent", "IssuesCommentEvent", "PullRequestEvent", "PullRequestReviewEvent", "PullRequestReviewCommentEvent":
			return true
		default:
			return false
		}
	}
}

func printPullRequestReport(events []*github.Event, repo string) {
	var newPRs = []*github.Event{}
	var mergedPRs = []*github.Event{}
	for _, event := range events {
		payload, _ := event.ParsePayload()
		prEvent := payload.(*github.PullRequestEvent)
		if *prEvent.Action == "opened" {
			newPRs = append(newPRs, event)
		}
		if *prEvent.Action == "closed" && *prEvent.PullRequest.Merged {
			mergedPRs = append(mergedPRs, event)
		}
	}
	if len(mergedPRs) != 0 {
		fmt.Printf("===========\n")
		fmt.Printf("PR'S MERGED\n")
		fmt.Printf("===========\n")
		var mergedPRText = sort.StringSlice{}
		for _, prEvent := range mergedPRs {
			mergedPRText = append(mergedPRText, printPullRequestEvent(prEvent, repo))
		}
		mergedPRText.Sort()
		for _, txt := range mergedPRText {
			fmt.Print(txt)
		}
	}
	if len(newPRs) != 0 {
		fmt.Printf("=================\n")
		fmt.Printf("NEW PULL REQUESTS\n")
		fmt.Printf("=================\n")
		var newPRText = sort.StringSlice{}
		for _, prEvent := range newPRs {
			newPRText = append(newPRText, printPullRequestEvent(prEvent, repo))
		}
		newPRText.Sort()
		for _, txt := range newPRText {
			fmt.Print(txt)
		}
	}
}

func printPullRequestEvent(event *github.Event, repo string) string {
	payload, _ := event.ParsePayload()
	prEvent := payload.(*github.PullRequestEvent)

	return fmt.Sprintf("- **%s** PR#%d %s: [%s](%s)\n",
		repo,
		*prEvent.Number,
		*prEvent.PullRequest.User.Login,
		*prEvent.PullRequest.Title,
		*prEvent.PullRequest.HTMLURL)
}

func getPullRequestURLFromReview(event *github.Event) (string, string) {
	if *event.Type == "PullRequestReviewEvent" {
		payload, _ := event.ParsePayload()
		rev := payload.(*github.PullRequestReviewEvent)
		return *rev.PullRequest.HTMLURL, *rev.PullRequest.Title
	}
	if *event.Type == "PullRequestReviewCommentEvent" {
		payload, _ := event.ParsePayload()
		rev := payload.(*github.PullRequestReviewCommentEvent)
		return *rev.PullRequest.HTMLURL, *rev.PullRequest.Title
	}
	return "", ""
}

func printPullRequestReviewReport(events []*github.Event, repo string) {
	reviewMap := make(map[string]int)
	titleMap := make(map[string]string)
	for _, event := range events {
		url, title := getPullRequestURLFromReview(event)
		reviewMap[url] = reviewMap[url] + 1
		titleMap[url] = title
	}
	fmt.Printf("==========================\n")
	fmt.Printf("PR REVIEW/COMMENT ACTIVITY\n")
	fmt.Printf("==========================\n")
	var revText = sort.StringSlice{}
	for url, val := range reviewMap {
		revText = append(revText, fmt.Sprintf("Actions:%d - **%s** [%s](%s)\n", val, repo, titleMap[url], url))
	}
	revText.Sort()
	for _, txt := range revText {
		fmt.Print(txt)
	}
}

func printIssueEventReport(events []*github.Event, repo string) {
	var newIssues = []*github.Event{}
	var closedIssues = []*github.Event{}
	for _, event := range events {
		payload, _ := event.ParsePayload()
		issuesEvent := payload.(*github.IssuesEvent)
		if *issuesEvent.Action == "opened" {
			newIssues = append(newIssues, event)
		}
		if *issuesEvent.Action == "closed" {
			closedIssues = append(closedIssues, event)
		}
	}
	if len(newIssues) != 0 {
		fmt.Printf("==========\n")
		fmt.Printf("NEW ISSUES\n")
		fmt.Printf("==========\n")
		var newIssueText = sort.StringSlice{}
		for _, issueEvent := range newIssues {
			newIssueText = append(newIssueText, printIssueEvent(issueEvent, repo))
		}
		newIssueText.Sort()
		for _, txt := range newIssueText {
			fmt.Print(txt)
		}
	}
	if len(closedIssues) != 0 {
		fmt.Printf("=============\n")
		fmt.Printf("CLOSED ISSUES\n")
		fmt.Printf("=============\n")
		var closedIssueText = sort.StringSlice{}
		for _, issueEvent := range closedIssues {
			closedIssueText = append(closedIssueText, printIssueEvent(issueEvent, repo))
		}
		closedIssueText.Sort()
		for _, txt := range closedIssueText {
			fmt.Print(txt)
		}
	}
}

func printIssueEvent(event *github.Event, repo string) string {
	payload, _ := event.ParsePayload()
	issuesEvent := payload.(*github.IssuesEvent)
	return fmt.Sprintf("- **%s** ISSUE#%d %s: [%s](%s)\n",
		repo,
		*issuesEvent.Issue.Number,
		*issuesEvent.Issue.User.Login,
		*issuesEvent.Issue.Title,
		*issuesEvent.Issue.HTMLURL)
}

func printIssueCommentEventReport(events []*github.Event, repo string) {
	commentMap := make(map[string]int)
	titleMap := make(map[string]string)
	for _, event := range events {
		payload, _ := event.ParsePayload()
		commentEvent := payload.(*github.IssueCommentEvent)
		url := *commentEvent.Issue.HTMLURL
		commentMap[url] = commentMap[url] + 1
		titleMap[url] = *commentEvent.Issue.Title
	}
	fmt.Printf("======================\n")
	fmt.Printf("ISSUE COMMENT ACTIVITY\n")
	fmt.Printf("======================\n")
	var commentText = sort.StringSlice{}
	for url, val := range commentMap {
		commentText = append(commentText, fmt.Sprintf("Comments:%d - **%s** [%s](%s)\n", val, repo, titleMap[url], url))
	}
	commentText.Sort()
	for _, txt := range commentText {
		fmt.Print(txt)
	}
}

func GenerateEventReport(events []*github.Event, repo string) {
	defer fmt.Printf("%s\n", REPORT_SEPERATOR)
	eventMap := make(map[string][]*github.Event)
	for _, event := range events {
		eventList, ok := eventMap[*event.Type]
		if !ok {
			eventList = []*github.Event{}
		}
		eventList = append(eventList, event)
		eventMap[*event.Type] = eventList
	}
	if len(eventMap) == 0 {
		fmt.Printf("No Events\n")
		return
	}
	prEvents, _ := eventMap["PullRequestEvent"]
	if len(prEvents) != 0 {
		printPullRequestReport(prEvents, repo)
	}
	revEvents, _ := eventMap["PullRequestReviewEvent"]
	revCommentEvents, _ := eventMap["PullRequestReviewCommentEvent"]
	revEvents = append(revEvents, revCommentEvents...)
	if len(revEvents) != 0 {
		printPullRequestReviewReport(revEvents, repo)
	}
	issueEvents, _ := eventMap["IssuesEvent"]
	if len(issueEvents) != 0 {
		printIssueEventReport(issueEvents, repo)
	}
	issueCommentEvents, _ := eventMap["IssueCommentEvent"]
	if len(issueCommentEvents) != 0 {
		printIssueCommentEventReport(issueCommentEvents, repo)
	}
	fmt.Printf("============\n")
	fmt.Printf("EVENT REPORT\n")
	fmt.Printf("============\n")
	for key, val := range eventMap {
		fmt.Printf("%d %s\n", len(val), key)
	}
}
