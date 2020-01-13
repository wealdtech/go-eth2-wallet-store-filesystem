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
	"testing"

	"github.com/stretchr/testify/assert"
	filesystem "github.com/wealdtech/go-eth2-wallet-store-filesystem"
)

func TestNew(t *testing.T) {
	store := filesystem.New()
	assert.Equal(t, "filesystem", store.Name())
	store = filesystem.New(filesystem.WithLocation("test"))
	assert.Equal(t, "filesystem", store.Name())
	store = filesystem.New(filesystem.WithLocation("test"), filesystem.WithPassphrase([]byte("secret")))
	assert.Equal(t, "filesystem", store.Name())
}
