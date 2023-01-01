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

package main

import (
	"errors"

	"github.com/sbuckfelder/github-monitoring-tool/proxy"
	"github.com/urfave/cli"
)

const (
	ORG_NAME     string = "org"
	REPOALL_NAME string = "repoall"
	REPO_NAME    string = "repo"
	SINCE_NAME   string = "since"
	DATE_NAME    string = "date"
	HOURS_NAME   string = "hours"
)

var eventsCommand = cli.Command{
	Name:  "events",
	Usage: "List events for a github org/repo provide one of [hours,date,since]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     ORG_NAME,
			Usage:    "github org repo belongs to",
			Required: true,
		},
		cli.BoolFlag{
			Name:     REPOALL_NAME,
			Usage:    "get events for all repos in an organization, cannot be used with repo flag",
			Required: false,
		},
		cli.StringFlag{
			Name:     REPO_NAME,
			Usage:    "github repo to list events for, cannot be used with repoall flag",
			Required: false,
		},
		cli.StringFlag{
			Name:     SINCE_NAME,
			Usage:    "timestamp in RFC3339 format e.g. 2006-01-02T15:04:05Z",
			Required: false,
		},
		cli.IntFlag{
			Name:     HOURS_NAME,
			Usage:    "number of hours to look back",
			Required: false,
		},
		cli.StringFlag{
			Name:     DATE_NAME,
			Usage:    "date for events format YYYY-MM-DD",
			Required: false,
		},
	},
	Action: func(ctx *cli.Context) error {
		orgInput := ctx.String(ORG_NAME)
		repoAllInput := ctx.Bool(REPOALL_NAME)
		repoInput := ctx.String(REPO_NAME)
		sinceInput := ctx.String(SINCE_NAME)
		dateInput := ctx.String(DATE_NAME)
		hoursInput := ctx.Int(HOURS_NAME)

		if repoAllInput && repoInput != "" {
			return errors.New("Both 'repo' and 'repoall' flag cannot be set")
		}

		if !repoAllInput && repoInput == "" {
			return errors.New("Either 'repo' or 'repoall' needs to be set")
		}

		if sinceInput == "" && dateInput == "" && hoursInput == 0 {
			return errors.New("No time flag [since,date,hours] is set")
		}

		if sinceInput != "" && dateInput != "" {
			return errors.New("Cannot have both since and date set")
		}

		if sinceInput != "" && hoursInput != 0 {
			return errors.New("Cannot have both since and hours set")
		}

		if dateInput != "" && hoursInput != 0 {
			return errors.New("Cannot have both date and hours set")
		}

		ghProxy, err := proxy.NewProxy()
		if err != nil {
			panic("Failed to create GitHub client")
		}

		var repos = []string{}

		if repoAllInput {
			repos = ghProxy.GetReposForOrg(orgInput)
		} else {
			repos = append(repos, repoInput)
		}

		if dateInput != "" {
			ghProxy.GetEventsForDate(orgInput, repos, dateInput)
		}

		if sinceInput != "" {
			ghProxy.GetEventsSinceRFC3339(orgInput, repos, sinceInput)
		}

		if hoursInput != 0 {
			ghProxy.GetEventsForHours(orgInput, repos, hoursInput)
		}

		return nil
	},
}
