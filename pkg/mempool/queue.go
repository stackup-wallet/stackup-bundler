package mempool

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
	"github.com/wangjia184/sortedset"
)

type maxBatchSize int

const defaultMaxBatchSize = 10

type setsByEntryPoint map[string]*sortedset.SortedSet

type userOpQueues struct {
	maxBatchSize
	setsByEntryPoint
}

func (q *userOpQueues) findOrCreateSet(entryPoint common.Address) *sortedset.SortedSet {
	ep := entryPoint.String()
	if _, ok := q.setsByEntryPoint[ep]; !ok {
		q.setsByEntryPoint[ep] = sortedset.New()
	}

	return q.setsByEntryPoint[ep]
}

func (q *userOpQueues) AddOp(entryPoint common.Address, op *userop.UserOperation) bool {
	set := q.findOrCreateSet(entryPoint)

	return set.AddOrUpdate(op.Sender.String(), sortedset.SCORE(op.MaxPriorityFeePerGas.Int64()), op)
}

func (q *userOpQueues) GetOp(entryPoint common.Address, sender common.Address) *userop.UserOperation {
	set := q.findOrCreateSet(entryPoint)
	node := set.GetByKey(sender.String())
	if node == nil {
		return nil
	}

	return node.Value.(*userop.UserOperation)
}

func (q *userOpQueues) Next(entryPoint common.Address) []*userop.UserOperation {
	nodes := q.findOrCreateSet(entryPoint).GetByRankRange(-1, -defaultMaxBatchSize, false)
	batch := []*userop.UserOperation{}
	for _, n := range nodes {
		batch = append(batch, n.Value.(*userop.UserOperation))
	}

	return batch
}

func (q *userOpQueues) RemoveOps(entryPoint common.Address, senders []common.Address) {
	set := q.findOrCreateSet(entryPoint)
	for _, s := range senders {
		set.Remove(s.String())
	}
}

func newUserOpQueue() *userOpQueues {
	return &userOpQueues{
		maxBatchSize:     defaultMaxBatchSize,
		setsByEntryPoint: make(setsByEntryPoint),
	}
}
