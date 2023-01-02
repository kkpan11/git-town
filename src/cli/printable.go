package cli

import (
	"sort"
	"strconv"
	"strings"
)

// BranchAncestryConfig defines the configuration values needed by the `cli` package.
type BranchAncestryConfig interface {
	BranchAncestryRoots() []string
	ChildBranches(string) []string
}

// PrintableBranchAncestry provides the branch ancestry in CLI printable format.
func PrintableBranchAncestry(config BranchAncestryConfig) string {
	roots := config.BranchAncestryRoots()
	trees := make([]string, len(roots))
	for r, root := range roots {
		trees[r] = PrintableBranchTree(root, config)
	}
	return strings.Join(trees, "\n\n")
}

// PrintableBranchTree returns a user printable branch tree.
func PrintableBranchTree(branchName string, config BranchAncestryConfig) string {
	result := branchName
	childBranches := config.ChildBranches(branchName)
	sort.Strings(childBranches)
	for _, childBranch := range childBranches {
		result += "\n" + Indent(PrintableBranchTree(childBranch, config))
	}
	return result
}

// PrintableNewBranchPushFlag returns a user printable new branch push flag.
func PrintableNewBranchPushFlag(flag bool) string {
	return strconv.FormatBool(flag)
}

// PrintableOfflineFlag provides a printable version of the given offline flag.
func PrintableOfflineFlag(flag bool) string {
	return strconv.FormatBool(flag)
}
