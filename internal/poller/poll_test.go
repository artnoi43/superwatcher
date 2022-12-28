package poller

import (
	"testing"

	"github.com/artnoi43/superwatcher"
)

func TestLastGoodBlock(t *testing.T) {
	tests := map[*superwatcher.FilterResult]uint64{
		{
			FromBlock: 101,
			ToBlock:   120,
			GoodBlocks: []*superwatcher.BlockInfo{
				{Number: 102},
				{Number: 104},
				{Number: 105},
			},
			ReorgedBlocks: []*superwatcher.BlockInfo{
				{Number: 105},
			},
		}: 105,
		{
			FromBlock: 101,
			ToBlock:   120,
			GoodBlocks: []*superwatcher.BlockInfo{
				{Number: 102},
				{Number: 104},
				{Number: 105},
			},
			ReorgedBlocks: []*superwatcher.BlockInfo{
				{Number: 107},
			},
		}: 105,
		{
			FromBlock: 101,
			ToBlock:   120,
			GoodBlocks: []*superwatcher.BlockInfo{
				{Number: 107},
				{Number: 108},
			},
			ReorgedBlocks: []*superwatcher.BlockInfo{
				{Number: 105},
			},
		}: 104,
		{
			FromBlock: 101,
			ToBlock:   120,
			GoodBlocks: []*superwatcher.BlockInfo{
				{Number: 103},
				{Number: 104},
				{Number: 105},
			},
			ReorgedBlocks: []*superwatcher.BlockInfo{
				{Number: 103},
			},
		}: 102,
		{
			FromBlock: 101,
			ToBlock:   120,
			GoodBlocks: []*superwatcher.BlockInfo{
				{Number: 103},
				{Number: 104},
				{Number: 105},
			},
			ReorgedBlocks: []*superwatcher.BlockInfo{
				{Number: 104},
			},
		}: 103,
	}

	for filterResult, expected := range tests {
		actual := superwatcher.LastGoodBlock(filterResult)
		if actual != expected {
			t.Errorf("expecting lastGoodBlock %d, got %d", expected, actual)
		}
	}
}
