package superwatcher

import "context"

// Engine receives PollResult emitted from Emitter
// and executes business service logic on PollResult with ServiceEngine.
type Engine interface {
	// Loop is the entry point for Engine.
	// Call it in a different Goroutine than Emitter.Loop to make both run concurrently.
	Loop(context.Context) error
}
