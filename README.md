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
- token
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

#### Turn up the volume
As of now, I'm relying on `say` for this to actually notify the user, so you'll
need to be able to hear it.

#### Set up the Slack App
- Create a slack app and set up an incoming webhook so you can post messages to a slack channel

#### Start the server
In a terminal session, run `make start`

## How to get a Personal Access Token?
- Log in to your repository
- Click your icon in the top right corner
- Click "Developer Settings" on the bottom of the left panel
- Click "Personal access tokens" on the bottom of the left panel
- Click "Generate new token"
    - Name it whatever you want
    - Give the token the following scopes
        - repo (top level)
        - notifications
        - user (top level)
        - read:discussion
- Copy the token value and paste it into the app config "token" value