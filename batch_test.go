// Copyright 2023 Weald Technology Trading.
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
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	filesystem "github.com/wealdtech/go-eth2-wallet-store-filesystem"
	e2wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

func TestStoreRetrieveBatch(t *testing.T) {
	ctx := context.Background()

	path := filepath.Join(os.TempDir(), fmt.Sprintf("TestStoreRetrieveWallet-%d", rand.Int31()))
	defer os.RemoveAll(path)
	store := filesystem.New(filesystem.WithLocation(path))

	walletID := uuid.New()
	walletName := "test wallet"
	data := []byte(fmt.Sprintf(`{"uuid":%q,"name":%q}`, walletID, walletName))
	require.Nil(t, store.StoreWallet(walletID, walletName, data))

	batchData := []byte(`{"test":true}`)
	require.NoError(t, store.(e2wtypes.BatchStorer).StoreBatch(ctx, walletID, walletName, batchData))

	retrievedBatchData, err := store.(e2wtypes.BatchRetriever).RetrieveBatch(ctx, walletID)
	require.NoError(t, err)
	require.Equal(t, batchData, retrievedBatchData)
}

func TestStoreBatchNonExistentWallet(t *testing.T) {
	ctx := context.Background()

	path := filepath.Join(os.TempDir(), fmt.Sprintf("TestStoreRetrieveWallet-%d", rand.Int31()))
	defer os.RemoveAll(path)
	store := filesystem.New(filesystem.WithLocation(path))

	walletID := uuid.New()
	walletName := "test wallet"

	batchData := []byte(`{"test":true}`)
	require.ErrorContains(t, store.(e2wtypes.BatchStorer).StoreBatch(ctx, walletID, walletName, batchData), "wallet not found")
}

func TestRetrieveBatchNonExistentWallet(t *testing.T) {
	ctx := context.Background()

	path := filepath.Join(os.TempDir(), fmt.Sprintf("TestRetrieveNonExistentWallet-%d", rand.Int31()))
	defer os.RemoveAll(path)
	store := filesystem.New(filesystem.WithLocation(path))

	walletID := uuid.New()

	_, err := store.(e2wtypes.BatchRetriever).RetrieveBatch(ctx, walletID)
	require.ErrorContains(t, err, "wallet not found")
}

func TestRetrieveNonExistentBatch(t *testing.T) {
	ctx := context.Background()

	path := filepath.Join(os.TempDir(), fmt.Sprintf("TestRetrieveNonExistentWallet-%d", rand.Int31()))
	defer os.RemoveAll(path)
	store := filesystem.New(filesystem.WithLocation(path))

	walletID := uuid.New()
	walletName := "test wallet"
	data := []byte(fmt.Sprintf(`{"uuid":%q,"name":%q}`, walletID, walletName))
	require.Nil(t, store.StoreWallet(walletID, walletName, data))

	_, err := store.(e2wtypes.BatchRetriever).RetrieveBatch(ctx, walletID)
	require.ErrorContains(t, err, "no such file or directory")
}
