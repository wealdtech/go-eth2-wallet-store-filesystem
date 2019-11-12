// Copyright © 2019 Weald Technology Trading
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/wealdtech/go-ecodec"
)

// StoreAccount stores an account.  It will fail if it cannot store the data.
// Note this will overwrite an existing account with the same ID.  It will not, however, allow multiple accounts with the same
// name to co-exist in the same wallet.
func (s *Store) StoreAccount(walletID uuid.UUID, walletName string, accountID uuid.UUID, accountName string, data []byte) error {
	// Ensure the wallet exists
	walletPath := s.walletPath(walletName)
	_, err := os.Stat(walletPath)
	if err != nil {
		return fmt.Errorf("no wallet at %q", walletPath)
	}

	// See if an account with this name already exists
	existingAccount, err := s.RetrieveAccount(walletID, walletName, accountName)
	if err == nil {
		// It does; they need to have the same ID for us to overwrite it
		info := &struct {
			ID string `json:"id"`
		}{}
		err := json.Unmarshal(existingAccount, info)
		if err != nil {
			return err
		}
		if info.ID != accountID.String() {
			return fmt.Errorf("account %q already exists", accountName)
		}
	}

	if len(s.passphrase) > 0 {
		data, err = ecodec.Encrypt(data, s.passphrase)
		if err != nil {
			return err
		}
	}

	// Store the data
	path := s.accountPath(walletName, accountID.String())
	return ioutil.WriteFile(filepath.FromSlash(path), data, 0700)
}

// RetrieveAccount retrieves account-level data.  It will fail if it cannot retrieve the data.
func (s *Store) RetrieveAccount(walletID uuid.UUID, walletName string, accountName string) ([]byte, error) {
	for acc := range s.RetrieveAccounts(walletID, walletName) {
		info := &struct {
			Name string `json:"name"`
		}{}
		err := json.Unmarshal(acc, info)
		if err == nil && info.Name == accountName {
			return acc, nil
		}
	}
	return nil, fmt.Errorf("no account %q", accountName)
}

// RetrieveAccounts retrieves all account-level data for a wallet.
func (s *Store) RetrieveAccounts(id uuid.UUID, name string) <-chan []byte {
	ch := make(chan []byte, 1024)
	go func() {
		files, err := ioutil.ReadDir(s.walletPath(name))
		if err == nil {
			for _, file := range files {
				if file.Name() == "_header.json" {
					continue
				}
				data, err := ioutil.ReadFile(s.accountPath(name, strings.TrimSuffix(file.Name(), ".json")))
				if err != nil {
					continue
				}
				if len(s.passphrase) > 0 {
					data, err = ecodec.Decrypt(data, s.passphrase)
					if err != nil {
						continue
					}
				}
				ch <- data
			}
		}
		close(ch)
	}()
	return ch
}
