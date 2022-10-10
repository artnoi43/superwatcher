package reorg

/*
	[compareWithTracker]: Populate processLogs from fresh event logs,
	and populate freshHashesByBlockNumber with fresh block hashes.

	[ProcessReorgs]: Check if block hash saved in tracker matches the fresh block hash.
	If they are different, old logs from w.tracker will be tagged as Removed and
	PREPENDED in processLogs[blockNumber]

	Let's say we have these logs in the tracker:

	{block:68, hash:"0x68"}, {block: 69, hash:"0x69"}, {block:70, hash:"0x70"}

	And then we have these fresh logs:

	{block:68, hash:"0x68"}, {block: 69, hash:"0x112"}, {block:70, hash:"0x70"}

	The result processLogs will look like this map:
	{
		68: [{block:68, hash:"0x68"}]
		69: [{block: 69, hash:"0x69", removed: true}, {block: 69, hash:"0x112"}]
		70: [{block:70, hash:"0x70"}]
	}
*/

func compareWithTracker() { panic("not implemented") }

func ProcessReorgs() { panic("not implemented") }
