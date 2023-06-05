package mock

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"

	"github.com/soyart/superwatcher"
)

func TestFakeRedis(t *testing.T) {
	t.Run("testFakeRedisMem", func(t *testing.T) {
		testFakeRedisMem(t)
	})
	t.Run("testFakeRedisFile", func(t *testing.T) {
		testFakeRedisFile(t)
	})
}

func testFakeRedisMem(t *testing.T) {
	x := uint64(69)
	y := uint64(100)

	f := NewDataGatewayMem(x, true)
	lastRec, err := f.GetLastRecordedBlock(nil)
	if err != nil {
		t.Error("error in fakeRedis.GetLastRecordedBlock", err.Error())
	}
	if lastRec != x {
		t.Errorf("unexpected result from fakeRedis.GetLastRecordedBlock - expecting %d, got %d", x, lastRec)
	}
	if err := setAndGet(f, y); err != nil {
		t.Error(err.Error())
	}

	f = NewDataGatewayMem(x, false) // Never run before - Get before Set should fail
	lastRec, err = f.GetLastRecordedBlock(nil)
	if err != nil {
		if !errors.Is(err, superwatcher.ErrRecordNotFound) {
			t.Error("error in fakeRedis.GetLastRecordedBlock not ErrRecordNotFound", err.Error())
		}
	}
	if lastRec != 0 {
		t.Errorf("unexpected result from fakeRedis.GetLastRecordedBlock - expecting %d, got %d", 0, lastRec)
	}
	if err := setAndGet(f, x); err != nil {
		t.Error(err.Error())
	}
	if err := setAndGet(f, y); err != nil {
		t.Error(err.Error())
	}
}

func testFakeRedisFile(t *testing.T) {
	x := uint64(69)
	y := uint64(100)
	filename := "tmp/fakeredis.db" // Will be ./tmp/fakeredis.db RELATIVE TO THIS TEST FILE

	f := NewDataGatewayFile(filename, x, true)
	lastRec, err := f.GetLastRecordedBlock(nil)
	if err != nil {
		t.Error("error in fakeRedis.GetLastRecordedBlock", err.Error())
	}
	if lastRec != x {
		t.Errorf("unexpected result from fakeRedis.GetLastRecordedBlock - expecting %d, got %d", x, lastRec)
	}
	if err := setAndGet(f, y); err != nil {
		t.Error(err.Error())
	}

	f = NewDataGatewayFile(filename, x, false) // Never run before - Get before Set should fail
	lastRec, err = f.GetLastRecordedBlock(nil)
	if err != nil {
		if !errors.Is(err, superwatcher.ErrRecordNotFound) {
			t.Error("error in fakeRedis.GetLastRecordedBlock not ErrRecordNotFound", err.Error())
		}
	}
	if lastRec != 0 {
		t.Errorf("unexpected result from fakeRedis.GetLastRecordedBlock - expecting %d, got %d", 0, lastRec)
	}
	if err := setAndGet(f, x); err != nil {
		t.Error(err.Error())
	}
	if err := setAndGet(f, y); err != nil {
		t.Error(err.Error())
	}
}

func setAndGet(f superwatcher.StateDataGateway, x uint64) error {
	if err := f.SetLastRecordedBlock(nil, x); err != nil {
		return errors.Wrap(err, "error in fakeRedis.GetLastRecordedBlock")
	}

	lastRec, err := f.GetLastRecordedBlock(nil)
	if err != nil {
		return errors.Wrap(err, "error in fakeRedis.GetLastRecordedBlock")
	}

	if lastRec != x {
		return fmt.Errorf("unexpected result from fakeRedis.GetLastRecordedBlock - expecting %d, got %d", x, lastRec)
	}

	return nil
}
