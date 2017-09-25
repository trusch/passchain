package state

import (
	"github.com/tendermint/tmlibs/merkle"
)

const (
	accountPrefix = "account::"
	secretPrefix  = "secret::"
)

type State struct {
	Tree merkle.Tree
}

func NewStateFromTree(tree merkle.Tree) *State {
	return &State{tree}
}
