Feature: detached sync the current feature branch using the "rebase" sync-feature strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And Git Town setting "sync-feature-strategy" is "rebase"
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    When I run "git-town sync --detached"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git rebase main                                 |
      |         | git push --force-with-lease --force-if-includes |
      |         | git rebase origin/feature                       |
      |         | git push --force-with-lease --force-if-includes |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local         | local main commit     |
      |         | origin        | origin main commit    |
      | feature | local, origin | origin feature commit |
      |         |               | local main commit     |
      |         |               | local feature commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                                                           |
      | feature | git reset --hard {{ sha-before-run 'local feature commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin feature commit' }}:feature |
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist