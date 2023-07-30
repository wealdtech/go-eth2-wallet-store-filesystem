// Copyright 2019 - 2023 Weald Technology Trading
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
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// StoreAccount stores an account.  It will fail if it cannot store the data.
// Note this will overwrite an existing account with the same ID.  It will not, however, allow multiple accounts with the same
// name to co-exist in the same wallet.
func (s *Store) StoreAccount(walletID uuid.UUID, accountID uuid.UUID, data []byte) error {
	// Ensure the wallet exists.
	_, err := s.RetrieveWalletByID(walletID)
	if err != nil {
		return errors.Wrap(err, "unable to retrieve wallet")
	}

	data, err = s.encryptIfRequired(data)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt account")
	}
	path := s.accountPath(walletID, accountID)

	return os.WriteFile(filepath.FromSlash(path), data, 0o600)
}

// RetrieveAccount retrieves account-level data.  It will return an error if it cannot retrieve the data.
func (s *Store) RetrieveAccount(walletID uuid.UUID, accountID uuid.UUID) ([]byte, error) {
	path := s.accountPath(walletID, accountID)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "account not found")
	}
	data, err = s.decryptIfRequired(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt account")
	}

	return data, nil
}

// RetrieveAccounts retrieves all account-level data for a wallet.
func (s *Store) RetrieveAccounts(walletID uuid.UUID) <-chan []byte {
	ch := make(chan []byte, 1024)
	go func() {
		defer close(ch)
		files, err := os.ReadDir(s.walletPath(walletID))
		if err != nil {
			return
		}

		walletName := walletID.String()
		for _, file := range files {
			switch file.Name() {
			case walletName, "index", "batch":
				// Not accounts.
				continue
			default:
				accountID, err := uuid.Parse(file.Name())
				if err != nil {
					continue
				}
				data, err := os.ReadFile(s.accountPath(walletID, accountID))
				if err != nil {
					continue
				}
				data, err = s.decryptIfRequired(data)
				if err != nil {
					continue
				}
				ch <- data
			}
		}
	}()

	return ch
}
