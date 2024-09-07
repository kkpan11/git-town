Feature: ship a coworker's feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE            | AUTHOR                            |
      | feature | local, origin | developer commit 1 | developer <developer@example.com> |
      |         |               | developer commit 2 | developer <developer@example.com> |
      |         |               | coworker commit    | coworker <coworker@example.com>   |
    And the current branch is "feature"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship"

  Scenario: result
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE            | AUTHOR                            |
      | main   | local, origin | developer commit 1 | developer <developer@example.com> |
      |        |               | developer commit 2 | developer <developer@example.com> |
      |        |               | coworker commit    | coworker <coworker@example.com>   |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                        |
      | main   | git branch feature {{ sha 'coworker commit' }} |
      |        | git push -u origin feature                     |
      |        | git checkout feature                           |
    And the current branch is now "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE            |
      | main    | local, origin | developer commit 1 |
      |         |               | developer commit 2 |
      |         |               | coworker commit    |
      | feature | local, origin | developer commit 1 |
      |         |               | developer commit 2 |
      |         |               | coworker commit    |
    And the initial branches and lineage exist