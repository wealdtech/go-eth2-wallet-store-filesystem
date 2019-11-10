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

	"github.com/wealdtech/go-ecodec"
	types "github.com/wealdtech/go-eth2-wallet-types"
)

// StoreAccount stores an account.  It will fail if it cannot store the data.
// Note this will overwrite an existing account with the same ID.  It will not, however, allow multiple accounts with the same
// name to co-exist in the same wallet.
func (s *Store) StoreAccount(wallet types.Wallet, account types.Account, data []byte) error {
	// Ensure the wallet exists
	walletPath := s.walletPath(wallet.Name())
	_, err := os.Stat(walletPath)
	if err != nil {
		return fmt.Errorf("no wallet at %q", walletPath)
	}

	// See if an account with this name already exists
	existingAccount, err := wallet.AccountByName(account.Name())
	if err == nil {
		// It does; they need to have the same ID for us to overwrite it
		if existingAccount.ID().String() != account.ID().String() {
			return fmt.Errorf("account %q already exists", account.Name())
		}
	}

	if len(s.passphrase) > 0 {
		data, err = ecodec.Encrypt(data, s.passphrase)
		if err != nil {
			return err
		}
	}

	// Store the data
	path := s.accountPath(wallet.Name(), account.ID().String())
	return ioutil.WriteFile(filepath.FromSlash(path), data, 0700)
}

// RetrieveAccount retrieves account-level data.  It will fail if it cannot retrieve the data.
func (s *Store) RetrieveAccount(wallet types.Wallet, name string) ([]byte, error) {
	type accountName struct {
		Name string `json:"name"`
	}

	for acc := range s.RetrieveAccounts(wallet) {
		info := &accountName{}
		err := json.Unmarshal(acc, info)
		if err == nil && info.Name == name {
			return acc, nil
		}
	}
	return nil, fmt.Errorf("no account %q", name)
}

// RetrieveAccounts retrieves all account-level data for a wallet.
func (s *Store) RetrieveAccounts(wallet types.Wallet) <-chan []byte {
	ch := make(chan []byte, 1024)
	go func() {
		files, err := ioutil.ReadDir(s.walletPath(wallet.Name()))
		if err == nil {
			for _, file := range files {
				if file.Name() == "_header.json" {
					continue
				}
				data, err := ioutil.ReadFile(s.accountPath(wallet.Name(), strings.TrimSuffix(file.Name(), ".json")))
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
