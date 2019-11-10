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

package filesystem_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	nd "github.com/wealdtech/go-eth2-wallet-nd"
	filesystem "github.com/wealdtech/go-eth2-wallet-store-filesystem"
)

func TestStoreRetrieveAccount(t *testing.T) {
	rand.Seed(time.Now().Unix())
	path := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d", t.Name(), rand.Int31()))
	defer os.RemoveAll(path)
	store := filesystem.New(filesystem.WithLocation(path))
	encryptor := keystorev4.New()

	wallet, err := nd.CreateWallet("test wallet", store, encryptor)
	require.Nil(t, err)
	err = wallet.Unlock(nil)
	require.Nil(t, err)

	accountName := fmt.Sprintf("%d", rand.Int31())
	account, err := wallet.CreateAccount(accountName, []byte{})
	require.Nil(t, err)

	data, err := json.Marshal(account)
	require.Nil(t, err)
	err = store.StoreAccount(wallet, account, data)
	require.Nil(t, err)
	_, err = store.RetrieveAccount(wallet, accountName)
	require.Nil(t, err)

	store.RetrieveWallets()

	_, err = store.RetrieveAccount(wallet, "not present")
	assert.NotNil(t, err)

	_, err = wallet.CreateAccount(accountName, []byte{})
	assert.NotNil(t, err)
}

func TestDuplicateAccounts(t *testing.T) {
	rand.Seed(time.Now().Unix())
	path := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d", t.Name(), rand.Int31()))
	defer os.RemoveAll(path)
	store := filesystem.New(filesystem.WithLocation(path))
	encryptor := keystorev4.New()

	wallet, err := nd.CreateWallet("test wallet", store, encryptor)
	require.Nil(t, err)
	err = wallet.Unlock(nil)
	require.Nil(t, err)

	accountName := fmt.Sprintf("%d", rand.Int31())
	_, err = wallet.CreateAccount(accountName, []byte{})
	require.Nil(t, err)

	// Try to create another account with the same name; should fail
	_, err = wallet.CreateAccount(accountName, []byte{})
	require.NotNil(t, err)
}
