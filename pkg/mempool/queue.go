package mempool

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
	"github.com/wangjia184/sortedset"
)

type set struct {
	all      *sortedset.SortedSet
	entities map[common.Address]*sortedset.SortedSet
}

func (s *set) getEntitiesSortedSet(entity common.Address) *sortedset.SortedSet {
	if _, ok := s.entities[entity]; !ok {
		s.entities[entity] = sortedset.New()
	}

	return s.entities[entity]
}

type userOpQueues struct {
	setsByEntryPoint sync.Map
}

func (q *userOpQueues) getEntryPointSet(entryPoint common.Address) *set {
	val, ok := q.setsByEntryPoint.Load(entryPoint)
	if !ok {
		val = &set{
			all:      sortedset.New(),
			entities: make(map[common.Address]*sortedset.SortedSet),
		}
		q.setsByEntryPoint.Store(entryPoint, val)
	}

	return val.(*set)
}

func (q *userOpQueues) AddOp(entryPoint common.Address, op *userop.UserOperation) {
	eps := q.getEntryPointSet(entryPoint)
	key := string(getUniqueKey(entryPoint, op.Sender, op.Nonce))

	eps.all.AddOrUpdate(key, sortedset.SCORE(eps.all.GetCount()), op)
	eps.getEntitiesSortedSet(op.Sender).AddOrUpdate(key, sortedset.SCORE(op.Nonce.Int64()), op)
	if factory := op.GetFactory(); factory != common.HexToAddress("0x") {
		fss := eps.getEntitiesSortedSet(factory)
		fss.AddOrUpdate(key, sortedset.SCORE(fss.GetCount()), op)
	}
	if paymaster := op.GetPaymaster(); paymaster != common.HexToAddress("0x") {
		pss := eps.getEntitiesSortedSet(paymaster)
		pss.AddOrUpdate(key, sortedset.SCORE(pss.GetCount()), op)
	}
}

func (q *userOpQueues) GetOps(entryPoint common.Address, entity common.Address) []*userop.UserOperation {
	eps := q.getEntryPointSet(entryPoint)
	ess := eps.getEntitiesSortedSet(entity)
	nodes := ess.GetByRankRange(-1, -ess.GetCount(), false)
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
		key := string(getUniqueKey(entryPoint, op.Sender, op.Nonce))
		eps.all.Remove(key)
		eps.getEntitiesSortedSet(op.Sender).Remove(key)
		eps.getEntitiesSortedSet(op.GetFactory()).Remove(key)
		eps.getEntitiesSortedSet(op.GetPaymaster()).Remove(key)
	}
}

func newUserOpQueue() *userOpQueues {
	return &userOpQueues{}
}
