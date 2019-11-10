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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/wealdtech/go-ecodec"
	types "github.com/wealdtech/go-eth2-wallet-types"
)

// StoreWallet stores wallet-level data.  It will fail if it cannot store the data.
// Note that this will overwrite any existing data; it is up to higher-level functions to check for the presence of a wallet with
// the wallet name and handle clashes accordingly.
func (s *Store) StoreWallet(wallet types.Wallet, data []byte) error {
	path := s.walletPath(wallet.Name())
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0700)
		if err != nil {
			return fmt.Errorf("failed to create wallet at %s", path)
		}
	}

	if len(s.passphrase) > 0 {
		data, err = ecodec.Encrypt(data, s.passphrase)
		if err != nil {
			return errors.Wrap(err, "failed to encrypt wallet")
		}
	}
	return ioutil.WriteFile(s.walletHeaderPath(wallet.Name()), data, 0700)
}

// RetrieveWallet retrieves wallet-level data.  It will fail if it cannot retrieve the data.
func (s *Store) RetrieveWallet(walletName string) ([]byte, error) {
	path := s.walletPath(walletName)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("no wallet at %s", path)
	}

	data, err := ioutil.ReadFile(s.walletHeaderPath(walletName))
	if err != nil {
		return nil, err
	}

	if len(s.passphrase) > 0 {
		data, err = ecodec.Decrypt(data, s.passphrase)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decrypt wallet")
		}
	}
	return data, nil
}

// RetrieveWallets retrieves wallet-level data for all wallets.
func (s *Store) RetrieveWallets() <-chan []byte {
	ch := make(chan []byte, 1024)
	go func() {
		dirs, err := ioutil.ReadDir(s.location)
		if err == nil {
			for _, dir := range dirs {
				data, err := s.RetrieveWallet(dir.Name())
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
