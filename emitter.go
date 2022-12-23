package superwatcher

import (
	"context"
)

// Emitter receives results from Poller and emits them to Engine.
// Emitter is aware of the current service states (via StateDataGateway),
// and uses that information to determine fromBlock and toBlock for Poller.Poll.
type Emitter interface {
	// Loop is the entry point for Emitter.
	// Call it in a different loop than Engine.Loop to make both run concurrently.
	Loop(context.Context) error

	SyncsWithEngine() // Waits until engine is done processing the last batch
	Shutdown()        // Shutdown and closing Go channels

	Poller() EmitterPoller   // Poller returns the current Poller in use by Emitter
	SetPoller(EmitterPoller) // Change emitter's Poller to new one
}
