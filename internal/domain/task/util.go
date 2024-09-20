package task

import (
	"math/rand"
	"time"
)

//go:generate mockery --name Random --output=mocks/
type Random interface {
	Int63n(value int64) int64
}

type gorand struct {
	rand *rand.Rand
}

func NewRandom() Random {
	return NewRandomWithSeed(time.Now().UnixNano())
}

func NewRandomWithSeed(seed int64) Random {
	return &gorand{rand: rand.New(rand.NewSource(seed))} //nolint:gosec // we don't need a secure generator here
}

func (r *gorand) Int63n(value int64) int64 {
	return r.rand.Int63n(value)
}
