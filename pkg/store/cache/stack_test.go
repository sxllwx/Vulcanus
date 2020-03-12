package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/sxllwx/vulcanus/pkg/store"
)

type StackTestSuit struct {
	suite.Suite

	stack      store.Stack
	blockStack store.Stack
}

func (s *StackTestSuit) SetupSuite() {

	parent := context.Background()
	s.stack = NewStack(parent)
	s.blockStack = NewBlockStack(parent, func(i int) bool {
		return i > 5
	})

}

func (s *StackTestSuit) TestStack() {

	e, err := s.stack.Pop()
	assert.NotNil(s.T(), err)
	s.T().Logf("the stack response pop %v", err)

	assert.Nil(s.T(), s.stack.Push(0))
	assert.Nil(s.T(), s.stack.Push(1))
	assert.Nil(s.T(), s.stack.Push(2))
	assert.Nil(s.T(), s.stack.Push(3))
	assert.Nil(s.T(), s.stack.Push(4))
	assert.Nil(s.T(), s.stack.Push(5))

	e, err = s.stack.Pop()
	assert.Nil(s.T(), err)
	s.T().Logf("got element: %v", e)

	e, err = s.stack.Pop()
	assert.Nil(s.T(), err)
	s.T().Logf("got element: %v", e)

	e, err = s.stack.Pop()
	assert.Nil(s.T(), err)
	s.T().Logf("got element: %v", e)

	e, err = s.stack.Pop()
	assert.Nil(s.T(), err)
	s.T().Logf("got element: %v", e)
}

func (s *StackTestSuit) TestBlockStack() {

	go func() {
		s.T().Logf("try to get element")
		e, err := s.blockStack.Pop()
		assert.Nil(s.T(), err)
		s.T().Logf("block got element: %v", e)
	}()

	time.Sleep(time.Second)

	assert.Nil(s.T(), s.blockStack.Push(0))
	assert.Nil(s.T(), s.blockStack.Push(1))
	assert.Nil(s.T(), s.blockStack.Push(2))

}

func (s *StackTestSuit) TearDownSuite() {

	s.stack.Close()

	rest, err := s.stack.Rest()
	assert.Nil(s.T(), err)
	s.T().Logf("got rest: %v", rest)

	s.blockStack.Close()
}

func TestStack(t *testing.T) {
	suite.Run(t, new(StackTestSuit))
}
