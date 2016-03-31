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

## Solutions

I don't know of an obvious solution.
Suggestions, with pros and cons, will be added here in time.

