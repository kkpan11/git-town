Feature: display all executed Git commands

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
    And the current branch is "existing"

  Scenario: result
    When I run "git-town append new --verbose"
    Then it runs the commands
      | BRANCH   | TYPE     | COMMAND                                                         |
      |          | backend  | git version                                                     |
      |          | backend  | git rev-parse --show-toplevel                                   |
      |          | backend  | git config -lz --includes --global                              |
      |          | backend  | git config -lz --includes --local                               |
      |          | backend  | git status --long --ignore-submodules                           |
      |          | backend  | git remote                                                      |
      |          | backend  | git rev-parse --abbrev-ref HEAD                                 |
      | existing | frontend | git fetch --prune --tags                                        |
      |          | backend  | git stash list                                                  |
      |          | backend  | git branch -vva --sort=refname                                  |
      |          | backend  | git rev-parse --verify --abbrev-ref @{-1}                       |
      |          | backend  | git remote get-url origin                                       |
      |          | backend  | git log main..existing --format=%h                              |
      |          | backend  | git log --format=%B -n 1 {{ sha-before-run 'existing commit' }} |
      | existing | frontend | git checkout main                                               |
      | main     | frontend | git rebase origin/main                                          |
      |          | backend  | git rev-list --left-right main...origin/main                    |
      | main     | frontend | git checkout existing                                           |
      | existing | frontend | git merge --no-edit --ff origin/existing                        |
      |          | frontend | git merge --no-edit --ff main                                   |
      |          | backend  | git rev-list --left-right existing...origin/existing            |
      |          | backend  | git show-ref --verify --quiet refs/heads/existing               |
      | existing | frontend | git checkout -b new                                             |
      |          | backend  | git show-ref --verify --quiet refs/heads/existing               |
      |          | backend  | git config git-town-branch.new.parent existing                  |
      |          | backend  | git show-ref --verify --quiet refs/heads/existing               |
      |          | backend  | git branch -vva --sort=refname                                  |
      |          | backend  | git config -lz --includes --global                              |
      |          | backend  | git config -lz --includes --local                               |
      |          | backend  | git stash list                                                  |
    And it prints:
      """
      Ran 30 shell commands.
      """
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town append new"
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH   | TYPE     | COMMAND                                       |
      |          | backend  | git version                                   |
      |          | backend  | git rev-parse --show-toplevel                 |
      |          | backend  | git config -lz --includes --global            |
      |          | backend  | git config -lz --includes --local             |
      |          | backend  | git status --long --ignore-submodules         |
      |          | backend  | git stash list                                |
      |          | backend  | git branch -vva --sort=refname                |
      |          | backend  | git remote get-url origin                     |
      |          | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |          | backend  | git remote get-url origin                     |
      | new      | frontend | git checkout existing                         |
      | existing | frontend | git branch -D new                             |
      |          | backend  | git config --unset git-town-branch.new.parent |
    And it prints:
      """
      Ran 13 shell commands.
      """
    And the current branch is still "existing"
    And the initial commits exist now
    And the initial branches and lineage exist now
