package filesystem

import (
	"path/filepath"
)

func (s *Store) walletPath(walletName string) string {
	return filepath.FromSlash(filepath.Join(s.location, walletName))
}

func (s *Store) walletHeaderPath(walletName string) string {
	return filepath.FromSlash(filepath.Join(s.location, walletName, "_header.json"))
}

func (s *Store) accountPath(walletName string, accountName string) string {
	return filepath.FromSlash(filepath.Join(s.location, walletName, accountName+".json"))
}
