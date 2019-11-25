package filesystem

import (
	"path/filepath"

	"github.com/google/uuid"
)

func (s *Store) walletPath(walletID uuid.UUID) string {
	return filepath.FromSlash(filepath.Join(s.location, walletID.String()))
}

func (s *Store) walletHeaderPath(walletID uuid.UUID) string {
	return filepath.FromSlash(filepath.Join(s.location, walletID.String(), walletID.String()))
}

func (s *Store) accountPath(walletID uuid.UUID, accountID uuid.UUID) string {
	return filepath.FromSlash(filepath.Join(s.location, walletID.String(), accountID.String()))
}
