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
	"github.com/shibukawa/configdir"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

// options are the options for the filesystem store.
type options struct {
	passphrase []byte
	location   string
}

// Option gives options to New
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// WithPassphrase sets the encryption for the store.
func WithPassphrase(passphrase []byte) Option {
	return optionFunc(func(o *options) {
		o.passphrase = passphrase
	})
}

// WithLocation sets the on-filesystem location for the store.
func WithLocation(b string) Option {
	return optionFunc(func(o *options) {
		o.location = b
	})
}

// Store is the store for the wallet.
type Store struct {
	location   string
	passphrase []byte
}

func defaultLocation() string {
	configDirs := configdir.New("ethereum2", "wallets")
	return configDirs.QueryFolders(configdir.Global)[0].Path
}

// New creates a new filesystem store.
// If the path is not supplied a default path is used.
func New(opts ...Option) wtypes.Store {
	options := options{
		location: defaultLocation(),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	return &Store{
		location:   options.location,
		passphrase: options.passphrase,
	}
}

// Name returns the name of this store.
func (s *Store) Name() string {
	return "filesystem"
}

// Location returns the location of this store.
func (s *Store) Location() string {
	return s.location
}
