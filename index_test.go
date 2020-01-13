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
	"github.com/stretchr/testify/require"
	filesystem "github.com/wealdtech/go-eth2-wallet-store-filesystem"
	"github.com/wealdtech/go-indexer"
)

func TestStoreRetrieveIndex(t *testing.T) {
	rand.Seed(time.Now().Unix())
	path := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d", t.Name(), rand.Int31()))
	//	defer os.RemoveAll(path)
	store := filesystem.New(filesystem.WithLocation(path))

	walletID := uuid.New()
	walletName := "test wallet"
	walletData := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, walletName, walletID.String()))
	accountID := uuid.New()
	accountName := "test account"
	accountData := []byte(fmt.Sprintf(`{"name":%q,"uuid":%q}`, accountName, accountID.String()))

	index := indexer.New()
	index.Add(accountID, accountName)

	err := store.StoreWallet(walletID, walletName, walletData)
	require.Nil(t, err)
	err = store.StoreAccount(walletID, accountID, accountData)
	require.Nil(t, err)

	serializedIndex, err := index.Serialize()
	require.Nil(t, err)
	err = store.StoreAccountsIndex(walletID, serializedIndex)
	require.Nil(t, err)

	fetchedIndex, err := store.RetrieveAccountsIndex(walletID)
	require.Nil(t, err)

	reIndex, err := indexer.Deserialize(fetchedIndex)
	require.Nil(t, err)

	fetchedAccountName, exists := reIndex.Name(accountID)
	require.Equal(t, true, exists)
	require.Equal(t, accountName, fetchedAccountName)

	fetchedAccountID, exists := reIndex.ID(accountName)
	require.Equal(t, true, exists)
	require.Equal(t, accountID, fetchedAccountID)
}
