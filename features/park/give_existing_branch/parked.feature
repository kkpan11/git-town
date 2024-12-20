Feature: park another already parked branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the current branch is "main"
    When I run "git-town park parked"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      branch "parked" is already parked
      """
    And the parked branches are still "parked"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the parked branches are still "parked"
