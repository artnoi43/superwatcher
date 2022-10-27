package datagateway

import "context"

type StateDataGateway interface {
	GetLastRecordedBlock(context.Context) (uint64, error)
	SetLastRecordedBlock(context.Context, uint64) error

	Shutdown()
}
