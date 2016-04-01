# wtf [![Build Status](https://travis-ci.org/peterbourgon/wtf.svg?branch=master)](https://travis-ci.org/peterbourgon/wtf)

wtf is a package that demonstrates what I think is a serious problem with how Go vendoring currently works.
As a result of this problem, the package fails to build.

## Problem

1. package wtf is a library, designed to be imported by other code.
   Vendoring is the responsibility of binary (i.e. package main) authors.
   Therefore, package wtf doesn't (and shouldn't) vendor its dependencies.

2. package wtf provides a type that should satisfy a third party interface.
   Specifically, [etcdserverpb.KVServer](https://godoc.org/github.com/coreos/etcd/etcdserver/etcdserverpb#KVServer).

3. Methods in KVServer take [context.Context](https://godoc.org/golang.org/x/net/context) as a parameter.
   Therefore, wtf.go imports `"golang.org/x/net/context"`.
   Note that etcdserverpb [imports the same literal path](https://github.com/coreos/etcd/blob/1d698f093f490cade87b7e5b33e8525367f49f0c/etcdserver/etcdserverpb/rpc.pb.go#L18).

4. Note that etcd has [a vendor folder](https://github.com/coreos/etcd/tree/master/vendor) in the root of its repository, and that vendor folder includes golang.org/x/net/context.

At build time, the go toolchain resolves the import path `"golang.org/x/net/context"` as specified in etcdserverpb to the copy stored at $GOPATH/src/github.com/coreos/etcd/vendor/golang.org/x/net/context.
But, it resolves the same import path as specified in wtf.go to _a different copy_, by default the plain $GOPATH/src/golang.org/x/net/context.
Note that it's _impossible_ for the import path specified in wtf.go to resolve to the vendored copy in github.com/coreos/etcd/vendor.

Because the same import path resolves to two different concrete packages, and no type coercion is possible, the build fails:

```
./wtf.go:38: cannot use Server literal (type Server) as type etcdserverpb.KVServer in assignment:
	Server does not implement etcdserverpb.KVServer (wrong type for Compact method)
		have Compact("golang.org/x/net/context".Context, *etcdserverpb.CompactionRequest) (*etcdserverpb.CompactionResponse, error)
		want Compact("github.com/coreos/etcd/vendor/golang.org/x/net/context".Context, *etcdserverpb.CompactionRequest) (*etcdserverpb.CompactionResponse, error)
```

## More context

- [My golang-nuts post](https://groups.google.com/d/msg/golang-nuts/AnMr9NL6dtc/UnyUUKcMCAAJ)
- [A similar golang-dev post](https://groups.google.com/forum/#!msg/golang-dev/WebP4dLV1b0/Lhk4hpwJEgAJ)
- https://github.com/coreos/etcd/issues/4913

## Solutions

### Restructure etcd repos

coreos/etcd could be split into two repos:

- coreos/etcdbin, containing all `package main` binaries, with vendoring
- coreos/etcdlib, containing all non-`package main` libraries, without vendoring

This imposes a burden on the etcd developers, but appears to be the best solution for consumers.

### Use a vendoring tool to build

Daniel Theophanes notes that

>  The package won't [build] with a git pull or `go get`, but it will compile correctly when you run `govendor get github.com/peterbourgon/wtf`. Ensure you've run `go get -u github.com/kardianos/govendor` first.

Requiring a third-party build tool to build a project that just happens to depend on a package that uses vendoring is obviously not an ideal outcome. Go developers have the expectation that `go get` will work in the general case.

### Others

I'll be adding more solutions, with pros and cons, as they are discovered.
