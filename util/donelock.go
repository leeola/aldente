package util

import "sync"

type DoneLock struct {
	*sync.Mutex

	Done bool
}

func NewDoneLock() *DoneLock {
	return &DoneLock{
		Mutex: &sync.Mutex{},
	}
}
