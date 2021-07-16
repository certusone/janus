package janus

import (
	"context"
	"fmt"
	etcd "go.etcd.io/etcd/v3/clientv3"
	"time"
)

type (
	// Locker can be used to lock signature steps in Tendermint using etcd to prevent double signing
	Locker struct {
		client  *etcd.Client
		timeout time.Duration
	}
)

const (
	lockKeyFormat = "val_lock/h_%d/r_%d/%s"
)

// NewLocker returns an instance of Locker
func NewLocker(endpoints []string, timeout time.Duration) (*Locker, error) {
	c, err := etcd.New(etcd.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})

	// etcd etcd >= v3.2.10, grpc/grpc-go >= v1.7.3
	if err != nil {
		return &Locker{timeout: timeout}, err
	}

	return &Locker{
		c,
		timeout,
	}, nil
}

// TryLock tries to lock a specific signature step
// Returns false if the signature step is already locked
func (l *Locker) TryLock(validatorName, lockType string, height int64, round int) (bool, error) {
	key := fmt.Sprintf(lockKeyFormat, height, round, lockType)

	cmp := etcd.Compare(etcd.CreateRevision(key), "=", 0)
	put := etcd.OpPut(key, validatorName)

	ctx, _ := context.WithTimeout(context.Background(), l.timeout)
	resp, err := l.client.Txn(ctx).If(cmp).Then(put).Commit()
	if err != nil {
		return false, err
	}

	return resp.Succeeded, nil
}

// TryLockSetHash tries to lock a specific signature step for this h/r pair.
// This should be used for the prevote step.
// Returns false if the signature step is already locked
func (l *Locker) TryLockSetHash(validatorName, lockType string, height int64, round int, blockHash string) (bool, error) {
	key := fmt.Sprintf(lockKeyFormat, height, round, lockType)

	cmp := etcd.Compare(etcd.CreateRevision(key), "=", 0)
	put := etcd.OpPut(key, validatorName)

	ctx, _ := context.WithTimeout(context.Background(), l.timeout)
	resp, err := l.client.Txn(ctx).If(cmp).Then(put).Commit()
	if err != nil {
		return false, err
	}

	return resp.Succeeded, nil
}

// TryLockCheckHash tries to lock a specific signature step.
// This should be used for the precommit step.
// Returns false if the signature step is already locked
func (l *Locker) TryLockCheckHash(validatorName, lockType string, height int64, round int, blockHash string) (bool, error) {
	key := fmt.Sprintf(lockKeyFormat, height, round, lockType)

	cmp := etcd.Compare(etcd.CreateRevision(key), "=", 0)
	put := etcd.OpPut(key, validatorName)

	ctx, _ := context.WithTimeout(context.Background(), l.timeout)
	resp, err := l.client.Txn(ctx).If(cmp).Then(put).Commit()
	if err != nil {
		return false, err
	}

	return resp.Succeeded, nil
}

// Disconnect closes the underlying etcd connection
func (l *Locker) Disconnect() {
	l.client.Close()
}
