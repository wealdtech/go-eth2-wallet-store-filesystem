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

package filesystem_test

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	filesystem "github.com/wealdtech/go-eth2-wallet-store-filesystem"
)

func TestStoreRetrieveEncryptedWallet(t *testing.T) {
	rand.Seed(time.Now().Unix())
	path := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d", t.Name(), rand.Int31()))
	defer os.RemoveAll(path)
	store := filesystem.New(filesystem.WithLocation(path), filesystem.WithPassphrase([]byte("test")))

	walletID := uuid.New()
	walletName := "test"
	data := []byte(fmt.Sprintf(`{"uuid":%q,"name":%q}`, walletID, walletName))

	err := store.StoreWallet(walletID, walletName, data)
	require.Nil(t, err)
	retData, err := store.RetrieveWallet(walletName)
	require.Nil(t, err)
	assert.Equal(t, data, retData)

	wallets := false
	for range store.RetrieveWallets() {
		wallets = true
	}
	assert.True(t, wallets)

	store.RetrieveWallets()
}

func TestStoreRetrieveEncryptedAccount(t *testing.T) {
	rand.Seed(time.Now().Unix())
	path := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d", t.Name(), rand.Int31()))
	defer os.RemoveAll(path)
	store := filesystem.New(filesystem.WithLocation(path), filesystem.WithPassphrase([]byte("test")))

	walletID := uuid.New()
	walletName := "test wallet"
	walletData := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, walletName, walletID.String()))
	accountID := uuid.New()
	accountName := "test account"
	accountData := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, accountName, accountID.String()))

	err := store.StoreWallet(walletID, walletName, walletData)
	require.Nil(t, err)

	err = store.StoreAccount(walletID, accountID, accountData)
	require.Nil(t, err)
	retData, err := store.RetrieveAccount(walletID, accountID)
	require.Nil(t, err)
	require.Equal(t, accountData, retData)

	accounts := false
	for range store.RetrieveAccounts(walletID) {
		accounts = true
	}
	assert.True(t, accounts)
}

func TestBadWalletKey(t *testing.T) {
	rand.Seed(time.Now().Unix())
	path := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d", t.Name(), rand.Int31()))
	defer os.RemoveAll(path)
	store := filesystem.New(filesystem.WithLocation(path), filesystem.WithPassphrase([]byte("test")))

	walletID := uuid.New()
	walletName := "test wallet"
	data := []byte(fmt.Sprintf(`{"uuid":%q,"name":%q}`, walletID, walletName))

	err := store.StoreWallet(walletID, walletName, data)
	require.Nil(t, err)

	// Open wallet with store with different key; should fail
	store = filesystem.New(filesystem.WithLocation(path), filesystem.WithPassphrase([]byte("badkey")))
	_, err = store.RetrieveWallet(walletName)
	require.NotNil(t, err)
}
