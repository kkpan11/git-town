package configfile_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/config/configfile"
	"github.com/git-town/git-town/v16/internal/config/gitconfig"
	"github.com/git-town/git-town/v16/internal/git"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestSave(t *testing.T) {
	t.Parallel()

	t.Run("RenderPerennialBranches", func(t *testing.T) {
		t.Parallel()
		t.Run("no perennial branches", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.LocalBranchNames{}
			have := configfile.RenderPerennialBranches(give)
			want := "[]"
			must.EqOp(t, want, have)
		})
		t.Run("one perennial branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.NewLocalBranchNames("one")
			have := configfile.RenderPerennialBranches(give)
			want := `["one"]`
			must.EqOp(t, want, have)
		})
		t.Run("multiple perennial branches", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.NewLocalBranchNames("one", "two")
			have := configfile.RenderPerennialBranches(give)
			want := `["one", "two"]`
			must.EqOp(t, want, have)
		})
	})

	t.Run("RenderTOML", func(t *testing.T) {
		t.Parallel()
		give := config.UnvalidatedConfig{
			UnvalidatedConfig: configdomain.UnvalidatedConfigData{
				MainBranch: Some(gitdomain.NewLocalBranchName("main")),
			},
			NormalConfig: config.NormalConfig{
				NormalConfigData: configdomain.NormalConfigData{
					DefaultBranchType:        configdomain.BranchTypeFeatureBranch,
					FeatureRegex:             None[configdomain.FeatureRegex](),
					HostingOriginHostname:    None[configdomain.HostingOriginHostname](),
					HostingPlatform:          None[configdomain.HostingPlatform](),
					Lineage:                  configdomain.NewLineage(),
					NewBranchType:            configdomain.BranchTypePrototypeBranch,
					ObservedBranches:         gitdomain.LocalBranchNames{},
					Offline:                  false,
					ParkedBranches:           gitdomain.LocalBranchNames{},
					PerennialBranches:        gitdomain.NewLocalBranchNames("one", "two"),
					PerennialRegex:           None[configdomain.PerennialRegex](),
					PushHook:                 true,
					PushNewBranches:          false,
					ShipStrategy:             configdomain.ShipStrategySquashMerge,
					ShipDeleteTrackingBranch: true,
					SyncFeatureStrategy:      configdomain.SyncFeatureStrategyMerge,
					SyncPerennialStrategy:    configdomain.SyncPerennialStrategyRebase,
					SyncTags:                 true,
					SyncUpstream:             true,
				},
			},
		}
		have := configfile.RenderTOML(&give)
		want := `
# Git Town configuration file
#
# Run "git town config setup" to add additional entries
# to this file after updating Git Town.

[branches]

# The main branch is the branch from which you cut new feature branches,
# and into which you ship feature branches when they are done.
# This branch is often called "main", "master", or "development".
main = "main"

# Perennial branches are long-lived branches.
# They are never shipped and have no ancestors.
# Typically, perennial branches have names like
# "development", "staging", "qa", "production", etc.
#
# See also the "perennial-regex" setting.
perennials = ["one", "two"]

# All branches whose name matches this regular expression
# are also considered perennial branches.
#
# If you are not sure, leave this empty.
perennial-regex = ""

[create]

# The "new-branch-type" setting determines which branch type Git Town
# creates when you run "git town hack", "append", or "prepend".
#
# More info at https://www.git-town.com/preferences/new-branch-type.
new-branch-type = "prototype"

# Should Git Town push the new branches it creates
# immediately to origin even if they are empty?
#
# When enabled, you can run "git push" right away
# but creating new branches is slower and
# it triggers an unnecessary CI run on the empty branch.
#
# When disabled, many Git Town commands execute faster
# and Git Town will create the missing tracking branch
# on the first run of "git town sync".
push-new-branches = false

[hosting]

# Knowing the type of code hosting platform allows Git Town
# to open browser URLs and talk to the code hosting API.
# Most people can leave this on "auto-detect".
# Only change this if your code hosting server uses as custom URL.
# platform = ""

# When using SSH identities, define the hostname
# of your source code repository. Only change this
# if the auto-detection does not work for you.
# origin-hostname = ""

[ship]

# Should "git town ship" delete the tracking branch?
# You want to disable this if your code hosting platform
# (GitHub, GitLab, etc) deletes head branches when
# merging pull requests through its UI.
delete-tracking-branch = true

# Which method should Git Town use to ship feature branches?
#
# Options:
#
# - api: merge the proposal on your code hosting platform via the code hosting API
# - fast-forward: in your local repo, fast-forward the parent branch to point to the commits on the feature branch
# - squash-merge: in your local repo, squash-merge the feature branch into its parent branch
#
# All options update proposals of child branches and remove the shipped branch locally and remotely.
strategy = "squash-merge"

[sync]

# How should Git Town synchronize feature branches?
# Feature branches are short-lived branches cut from
# the main branch and shipped back into the main branch.
# Typically you develop features and bug fixes on them,
# hence their name.
feature-strategy = "merge"

# How should Git Town synchronize perennial branches?
# Perennial branches have no parent branch.
# The only updates they receive are additional commits
# made to their tracking branch somewhere else.
perennial-strategy = "rebase"
prototype-strategy = ""

# The "push-hook" setting determines whether Git Town
# permits or prevents Git hooks while pushing branches.
# Hooks are enabled by default. If your Git hooks are slow,
# you can disable them to speed up branch syncing.
#
# When disabled, Git Town pushes using the "--no-verify" switch.
# More info at https://www.git-town.com/preferences/push-hook.
push-hook = true

# Should "git town sync" sync tags with origin?
tags = true

# Should "git town sync" also fetch updates from the upstream remote?
#
# If an "upstream" remote exists, and this setting is enabled,
# "git town sync" will also update the local main branch
# with commits from the main branch at the upstream remote.
#
# This is useful if the repository you work on is a fork,
# and you want to keep it in sync with the repo it was forked from.
upstream = true
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("Save", func(t *testing.T) {
		t.Parallel()
		var gitAccess gitconfig.Access
		config := config.DefaultUnvalidatedConfig(gitAccess, git.EmptyVersion())
		config.UnvalidatedConfig.MainBranch = Some(gitdomain.NewLocalBranchName("main"))
		err := configfile.Save(&config)
		defer os.Remove(configfile.FileName)
		must.NoError(t, err)
		bytes, err := os.ReadFile(configfile.FileName)
		must.NoError(t, err)
		have := string(bytes)
		want := `
# Git Town configuration file
#
# Run "git town config setup" to add additional entries
# to this file after updating Git Town.

[branches]

# The main branch is the branch from which you cut new feature branches,
# and into which you ship feature branches when they are done.
# This branch is often called "main", "master", or "development".
main = "main"

# Perennial branches are long-lived branches.
# They are never shipped and have no ancestors.
# Typically, perennial branches have names like
# "development", "staging", "qa", "production", etc.
#
# See also the "perennial-regex" setting.
perennials = []

# All branches whose name matches this regular expression
# are also considered perennial branches.
#
# If you are not sure, leave this empty.
perennial-regex = ""

[create]

# The "new-branch-type" setting determines which branch type Git Town
# creates when you run "git town hack", "append", or "prepend".
#
# More info at https://www.git-town.com/preferences/new-branch-type.
new-branch-type = "feature"

# Should Git Town push the new branches it creates
# immediately to origin even if they are empty?
#
# When enabled, you can run "git push" right away
# but creating new branches is slower and
# it triggers an unnecessary CI run on the empty branch.
#
# When disabled, many Git Town commands execute faster
# and Git Town will create the missing tracking branch
# on the first run of "git town sync".
push-new-branches = false

[hosting]

# Knowing the type of code hosting platform allows Git Town
# to open browser URLs and talk to the code hosting API.
# Most people can leave this on "auto-detect".
# Only change this if your code hosting server uses as custom URL.
# platform = ""

# When using SSH identities, define the hostname
# of your source code repository. Only change this
# if the auto-detection does not work for you.
# origin-hostname = ""

[ship]

# Should "git town ship" delete the tracking branch?
# You want to disable this if your code hosting platform
# (GitHub, GitLab, etc) deletes head branches when
# merging pull requests through its UI.
delete-tracking-branch = true

# Which method should Git Town use to ship feature branches?
#
# Options:
#
# - api: merge the proposal on your code hosting platform via the code hosting API
# - fast-forward: in your local repo, fast-forward the parent branch to point to the commits on the feature branch
# - squash-merge: in your local repo, squash-merge the feature branch into its parent branch
#
# All options update proposals of child branches and remove the shipped branch locally and remotely.
strategy = "api"

[sync]

# How should Git Town synchronize feature branches?
# Feature branches are short-lived branches cut from
# the main branch and shipped back into the main branch.
# Typically you develop features and bug fixes on them,
# hence their name.
feature-strategy = "merge"

# How should Git Town synchronize perennial branches?
# Perennial branches have no parent branch.
# The only updates they receive are additional commits
# made to their tracking branch somewhere else.
perennial-strategy = "rebase"
prototype-strategy = "rebase"

# The "push-hook" setting determines whether Git Town
# permits or prevents Git hooks while pushing branches.
# Hooks are enabled by default. If your Git hooks are slow,
# you can disable them to speed up branch syncing.
#
# When disabled, Git Town pushes using the "--no-verify" switch.
# More info at https://www.git-town.com/preferences/push-hook.
push-hook = true

# Should "git town sync" sync tags with origin?
tags = true

# Should "git town sync" also fetch updates from the upstream remote?
#
# If an "upstream" remote exists, and this setting is enabled,
# "git town sync" will also update the local main branch
# with commits from the main branch at the upstream remote.
#
# This is useful if the repository you work on is a fork,
# and you want to keep it in sync with the repo it was forked from.
upstream = true
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("TOMLComment", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			have := configfile.TOMLComment("")
			want := ""
			must.Eq(t, want, have)
		})
		t.Run("single line", func(t *testing.T) {
			t.Parallel()
			have := configfile.TOMLComment("line 1")
			want := "# line 1"
			must.Eq(t, want, have)
		})
		t.Run("multiple lines", func(t *testing.T) {
			t.Parallel()
			have := configfile.TOMLComment("line 1\nline 2\nline 3")
			want := "# line 1\n# line 2\n# line 3"
			must.Eq(t, want, have)
		})
		t.Run("multiple lines with terminating newline", func(t *testing.T) {
			t.Parallel()
			have := configfile.TOMLComment("line 1\nline 2\nline 3\n")
			want := "# line 1\n# line 2\n# line 3\n#"
			must.Eq(t, want, have)
		})
	})
}
