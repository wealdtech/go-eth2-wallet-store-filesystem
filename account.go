// Copyright 2019, 2020 Weald Technology Trading
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
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/google/uuid"
)

// StoreAccount stores an account.  It will fail if it cannot store the data.
// Note this will overwrite an existing account with the same ID.  It will not, however, allow multiple accounts with the same
// name to co-exist in the same wallet.
func (s *Store) StoreAccount(walletID uuid.UUID, accountID uuid.UUID, data []byte) error {
	// Ensure the wallet exists
	_, err := s.RetrieveWalletByID(walletID)
	if err != nil {
		return errors.New("unknown wallet")
	}

	data, err = s.encryptIfRequired(data)
	if err != nil {
		return err
	}
	path := s.accountPath(walletID, accountID)
	return ioutil.WriteFile(filepath.FromSlash(path), data, 0700)
}

// RetrieveAccount retrieves account-level data.  It will return an error if it cannot retrieve the data.
func (s *Store) RetrieveAccount(walletID uuid.UUID, accountID uuid.UUID) ([]byte, error) {
	path := s.accountPath(walletID, accountID)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		// TODO handle specific errors such as not found
		return nil, errors.New("account not found")
	}
	data, err = s.decryptIfRequired(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// RetrieveAccounts retrieves all account-level data for a wallet.
func (s *Store) RetrieveAccounts(walletID uuid.UUID) <-chan []byte {
	ch := make(chan []byte, 1024)
	go func() {
		files, err := ioutil.ReadDir(s.walletPath(walletID))
		if err == nil {
			for _, file := range files {
				if file.Name() == walletID.String() {
					continue
				}
				if file.Name() == "index" {
					continue
				}
				accountID, err := uuid.Parse(file.Name())
				if err != nil {
					continue
				}
				data, err := ioutil.ReadFile(s.accountPath(walletID, accountID))
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
		close(ch)
	}()
	return ch
}
