# ddbsync

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/joshmyers/ddbsync)
[![Build Status](https://circleci.com/gh/joshmyers/ddbsync/tree/master.svg?style=svg)](https://circleci.com/gh/joshmyers/ddbsync/tree/master)
[![Coverage Status](https://coveralls.io/repos/joshmyers/ddbsync/badge.svg?branch=master)](https://coveralls.io/r/joshmyers/ddbsync?branch=master)

DynamoDB/sync

This package is designed to emulate the behaviour of `pkg/sync` on top of Amazon's DynamoDB. If you need a distributed locking mechanism, consider using this package and DynamoDB before standing up paxos or Zookeeper.


## Dependency Management

This project uses [Glide](https://github.com/Masterminds/glide) to manage it's dependencies, and does not commit the vendor directory. As a result this is not `go get` installable, please refer to the Glide docs for Glide install instructions, and make sure to `glide install` to create a `vendor` directory before attempting to build against this repo.

## Usage

Create a DynamoDB table named *Locks*.

```bash
$ export AWS_ACCESS_KEY=access
$ export AWS_SECRET_KEY=secret
```

```go
// ./main.go

package main

import(
	"time"
	"github.com/joshmyers/ddbsync"
)

func main() {
	m := new(ddbsync.Mutex)
	m.Name = "some-name"
	m.TTL = int64(10 * time.Second)
	m.Lock()
	defer m.Unlock()
	// do important work here
	return
}
```

```bash
$ git clone http://github.com/joshmyers/ddbsync && cd ddbsync
$ glide install
$ go run main.go
```

## Related

[ddbsync](https://github.com/ryandotsmith/ddbsync)
[lock-smith](https://github.com/ryandotsmith/lock-smith)
