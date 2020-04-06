package cache

import (
	"context"
	cache2 "github.com/sxllwx/vulcanus/pkg/feature/cache"
	"net/http"
	_ "net/http/pprof"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type cacheTestSuit struct {
	suite.Suite

	stack      cache2.Stack
	blockStack cache2.Stack

	queue      cache2.DoubleEndQueue
	blockQueue cache2.DoubleEndQueue

	set      cache2.Set
	blockSet cache2.Set
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
	s.T().Logf("got restclient from stack: %v", rest)

	s.blockStack.Close()

	s.blockQueue.Close()
	rest, err = s.blockQueue.Rest()
	assert.Nil(s.T(), err)
	s.T().Logf("got restclient from block queue: %v", rest)

	s.blockSet.Close()
	rest, err = s.blockSet.Rest()
	assert.Nil(s.T(), err)
	s.T().Logf("got restclient from block set: %v", rest)
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(cacheTestSuit))
}
