// Copyright 2019 - 2023 Weald Technology Trading.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filesystem

import (
	"fmt"
	"os"
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

func (s *Store) walletIndexPath(walletID uuid.UUID) string {
	return filepath.FromSlash(filepath.Join(s.walletPath(walletID), "index"))
}

func (s *Store) walletBatchPath(walletID uuid.UUID) string {
	return filepath.FromSlash(filepath.Join(s.walletPath(walletID), "batch"))
}

func (s *Store) ensureWalletPathExists(walletID uuid.UUID) error {
	path := s.walletPath(walletID)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0o700)
		if err != nil {
			return fmt.Errorf("failed to create wallet directory at %s", path)
		}
	}

	return nil
}
