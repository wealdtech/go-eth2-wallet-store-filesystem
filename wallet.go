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
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// StoreWallet stores wallet-level data.  It will fail if it cannot store the data.
// Note that this will overwrite any existing data; it is up to higher-level functions to check for the presence of a wallet with
// the wallet name and handle clashes accordingly.
func (s *Store) StoreWallet(walletID uuid.UUID, _ string, data []byte) error {
	if err := s.ensureWalletPathExists(walletID); err != nil {
		return errors.Wrap(err, "wallet path does not exist")
	}
	data, err := s.encryptIfRequired(data)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt wallet")
	}

	return os.WriteFile(s.walletHeaderPath(walletID), data, 0o600)
}

// RetrieveWallet retrieves wallet-level data.  It will fail if it cannot retrieve the data.
func (s *Store) RetrieveWallet(walletName string) ([]byte, error) {
	for data := range s.RetrieveWallets() {
		info := &struct {
			Name string `json:"name"`
		}{}
		err := json.Unmarshal(data, info)
		if err == nil && info.Name == walletName {
			return data, nil
		}
	}

	return nil, errors.New("wallet not found")
}

// RetrieveWalletByID retrieves wallet-level data.  It will fail if it cannot retrieve the data.
func (s *Store) RetrieveWalletByID(walletID uuid.UUID) ([]byte, error) {
	for data := range s.RetrieveWallets() {
		info := &struct {
			ID uuid.UUID `json:"uuid"`
		}{}
		err := json.Unmarshal(data, info)
		if err == nil && info.ID == walletID {
			return data, nil
		}
	}

	return nil, errors.New("wallet not found")
}

// RetrieveWallets retrieves wallet-level data for all wallets.
func (s *Store) RetrieveWallets() <-chan []byte {
	ch := make(chan []byte, 1024)
	go func() {
		defer close(ch)
		dirs, err := os.ReadDir(s.location)
		if err != nil {
			return
		}
		for _, dir := range dirs {
			if !dir.IsDir() {
				continue
			}
			walletID, err := uuid.Parse(dir.Name())
			if err != nil {
				continue
			}
			data, err := os.ReadFile(s.walletHeaderPath(walletID))
			if err != nil {
				continue
			}
			data, err = s.decryptIfRequired(data)
			if err != nil {
				continue
			}
			ch <- data
		}
	}()

	return ch
}
