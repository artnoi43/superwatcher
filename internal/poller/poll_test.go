package poller

import (
	"testing"

	"github.com/artnoi43/superwatcher"
)

func TestLastGoodBlock(t *testing.T) {
	tests := map[*superwatcher.PollerResult]uint64{
		{
			FromBlock: 101,
			ToBlock:   120,
			GoodBlocks: []*superwatcher.Block{
				{Number: 102},
				{Number: 104},
				{Number: 105},
			},
			ReorgedBlocks: []*superwatcher.Block{
				{Number: 105},
			},
		}: 105,
		{
			FromBlock: 101,
			ToBlock:   120,
			GoodBlocks: []*superwatcher.Block{
				{Number: 102},
				{Number: 104},
				{Number: 105},
			},
			ReorgedBlocks: []*superwatcher.Block{
				{Number: 107},
			},
		}: 105,
		{
			FromBlock: 101,
			ToBlock:   120,
			GoodBlocks: []*superwatcher.Block{
				{Number: 107},
				{Number: 108},
			},
			ReorgedBlocks: []*superwatcher.Block{
				{Number: 105},
			},
		}: 104,
		{
			FromBlock: 101,
			ToBlock:   120,
			GoodBlocks: []*superwatcher.Block{
				{Number: 103},
				{Number: 104},
				{Number: 105},
			},
			ReorgedBlocks: []*superwatcher.Block{
				{Number: 103},
			},
		}: 102,
		{
			FromBlock: 101,
			ToBlock:   120,
			GoodBlocks: []*superwatcher.Block{
				{Number: 103},
				{Number: 104},
				{Number: 105},
			},
			ReorgedBlocks: []*superwatcher.Block{
				{Number: 104},
			},
		}: 103,
	}

	for PollerResult, expected := range tests {
		actual := superwatcher.LastGoodBlock(PollerResult)
		if actual != expected {
			t.Errorf("expecting lastGoodBlock %d, got %d", expected, actual)
		}
	}
}
