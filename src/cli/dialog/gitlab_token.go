package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialog/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const (
	gitLabTokenTitle = `GitLab API token`
	gitLabTokenHelp  = `
If you have an API token for GitLab,
and want to ship branches from the CLI,
please enter it now.

It's okay to leave this empty.

`
)

// GitLabToken lets the user enter the GitHub API token.
func GitLabToken(oldValue configdomain.GitLabToken, inputs components.TestInput) (configdomain.GitLabToken, bool, error) {
	token, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          gitLabTokenHelp,
		Prompt:        "Your GitLab API token: ",
		TestInput:     inputs,
		Title:         gitLabTokenTitle,
	})
	fmt.Printf("GitLab token: %s\n", components.FormattedSecret(token, aborted))
	return configdomain.GitLabToken(token), aborted, err
}
