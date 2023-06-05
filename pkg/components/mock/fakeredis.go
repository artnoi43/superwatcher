package mock

import "github.com/soyart/superwatcher"

// NewDataGatewayMem returns a `superwatcher.StateDataGateway`.
// If |ok| is false, `GetLastRecordedBlock` returns `ErrRecordNotFound`
// until the first call to `SetLastRecordedBlock` is made.
// If |ok| is true, `GetLastRecordedBlock` will keep returning |lastRecordedBlock|
// until the value is changed with `SetLastRecordedBlock`.
func NewDataGatewayMem(lastRecordedBlock uint64, ok bool) superwatcher.StateDataGateway {
	return &fakeRedisMem{
		lastRecordedBlock: lastRecordedBlock,
		ok:                ok,
	}
}

// NewDataGatewayFile returns a `superwatcher.StateDataGateway` with persistent file storage.
// If |ok| is false, `GetLastRecordedBlock` returns `ErrRecordNotFound`
// until the first call to `SetLastRecordedBlock` is made.
// If |ok| is true, `GetLastRecordedBlock` will keep returning |lastRecordedBlock|
// until the value is changed with `SetLastRecordedBlock`.
func NewDataGatewayFile(filename string, lastRecordedBlock uint64, ok bool) superwatcher.StateDataGateway {
	// Write lastRecordedBlock before first call to `GetLastRecordedBlock`
	if ok {
		if err := writeLastRecordedBlockToFile(filename, lastRecordedBlock); err != nil {
			panic("failed to write file to init fakeRedisFile: " + err.Error())
		}
	}
	return &fakeRedisFile{
		filename: filename,
		ok:       ok,
	}
}
