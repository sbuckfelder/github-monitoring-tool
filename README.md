## CLI Tool for Monitoring Projects Hosted on Github 
Tool is designed to provide light weight functionality to gather event, issue, and PR data for projects on GitHub.

## Installation
Currently the tool needs to be build from source using the `make` command.

## Usage

### Events
The event command works with a concept of a lookback.  There are three lookback modes available, they can not be used together:
- since: this returns all events from a timestamp till present.  The format is RFC3339 and only supports UTC.  e.g. 2006-01-02T15:04:05Z
- hours: this returns all events from a number of hours till present.  It accepts integers
- date: this returns all events for a date (UTC).  The format is YYYY-MM-DD.

The command takes in either the `--repo` flag or the `--repoall` flag, but they cannot be used together.  The `--repo` flag is used when targeting a single repo, while the `--repoall` flag will return events for all repositories under an organization.

**Example Usage**
Returns all events in the https://github.com/containerd/containerd repository that occurred in the last 2 hours
```
./ghmt events --org containerd --repo containerd --hours 2
```

Returns all events in the repositories that belong to https://github.com/containerd org that occurred since
```
./ghmt events --org containerd --repoall --since 2023-01-02T17:00:00Z
```

### PRs
The `pr` command writes a csv report to the command line and is still under development

**Example Usage**
Returns all open PRs with time stamps of when they were created, updated, and last commented on.
```
./ghmt events --org containerd --repo containerd
```






