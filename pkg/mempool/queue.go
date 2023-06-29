package mempool

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
	"github.com/wangjia184/sortedset"
)

type set struct {
	all     *sortedset.SortedSet
	senders map[common.Address]*sortedset.SortedSet
}

func (s *set) getSenderSortedSet(sender common.Address) *sortedset.SortedSet {
	if _, ok := s.senders[sender]; !ok {
		s.senders[sender] = sortedset.New()
	}

	return s.senders[sender]
}

type userOpQueues struct {
	setsByEntryPoint sync.Map
}

func (q *userOpQueues) getEntryPointSet(entryPoint common.Address) *set {
	val, ok := q.setsByEntryPoint.Load(entryPoint)
	if !ok {
		val = &set{
			all:     sortedset.New(),
			senders: make(map[common.Address]*sortedset.SortedSet),
		}
		q.setsByEntryPoint.Store(entryPoint, val)
	}

	return val.(*set)
}

func (q *userOpQueues) AddOp(entryPoint common.Address, op *userop.UserOperation) {
	eps := q.getEntryPointSet(entryPoint)
	sss := eps.getSenderSortedSet(op.Sender)
	key := string(getUniqueKey(entryPoint, op.Sender, op.Nonce))

	eps.all.AddOrUpdate(key, sortedset.SCORE(eps.all.GetCount()), op)
	sss.AddOrUpdate(key, sortedset.SCORE(op.Nonce.Int64()), op)
}

func (q *userOpQueues) GetOps(entryPoint common.Address, sender common.Address) []*userop.UserOperation {
	eps := q.getEntryPointSet(entryPoint)
	sss := eps.getSenderSortedSet(sender)
	nodes := sss.GetByRankRange(-1, -sss.GetCount(), false)
	batch := []*userop.UserOperation{}
	for _, n := range nodes {
		batch = append(batch, n.Value.(*userop.UserOperation))
	}

	return batch
}

func (q *userOpQueues) All(entryPoint common.Address) []*userop.UserOperation {
	eps := q.getEntryPointSet(entryPoint)
	nodes := eps.all.GetByRankRange(1, -1, false)
	batch := []*userop.UserOperation{}
	for _, n := range nodes {
		batch = append(batch, n.Value.(*userop.UserOperation))
	}

	return batch
}

func (q *userOpQueues) RemoveOps(entryPoint common.Address, ops ...*userop.UserOperation) {
	eps := q.getEntryPointSet(entryPoint)
	for _, op := range ops {
		sss := eps.getSenderSortedSet(op.Sender)
		key := string(getUniqueKey(entryPoint, op.Sender, op.Nonce))
		eps.all.Remove(key)
		sss.Remove(key)
	}
}

func newUserOpQueue() *userOpQueues {
	return &userOpQueues{}
}
