package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/messages"
)

const (
	shipDeleteTrackingBranchTitle = `Ship delete tracking branch`
	ShipDeleteTrackingBranchHelp  = `
Should "git town ship" delete the tracking branch?
You want to disable this if your code hosting platform
(GitHub, GitLab, etc) deletes head branches when
merging pull requests through its UI.

`
)

const (
	ShipDeleteTrackingBranchEntryYes shipDeleteTrackingBranchEntry = `yes, "git town ship" should delete tracking branches`
	ShipDeleteTrackingBranchEntryNo  shipDeleteTrackingBranchEntry = `no, my code hosting platform deletes tracking branches`
)

func ShipDeleteTrackingBranch(existing configdomain.ShipDeleteTrackingBranch, inputs components.TestInput) (configdomain.ShipDeleteTrackingBranch, bool, error) {
	entries := []shipDeleteTrackingBranchEntry{
		ShipDeleteTrackingBranchEntryYes,
		ShipDeleteTrackingBranchEntryNo,
	}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := components.RadioList(list.NewEntries(entries...), defaultPos, shipDeleteTrackingBranchTitle, ShipDeleteTrackingBranchHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf(messages.ShipDeletesTrackingBranches, components.FormattedSelection(selection.Short(), aborted))
	return selection.ShipDeleteTrackingBranch(), aborted, err
}

type shipDeleteTrackingBranchEntry string

func (self shipDeleteTrackingBranchEntry) ShipDeleteTrackingBranch() configdomain.ShipDeleteTrackingBranch {
	switch self {
	case ShipDeleteTrackingBranchEntryYes:
		return configdomain.ShipDeleteTrackingBranch(true)
	case ShipDeleteTrackingBranchEntryNo:
		return configdomain.ShipDeleteTrackingBranch(false)
	}
	panic("unhandled shipDeleteTrackingBranchEntry: " + self)
}

func (self shipDeleteTrackingBranchEntry) Short() string {
	start, _, _ := strings.Cut(self.String(), ",")
	return start
}

func (self shipDeleteTrackingBranchEntry) String() string {
	return string(self)
}
