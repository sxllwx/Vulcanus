package cache

import (
	"github.com/stretchr/testify/assert"
)

func (s *cacheTestSuit) TestStack() {

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

func (s *cacheTestSuit) TestBlockStack() {

	go func() {

		for {
			s.T().Logf("try to get element")
			e, err := s.blockStack.Pop()
			s.T().Logf("block got element: %v, err -> %v", e, err)
		}
	}()

	assert.Nil(s.T(), s.blockStack.Push(0))
	assert.Nil(s.T(), s.blockStack.Push(1))
	assert.Nil(s.T(), s.blockStack.Push(2))
	assert.Nil(s.T(), s.blockStack.Push(3))
	assert.Nil(s.T(), s.blockStack.Push(4))
	assert.Nil(s.T(), s.blockStack.Push(5))
	assert.Nil(s.T(), s.blockStack.Push(6))
	assert.Nil(s.T(), s.blockStack.Push(7))
	assert.Nil(s.T(), s.blockStack.Push(8))
	assert.Nil(s.T(), s.blockStack.Push(9))
	assert.Nil(s.T(), s.blockStack.Push(10))
}
