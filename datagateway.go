package superwatcher

import (
	"context"
)

// Note: Graceful shutdowns for the data gateway should be performed by service code.

// GetStateDataGateway is used by the emitter to get last recorded block.
type GetStateDataGateway interface {
	GetLastRecordedBlock(context.Context) (uint64, error)
}

// SetStateDataGateway is used by the engine to set last recorded block.
type SetStateDataGateway interface {
	SetLastRecordedBlock(context.Context, uint64) error
}

// superwatcher provides default implementation for StateDataGateway via watcherStateRedisCli
type StateDataGateway interface {
	GetStateDataGateway
	SetStateDataGateway
}

// Services can also inject functions or method into superwatcher components by wrapping the methods
// with GetStateDataGatewayFunc and SetStateDataGatewayFunc.

type FuncGetLastRecordedBlock func(context.Context) (uint64, error)
type FuncSetLastRecordedBlock func(context.Context, uint64) error

// Note: As of this writing, the emitter and engine implementations do not have fields for function types
// FuncGetLastRecordedBlock and FuncSetLastRecordedBlock.

// If you want to inject a function (not a whole struct),
// use the wrapper functions below.

// Example usage:
// ```
//  emitter := emitter.New(
//      nil,
//      nil,
//      GetStateDataGatewayFunc(someStruct.SomeFunc), // <<<<<< Use it like this
//      nil,
//      nil,
//      nil,
//      nil,
//      nil,
//  )
//
// ```

func GetStateDataGatewayFunc(f FuncGetLastRecordedBlock) GetStateDataGateway {
	return &dataGatewayWrapper{
		getFunc: f,
	}
}

func SetStateDataGatewayFunc(f FuncSetLastRecordedBlock) SetStateDataGateway {
	return &dataGatewayWrapper{
		setFunc: f,
	}
}

type dataGatewayWrapper struct {
	getFunc FuncGetLastRecordedBlock
	setFunc FuncSetLastRecordedBlock
}

func (w *dataGatewayWrapper) GetLastRecordedBlock(ctx context.Context) (uint64, error) {
	return w.getFunc(ctx)
}

func (w *dataGatewayWrapper) SetLastRecordedBlock(ctx context.Context, blockNumber uint64) error {
	return w.setFunc(ctx, blockNumber)
}
