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
	"io/ioutil"

	"github.com/google/uuid"
)

// StoreAccountsIndex stores the account index.
func (s *Store) StoreAccountsIndex(walletID uuid.UUID, data []byte) error {
	if err := s.ensureWalletPathExists(walletID); err != nil {
		return err
	}

	data, err := s.encryptIfRequired(data)
	if err != nil {
		return err
	}

	path := s.walletIndexPath(walletID)
	return ioutil.WriteFile(path, data, 0700)
}

// RetrieveAccountsIndex retrieves the account index.
func (s *Store) RetrieveAccountsIndex(walletID uuid.UUID) ([]byte, error) {
	path := s.walletIndexPath(walletID)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return s.decryptIfRequired(data)
}
