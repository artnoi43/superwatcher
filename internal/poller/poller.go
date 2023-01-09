package poller

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

// poller implements superwatcher.WatcherPoller.
// It filters event logs and detecting chain reorg, i.e. it produces superwatcher.PollerResult for the emitter.
// poller behavior can be changed on-the-fly using methods defined in this file.
type poller struct {
	sync.RWMutex

	addresses []common.Address
	topics    [][]common.Hash

	lastRecordedBlock uint64 // For clearing tracker if SetDoReorg(false) is called
	filterRange       uint64
	client            superwatcher.EthClient
	doReorg           bool
	doHeader          bool
	policy            superwatcher.Policy

	tracker  *blockTracker
	debugger *debugger.Debugger
}

func New(
	addresses []common.Address,
	topics [][]common.Hash,
	doReorg bool,
	doHeader bool,
	filterRange uint64,
	client superwatcher.EthClient,
	logLevel uint8,
	policy superwatcher.Policy,
) superwatcher.EmitterPoller {
	var tracker *blockTracker
	if doReorg {
		tracker = newTracker("poller", logLevel)
	}

	return &poller{
		addresses:   addresses,
		topics:      topics,
		filterRange: filterRange,
		client:      client,
		doReorg:     doReorg,
		doHeader:    doHeader,
		tracker:     tracker,
		debugger:    debugger.NewDebugger("poller", logLevel),
		policy:      policy,
	}
}

// SetDoReorg changes poller behavior regarding processing chain reorg.
// If |doReorg| is false, setDoReorg clears all tracker data and deletes the tracker with nil value,
// and changes p.doReorg policy to false.
// If |doReorg| is true, SetDoReorg creates new tracker and update p.doReorg policy if it was false.
func (p *poller) SetDoReorg(doReorg bool) {
	p.Lock()
	defer p.Unlock()

	switch p.doReorg {
	case true:
		if doReorg {
			return
		}

		p.debugger.Debug(1, "SetDoReorg(false) called - clearing tracker and deleting tracker")
		p.tracker.clearUntil(p.lastRecordedBlock)
		p.tracker = nil

	case false:
		// p.doReorg = false, doReorg = false
		if !doReorg {
			p.debugger.Debug(1, "SetDoReorg(false) called, but p.doReorg is already false - returning")
			return
		}

		// p.doReorg = false, doReorg = true
		if p.tracker == nil {
			p.debugger.Debug(1, "SetDoReorg(true) called - creating new tracker")
			p.tracker = newTracker("poller", p.debugger.Level)
		} else {
			p.debugger.Debug(1, "SetDoReorg(true) called but tracker is not nil - reusing tracker")
		}
	}

	p.doReorg = doReorg
}

func (p *poller) DoReorg() bool {
	p.RLock()
	defer p.RUnlock()

	return p.doReorg
}

func (p *poller) SetDoHeader(doHeader bool) {
	p.Lock()
	defer p.Unlock()

	if doHeader && p.policy >= superwatcher.PolicyExpensive {
		p.debugger.Debug(1, "SetDoHeader called, but PolicyExpensive is set, ignoring doHeader value")
	}

	p.doHeader = doHeader
}

func (p *poller) DoHeader() bool {
	p.RLock()
	defer p.RUnlock()

	return p.doHeader
}

func (p *poller) Addresses() []common.Address {
	p.RLock()
	defer p.RUnlock()

	return p.addresses
}

func (p *poller) Topics() [][]common.Hash {
	p.RLock()
	defer p.RUnlock()

	return p.topics
}

func (p *poller) AddAddresses(addresses ...common.Address) {
	p.Lock()
	defer p.Unlock()

	if len(p.addresses) == 0 {
		p.addresses = addresses
		return
	}

	p.addresses = append(p.addresses, addresses...)
}

func (p *poller) AddTopics(topics ...[]common.Hash) {
	p.Lock()
	defer p.Unlock()

	if len(p.topics) == 0 {
		p.topics = topics
		return
	}

	p.topics = append(p.topics, topics...)
}

func (p *poller) SetAddresses(addresses []common.Address) {
	p.Lock()
	defer p.Unlock()

	p.addresses = addresses
}

func (p *poller) SetTopics(topics [][]common.Hash) {
	p.Lock()
	defer p.Unlock()

	p.topics = topics
}

func (p *poller) SetPolicy(level superwatcher.Policy) error {
	p.Lock()
	defer p.Unlock()

	// Remove all blocks from tracker with 0 logs if not PolicyExpensive,
	// because if these blocks are left in tracker, poller with level < PolicyExpensive
	// will see that these blocks are ones with missing logs and stamped as reorged = true.
	if p.policy < level {
		return errors.Wrapf(
			errDowngradeLevel,
			"cannot downgrade from %s (%d) to %s (%d)",
			p.policy.String(), p.policy, level.String(), level,
		)
	}

	p.policy = level
	return nil
}

func (p *poller) Policy() superwatcher.Policy {
	p.RLock()
	defer p.RUnlock()

	return p.policy
}
