package wtf

import (
	"errors"

	"github.com/coreos/etcd/etcdserver/etcdserverpb"
	"golang.org/x/net/context"
)

// Server provides an alternate implementation of the etcd KVServer.
type Server struct{}

// Range implements KVServer.
func (Server) Range(context.Context, *etcdserverpb.RangeRequest) (*etcdserverpb.RangeResponse, error) {
	return nil, errors.New("not implemented")
}

// Put implements KVServer.
func (Server) Put(context.Context, *etcdserverpb.PutRequest) (*etcdserverpb.PutResponse, error) {
	return nil, errors.New("not implemented")
}

// DeleteRange implements KVServer.
func (Server) DeleteRange(context.Context, *etcdserverpb.DeleteRangeRequest) (*etcdserverpb.DeleteRangeResponse, error) {
	return nil, errors.New("not implemented")
}

// Txn implements KVServer.
func (Server) Txn(context.Context, *etcdserverpb.TxnRequest) (*etcdserverpb.TxnResponse, error) {
	return nil, errors.New("not implemented")
}

// Compact implements KVServer.
func (Server) Compact(context.Context, *etcdserverpb.CompactionRequest) (*etcdserverpb.CompactionResponse, error) {
	return nil, errors.New("not implemented")
}

var _ etcdserverpb.KVServer = Server{}
