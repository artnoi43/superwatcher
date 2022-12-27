<!-- markdownlint-configure-file { "MD013": { "code_blocks": false } } -->

# Event logs filtering in emitter

The emitter's main purpose is to filter Ethereum event logs.
For practical reasons, this is done in a progressive ranged fashion,
with the forward range being determined by [`config.Config.FilterRange`](../../config/config.go).

The emitter `e` filters logs using [`e.poller.Poll(fromBlock, toBlock)`](./poller.go),
which only takes in 2 numbers denoting the start and end of the filter range.
The start of the filter range is called `fromBlock`, and the end of each loop is
called `toBlock`.

This method is called repeatedly in a for loop by [`*emitter.loopEmit`](./loop_filterlogs.go),
which is at the center of this document.

**At the end of each call to `poller.Poll`, a `lastRecordedBlock` must be saved
to some database by user code, and that value is later used in the subsequent
loop as the new basis for next `fromBlock`**.

The 3 main possible emitter states as seen in `*emitter.loopEmit` are:

<!-- markdownlint-capture-->
<!-- markdownlint-disable MD013 -->

1. First loop (first run)

   The first loop after the superwatcher service (re)starts.
   This is different than the other states in that the emitter will _go back_
   a certain number of blocks to detect chain reorgs.
   This can further be divided into 2 possibilities:

   1.1 First run on the host (got `datagateway.ErrRecordNotFound` when attempting `GetLastRecordedBlock`)

   If the emitter has never been run on the host, then the emitter will _NOT_ go back.
   The config's `StartBlock` field will instead be used as the base for `fromBlock`, and `toBlock`
   is determined using the normal case logic.

   1.2 superwatcher service restarts (first loop, but there was a legit result from `GetLastRecordedBlock`)

   If the emitter was restarted, then it will definitely need to _go back to the max_,
   to detect chain reorg that might have happened while the service was down.

2. The whole recent range was reorged (previous `fromBlock` was reorged)

   This can be detected by checking the error returned from `poller.Poll` against `superwatcher.ErrFromBlockReorged`.
   If this happens, it means that the chain is now reorging, so we need to _go back_ until all block
   are in range are canon again. If the chain is reorging for multiple loops, then we'll see that
   the `fromBlock` keeps going back, while `toBlock` stays the same.

   Unlike case 1.2 where the emitter goes back all the way back to the max limit, this case sees
   the emitter tries going back by `n x goBack` (where `n` is the number of attempts) until the
   result value reaches the configured max allowed go back.
   This case uses values computed by [`fromBlockToBlockIsReorging`](./blocknum_utils.go).

3. Normal case

   This is the base case for most of the times.
   In this case, the emitter uses `lastRecordedBlock + 1 - goBack` as the next `fromBlock`,
   and uses `lastRecordedBlock + goBack` as a new `toBlock`. This case uses values
   computed by [`fromBlockToBlockNormal`](./blocknum_utils.go)

The code block below gives some overview over how these block numbers are computed in the 3 different cases.

<!-- markdownlint-restore -->

```text
Case 1.1
# In this case, the emitter will consider this a normal case (Case 3) with
# startBlock being used as lastRecordedBlock
# N/A lastRecordedBlock, startBlock = 80, filterRange = 10, maxRetries = 5
# 71 - 90   [normalCase] -> lastRecordedBlock = 90  lookBack = 10, fwdRange = 90 - 80   = 10
# 81 - 100  [normalCase] -> lastRecordedBlock = 100 lookBack = 10, fwdRange = 100 - 90  = 10
# 91 - 110  [normalCase] -> lastRecordedBlock = 110 lookBack = 10, fwdRange = 110 - 100 = 10

Case 1.2
# Start with going back for filterRange * goBackRetries blocks if watcher was restarted
# lastRecordedBlock = 80, filterRange = 10, maxRetries = 5
# 31 - 40  [goBackFirstStart] -> lastRecordedBlock = 40 lookBack = 50, fwdRange = 0
# 31 - 50  [normalCase]       -> lastRecordedBlock = 50 lookBack = 10, fwdRange = 50 - 40   = 10
# 41 - 60  [normalCase]       -> lastRecordedBlock = 60 lookBack = 10, fwdRange = 60 - 50   = 10

Case 2
# The lookBack range will grow after each retries, but not the forward range
# lastRecordedBlock = 80, filterRange = 10
# 71  - 90   [normalCase]                -> lastRecordedBlock = 90,  lookBack = 10, fwdRange = 90 - 80   = 10
# 81  - 100  [normalCase]                -> lastRecordedBlock = 100, lookBack = 10, fwdRange = 100 - 90  = 10
# 91  - 110  # 91 reorged in this loop   -> lastRecordedBlock = 110, lookBack = 10, fwdRange = 110 - 100 = 10
# 81  - 110  # 81 reorged in this loop   -> lastRecordedBlock = 110, lookBack = 15, fwdRange = 110 - 110 = 0
# 71  - 110  # 71 reorged in this loop   -> lastRecordedBlock = 110, lookBack = 20, fwdRange = 110 - 110 = 0
# 61  - 110  # none reorged in this loop -> lastRecordedBlock = 110, lookBack = 25, fwdRange = 110 - 110 = 0
# 101 - 120  [normalCase]                -> lastRecordedBlock = 120, lookBack = 10, fwdRange = 120 - 110 = 10

Case 3
# Call fromBlockToBlock in normal cases
# lastRecordedBlock = 80, filterRange = 10
# 71  - 90   [normalCase] -> lastRecordedBlock = 90,  lookBack = 10, fwdRange = 90 - 80    = 10
# 81  - 100  [normalCase] -> lastRecordedBlock = 100, lookBack = 10, fwdRange = 100 - 90   = 10
# 91  - 110  [normalCase] -> lastRecordedBlock = 110, lookBack = 10, fwdRange = 110 - 100  = 10
# 101 - 120  [normalCase] -> lastRecordedBlock = 120, lookBack = 10, fwdRange = 120 - 110  = 10
```
