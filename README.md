# repo-watcher
A Ruby application that provides desktop notifications from a GHe repo

## Notes
- This will only work with GitHub Enterprise repos.  GitHub hosted repos use a different API scheme that is not supported here.
- Disclaimer: I have only tested this with a branch "owned" by an organization, so there might be some bugs around repos not set up in this way.

## How to use:
#### Clone the repo
`git clone https://github.com/mike-webster/repo-watcher.git`

#### App Config
Fill out the following values appropriately
- ~~token~~
    - GitHub Personal Access Token
- org_name
    - The organization that owns the repository
- repo_host
    - The url for where the repository is hosted
- name
    - The name of the user running the application
- refresh_seconds
    - How often you want the application to check for activity
- repo_to_watch
    - The repo you want to monitor
- username
    - Your github username, this is used to silence your own events
- slack_webhook
    - The webhook you set up for your Slack app

#### Write your deployment manifest
- Configure and deploy your app anywhere that has access to your GitHub Enterprise repository.

#### Set up the Slack App
- Create a slack app and set up an incoming webhook so you can post messages to a slack channel

#### Start the server
In a terminal session, run `make start`

## How to configure your webhooks?
- TODO