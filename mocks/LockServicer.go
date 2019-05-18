package mocks

import (
	"github.com/stretchr/testify/mock"
	"sync"
	"time"
)

type LockServicer struct {
	mock.Mock
}

// NewLock provides a mock function with given fields: _a0, _a1, _a2
func (_m *LockServicer) NewLock(_a0 string, _a1 int64, _a2 time.Duration) sync.Locker {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 sync.Locker
	if rf, ok := ret.Get(0).(func(string, int64, time.Duration) sync.Locker); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(sync.Locker)
	}

	return r0
}
