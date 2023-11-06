package altmempools

import (
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/puzpuzpuz/xsync/v3"
)

// Directory maintains a collection of alternative mempool configurations. It allows a consumer to check if a
// known alternative mempool exists that will allow specific exceptions that the canonical mempool cannot
// accept.
type Directory struct {
	invalidStorageAccess *xsync.MapOf[string, []string]
}

type Config struct {
	Id   string
	Data map[string]any
}

func invalidStorageAccessID(entity string, contract string, slot string) string {
	return entity + contract + slot
}

func fetchMempoolConfig(url string) (map[string]any, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// New accepts an array of alternative mempool configs and returns a Directory.
func New(chain *big.Int, altMempools []*Config) (*Directory, error) {
	dir := &Directory{
		invalidStorageAccess: xsync.NewMapOf[string, []string](),
	}
	for _, alt := range altMempools {
		if err := Schema.Validate(alt.Data); err != nil {
			return nil, err
		}

		skip := true
		for _, item := range alt.Data["chainIds"].([]any) {
			allowed, err := hexutil.DecodeBig(item.(string))
			if err != nil {
				return nil, err
			}

			if chain.Cmp(allowed) == 0 {
				skip = false
			}
		}
		if skip {
			continue
		}

		for _, item := range alt.Data["allowlist"].([]any) {
			config := item.(map[string]any)
			switch config["rule"].(string) {
			case "invalidStorageAccess":
				{
					isaId := invalidStorageAccessID(
						config["entity"].(string),
						config["contract"].(string),
						config["slot"].(string),
					)
					curr, _ := dir.invalidStorageAccess.Load(isaId)
					dir.invalidStorageAccess.Store(isaId, append(curr, alt.Id))
				}
			}
		}
	}

	return dir, nil
}

// NewFromIPFS will pull alternative mempool configs from IPFS and returns a Directory. The mempool id is
// equal to an IPFS CID.
func NewFromIPFS(chain *big.Int, ipfsGateway string, ids []string) (*Directory, error) {
	var alts []*Config
	for _, id := range ids {
		data, err := fetchMempoolConfig(ipfsGateway + "/" + id)
		if err != nil {
			return nil, err
		}
		alts = append(alts, &Config{id, data})
	}

	return New(chain, alts)
}

// HasInvalidStorageAccessException will attempt to find all mempools ids that will accept the given invalid
// storage access pattern and return it. If none is found, an empty array will be returned.
func (d *Directory) HasInvalidStorageAccessException(entity string, contract string, slot string) []string {
	ids, _ := d.invalidStorageAccess.Load(invalidStorageAccessID(entity, contract, slot))
	return ids
}
