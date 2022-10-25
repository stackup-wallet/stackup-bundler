package redispool

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis/v9"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

const maxRetries = 10
const prioritySortedSetKey = "priority"
const expirySortedSetKey = "expiry"
const ttlSeconds = 300

func newRedisClient(connectionString string) *redis.Client {
	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}

func addFunc(rdb *redis.Client) mempool.Add {
	return func(sender string, op *userop.UserOperation, epAddr common.Address) (bool, error) {
		memOp, err := op.MarshalJSON()
		if err != nil {
			return false, err
		}
		txn, err := json.Marshal(&[]string{string(memOp), epAddr.String()})
		if err != nil {
			return false, err
		}

		ctx := context.Background()
		txf := func(tx *redis.Tx) error {
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, sender, txn, 0)
				pipe.ZAdd(ctx, prioritySortedSetKey, redis.Z{
					Member: sender,
					Score:  float64(op.MaxPriorityFeePerGas.Int64()),
				})
				pipe.ZAdd(ctx, expirySortedSetKey, redis.Z{
					Member: sender,
					Score:  float64(time.Now().Unix() + ttlSeconds),
				})

				return nil
			})
			return err
		}

		for i := 0; i < maxRetries; i++ {
			err := rdb.Watch(ctx, txf, sender)
			if err == nil {
				break
			}
			if i != maxRetries-1 && err == redis.TxFailedErr {
				continue
			}
			return false, err
		}
		return true, nil
	}
}

func getFunc(rdb *redis.Client) mempool.Get {
	return func(sender string) (*mempool.PendingTransaction, error) {
		val, err := rdb.Get(context.Background(), sender).Result()
		if err != nil {
			if err == redis.Nil {
				return nil, nil
			}

			return nil, err
		}

		memTxn := []string{}
		json.Unmarshal([]byte(val), &memTxn)

		op := make(map[string]interface{})
		json.Unmarshal([]byte(memTxn[0]), &op)
		ep := common.HexToAddress(memTxn[1])

		userOp, err := userop.New(op)
		if err != nil {
			return nil, err
		}
		return &mempool.PendingTransaction{
			UserOp:     userOp,
			EntryPoint: ep,
		}, nil
	}
}

func NewClientInterface(connectionString string) (*mempool.ClientInterface, error) {
	rdb := newRedisClient(connectionString)

	return &mempool.ClientInterface{
		Add: addFunc(rdb),
		Get: getFunc(rdb),
	}, nil
}
