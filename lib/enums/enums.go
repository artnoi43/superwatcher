package enums

type ChainType string

const (
	ChainEthereum ChainType = "ethereum"
	ChainBSC      ChainType = "bsc"
)

func (c ChainType) IsValid() bool {
	switch c {
	case ChainEthereum, ChainBSC:
		return true
	}

	return false
}
