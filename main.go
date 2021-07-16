package janus

import (
	"fmt"
	"github.com/tendermint/tendermint/libs/service"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"time"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"errors"

	"github.com/tendermint/tendermint/types"
)

var (
	// ErrSignatureRejected indicates that the signature is already locked or locking failed
	ErrSignatureRejected = errors.New("signature rejected")
)

// EtcdSigningWrapper implements PrivValidator by wrapping the etcd locking around a PrivValidator
type EtcdSigningWrapper struct {
	service.BaseService

	privVal         types.PrivValidator
	locker          *Locker
	timeout         time.Duration
	validatorName   string
	lockerEndpoints []string
}

// NewEtcdSigningWrapper returns an instance of EtcdSigningWrapper.
func NewEtcdSigningWrapper(
	logger log.Logger,
	privVal types.PrivValidator,
	timeout time.Duration,
	validatorName string,
	lockerEndpoints []string,
) *EtcdSigningWrapper {
	rs := &EtcdSigningWrapper{
		privVal:         privVal,
		timeout:         timeout,
		validatorName:   validatorName,
		lockerEndpoints: lockerEndpoints,
	}

	rs.BaseService = *service.NewBaseService(logger, "EtcdSigningWrapper", rs)

	return rs
}

// OnStart implements cmn.Service.
func (es *EtcdSigningWrapper) OnStart() error {
	// Start etcd ?
	locker, err := NewLocker(es.lockerEndpoints, es.timeout)
	if err != nil {
		return err
	}
	es.locker = locker

	return nil
}

// OnStop implements cmn.Service.
func (es *EtcdSigningWrapper) OnStop() {
	es.locker.Disconnect()
	es.locker = nil
}

// GetPubKey returns the public key of the validator.
// Implements PrivValidator.
func (es *EtcdSigningWrapper) GetPubKey() (crypto.PubKey, error) {
	return es.privVal.GetPubKey()
}

// SignVote signs a canonical representation of the vote, along with the
// chainID. Implements PrivValidator.
func (es *EtcdSigningWrapper) SignVote(chainID string, vote *tmproto.Vote) error {
	if err := es.privVal.SignVote(chainID, vote); err != nil {
		return fmt.Errorf("error signing vote: %v", err)
	}

	var voteType string
	switch vote.Type {
	case tmproto.PrevoteType:
		voteType = "prevote"
	case tmproto.PrecommitType:
		voteType = "precommit"
	default:
		return errors.New("invalid vote type")
	}

	var err error
	lockAcquired := false
	if vote.Type == tmproto.PrevoteType {
		lockAcquired, err = es.locker.TryLockSetHash(es.validatorName, fmt.Sprintf("vote_%s", voteType), vote.Height, int(vote.Round), vote.BlockID.String())
	} else {
		lockAcquired, err = es.locker.TryLockCheckHash(es.validatorName, fmt.Sprintf("vote_%s", voteType), vote.Height, int(vote.Round), vote.BlockID.String())
	}
	if err != nil {
		vote.Signature = nil
		return err
	}

	if !lockAcquired {
		vote.Signature = nil
		return ErrSignatureRejected
	}

	return nil
}

// SignProposal signs a canonical representation of the proposal, along with
// the chainID. Implements PrivValidator.
func (es *EtcdSigningWrapper) SignProposal(chainID string, proposal *tmproto.Proposal) error {
	if err := es.privVal.SignProposal(chainID, proposal); err != nil {
		return fmt.Errorf("error signing proposal: %v", err)
	}

	lockAcquired, err := es.locker.TryLock(es.validatorName, "proposal", proposal.Height, int(proposal.Round))
	if err != nil {
		proposal.Signature = nil
		return err
	}

	if !lockAcquired {
		proposal.Signature = nil
		return ErrSignatureRejected
	}

	return nil
}
