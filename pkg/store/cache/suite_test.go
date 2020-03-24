package cache

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/sxllwx/vulcanus/pkg/store"
)

type cacheTestSuit struct {
	suite.Suite

	stack      store.Stack
	blockStack store.Stack

	queue      store.DoubleEndQueue
	blockQueue store.DoubleEndQueue

	set      store.Set
	blockSet store.Set
}

func (s *cacheTestSuit) SetupSuite() {

	go http.ListenAndServe(":19191", nil)

	parent := context.Background()

	s.stack = NewStack(parent)
	s.blockStack = NewBlockStack(parent, func(i int) bool {
		return i > 5
	})

	s.queue = NewDeQueue(parent)
	s.blockQueue = NewBlockDeQueue(parent, func(i int) bool {
		return i > 4
	})

	s.set = NewSet(parent)
	s.blockSet = NewBlockSet(parent, func(i int) bool {
		return i > 4
	})
}

func (s *cacheTestSuit) TearDownSuite() {

	s.stack.Close()
	rest, err := s.stack.Rest()
	assert.Nil(s.T(), err)
	s.T().Logf("got rest from stack: %v", rest)

	s.blockStack.Close()

	s.blockQueue.Close()
	rest, err = s.blockQueue.Rest()
	assert.Nil(s.T(), err)
	s.T().Logf("got rest from block queue: %v", rest)

	s.blockSet.Close()
	rest, err = s.blockSet.Rest()
	assert.Nil(s.T(), err)
	s.T().Logf("got rest from block set: %v", rest)
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(cacheTestSuit))
}
