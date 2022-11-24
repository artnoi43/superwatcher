package ensengine

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
)

func TestArtifacts(t *testing.T) {
	l := int64(10)
	txHashes := make([]common.Hash, l)
	blockHashes := make([]common.Hash, l)
	names := make([]string, l)

	for i := int64(1); i < l; i++ {
		txHashes[i] = common.BigToHash(big.NewInt(i))
		blockHashes[i] = common.BigToHash(big.NewInt(i + 100))
		names[i] = fmt.Sprintf("ens%d", i)
	}

	artifacts := make([]superwatcher.Artifact, l)
	for i := int64(1); i < l; i++ {
		ensArtifact := ENSArtifact{
			ENS: entity.ENS{
				Name:             names[i],
				TxHash:           gslutils.StringerToLowerString(txHashes[i]),
				BlockHashCreated: gslutils.StringerToLowerString(blockHashes[i]),
			},
		}

		artifacts[i] = ensArtifact
	}

	var ensArtifacts = make([]ENSArtifact, l)
	for i := int64(1); i < l; i++ {
		txHash := txHashes[i]
		log := &types.Log{TxHash: txHash}
		out := spwArtifactsByTxHash(log, artifacts)
		if out == nil {
			t.Errorf("nil from spwArtifactsByTxHash for txHash %s", txHash.String())
		}

		expected := artifacts[i].(ENSArtifact)
		if outName, expectedName := out.ENS.Name, expected.ENS.Name; outName != expectedName {
			t.Errorf("invalid name, expecting %s got %s", expectedName, outName)
		}

		ensArtifacts[i] = artifacts[i].(ENSArtifact)

		// Will panic if prev not found
		prev := getPrevENSArtifactFromLogTxHash(log, ensArtifacts)
		t.Log(prev.ENS.Name)
	}
}
