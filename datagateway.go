package superwatcher

import (
	"context"

	"github.com/pkg/errors"
)

// ErrRecordNotFound is checked for in emitter.loopEmit.
// If the error is ErrRecordNotFound, the emitter assumes the service
// has never run on this host (hence no data in the database), and will not attempt to go back.
var ErrRecordNotFound = errors.New("record not found")

func WrapErrRecordNotFound(err error, keyNotFound string) error {
	err = errors.Wrap(ErrRecordNotFound, err.Error())
	return errors.Wrapf(err, "key %s not found", keyNotFound)
}

type (
	// GetStateDataGateway is used by the emitter to get last recorded block.
	GetStateDataGateway interface {
		GetLastRecordedBlock(context.Context) (uint64, error)
	}

	// SetStateDataGateway is used by the engine to set last recorded block.
	SetStateDataGateway interface {
		SetLastRecordedBlock(context.Context, uint64) error
	}

	// StateDataGateway is an interface that could both set and get lastRecordedBlock for superwatcher.
	// Note: Graceful shutdowns for the StateDataGateway should be performed by service code.
	StateDataGateway interface {
		GetStateDataGateway
		SetStateDataGateway
	}

	FuncGetLastRecordedBlock func(context.Context) (uint64, error)
	FuncSetLastRecordedBlock func(context.Context, uint64) error

	// Note: As of this writing, the emitter and engine implementations do not have fields for
	// function types FuncGetLastRecordedBlock and FuncSetLastRecordedBlock.
	// If you want to inject a function (not a whole struct), use the wrapper functions below
	// to wrap your functions or methods with dataGatewayWrapper.
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

	dataGatewayWrapper struct {
		getFunc FuncGetLastRecordedBlock
		setFunc FuncSetLastRecordedBlock
	}
)

func GetStateDataGatewayFunc(f FuncGetLastRecordedBlock) GetStateDataGateway {
	return &dataGatewayWrapper{getFunc: f}
}

func SetStateDataGatewayFunc(f FuncSetLastRecordedBlock) SetStateDataGateway {
	return &dataGatewayWrapper{setFunc: f}
}

func (w *dataGatewayWrapper) GetLastRecordedBlock(ctx context.Context) (uint64, error) {
	return w.getFunc(ctx)
}

func (w *dataGatewayWrapper) SetLastRecordedBlock(ctx context.Context, blockNumber uint64) error {
	return w.setFunc(ctx, blockNumber)
}
