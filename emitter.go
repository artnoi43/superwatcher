package superwatcher

import (
	"context"
)

// Emitter receives results from Poller and emits them to Engine.
// Emitter is aware of the current service states (via StateDataGateway),
// and uses that information to determine fromBlock and toBlock for Poller.Poll.
type Emitter interface {
	// Loop is the entry point for Emitter.
	// Users will call Loop in a different loop than Engine.Loop
	// to make both components run concurrently.
	Loop(context.Context) error
	// Waits until engine is done processing the last batch
	SyncsWithEngine()
	// Shutdown and closing Go channels
	Shutdown()
	// Poller returns the current Poller in use by Emitter
	Poller() EmitterPoller
	// Change emitter's Poller to new one
	SetPoller(EmitterPoller)
}
