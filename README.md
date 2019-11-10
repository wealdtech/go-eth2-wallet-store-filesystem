# go-eth2-wallet-store-filesystem

[![Tag](https://img.shields.io/github/tag/wealdtech/go-eth2-wallet-store-filesystem.svg)](https://github.com/wealdtech/go-eth2-wallet-store-filesystem/releases/)
[![License](https://img.shields.io/github/license/wealdtech/go-eth2-wallet-store-filesystem.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/wealdtech/go-eth2-wallet-store-filesystem?status.svg)](https://godoc.org/github.com/wealdtech/go-eth2-wallet-store-filesystem)
[![Travis CI](https://img.shields.io/travis/wealdtech/go-eth2-wallet-store-filesystem.svg)](https://travis-ci.org/wealdtech/go-eth2-wallet-store-filesystem)
[![codecov.io](https://img.shields.io/codecov/c/github/wealdtech/go-eth2-wallet-store-filesystem.svg)](https://codecov.io/github/wealdtech/go-eth2-wallet-store-filesystem)
[![Go Report Card](https://goreportcard.com/badge/github.com/wealdtech/go-eth2-wallet-store-filesystem)](https://goreportcard.com/report/github.com/wealdtech/go-eth2-wallet-store-filesystem)

Filesystem-based store for the [Ethereum 2 wallet](https://github.com/wealdtech/go-eth2-wallet).


## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

`go-eth2-wallet-store-filesystem` is a standard Go module which can be installed with:

```sh
go get github.com/wealdtech/go-eth2-wallet-store-filesystem
```

## Usage

In normal operation this module should not be used directly.  Instead, it should be configured to be used as part of [go-eth2-wallet](https://github.com/wealdtech/go-eth2-wallet).

The filesystem store has the following options:

  - `location`: the base directory in which to store wallets.  If this is not configured it defaults to a sensible operating system-specific values:
    - for Linux: $HOME/.config/ethereum2/wallets
    - for OSX: $HOME/Library/Application Support/ethereum2/wallets
    - for Windows: %APPDATA%\ethereum2\wallets
  - `passphrase`: a key used to encrypt all data written to the store.  If this is not configured data is written to the store unencrypted (although wallet- and account-specific private information may be protected by their own passphrases)

### Example

```go
package main

import (
	e2wallet "github.com/wealdtech/go-eth2-wallet"
	filesystem "github.com/wealdtech/go-eth2-wallet-store-filesystem"
)

func main() {
    // Set up and use a simple store
    store := filesystem.New()
    e2wallet.UseStore(store)

    // Set up and use an encrypted store
    store := filesystem.New(filesystem.WithPassphrase([]byte("my secret")))
    e2wallet.UseStore(store)

    // Set up and use an encrypted store at a custom location
    store := filesystem.New(filesystem.WithPassphrase([]byte("my secret")), filesystem.WithLocation("/home/user/wallets"))
    e2wallet.UseStore(store)

    // Use e2wallet operations as normal.
}
```

## Maintainers

Jim McDonald: [@mcdee](https://github.com/mcdee).

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/wealdtech/go-eth2-wallet-store-filesystem/issues).

## License

[Apache-2.0](LICENSE) © 2019 Weald Technology Trading Ltd
