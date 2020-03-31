package store

import (
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/config"
)

type Store struct {
	RawStore   *RawStore
	ChainStore *ChainStore
}

var s *Store

func Init(path string, genesisMilestone config.Milestone) *Store {
	s = new(Store)
	s.RawStore = NewRawStore(path, genesisMilestone)
	s.ChainStore = NewChainStore(path)
	return s
}

func (store *Store) Close() {
}
