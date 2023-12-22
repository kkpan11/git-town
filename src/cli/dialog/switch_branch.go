package dialog

import (
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

func NewBuilder(lineage configdomain.Lineage) Builder {
	return Builder{
		Entries: ModalSelectEntries{},
		Lineage: lineage,
	}
}

// queryBranch lets the user select a new branch via a visual dialog.
// Indicates via `validSelection` whether the user made a valid selection.
func SwitchBranch(roots gitdomain.LocalBranchNames, selected gitdomain.LocalBranchName, lineage configdomain.Lineage) (selection gitdomain.LocalBranchName, validSelection bool, err error) {
	builder := NewBuilder(lineage)
	err = builder.CreateEntries(roots, selected)
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), false, err
	}
	choice, err := ModalSelect(builder.Entries, selected.String())
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), false, err
	}
	if choice == nil {
		return gitdomain.EmptyLocalBranchName(), false, nil
	}
	return gitdomain.NewLocalBranchName(*choice), true, nil
}

// Builder builds up the switch-branch dialog entries.
type Builder struct {
	Entries ModalSelectEntries
	Lineage configdomain.Lineage
}

// AddEntryAndChildren adds the given branch and all its child branches to the given entries collection.
func (self *Builder) AddEntryAndChildren(branch gitdomain.LocalBranchName, indent int) error {
	self.Entries = append(self.Entries, ModalSelectEntry{
		Text:  strings.Repeat("  ", indent) + branch.String(),
		Value: branch.String(),
	})
	var err error
	for _, child := range self.Lineage.Children(branch) {
		err = self.AddEntryAndChildren(child, indent+1)
		if err != nil {
			return err
		}
	}
	return nil
}

// createEntries provides all the entries for the branch dialog.
func (self *Builder) CreateEntries(roots gitdomain.LocalBranchNames, selected gitdomain.LocalBranchName) error {
	var err error
	for _, root := range roots {
		err = self.AddEntryAndChildren(root, 0)
		if err != nil {
			return err
		}
	}
	if len(self.Entries) == 0 {
		self.Entries = append(self.Entries, ModalSelectEntry{
			Text:  string(selected),
			Value: string(selected),
		})
	}
	return nil
}
