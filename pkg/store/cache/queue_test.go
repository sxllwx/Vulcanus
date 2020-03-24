package cache

import (
	"github.com/stretchr/testify/assert"
)

func (s *cacheTestSuit) TestQueue() {

	e, err := s.queue.DeQueue()
	assert.NotNil(s.T(), err)
	s.T().Logf("the queue response pop %v", err)

	assert.Nil(s.T(), s.queue.EnQueue(0))
	assert.Nil(s.T(), s.queue.EnQueue(1))
	assert.Nil(s.T(), s.queue.EnQueue(2))
	assert.Nil(s.T(), s.queue.EnQueue(3))
	assert.Nil(s.T(), s.queue.EnQueue(4))
	assert.Nil(s.T(), s.queue.EnQueue(5))

	e, err = s.queue.DeQueue()
	assert.Nil(s.T(), err)
	s.T().Logf("got element: %v", e)

	e, err = s.queue.DeQueue()
	assert.Nil(s.T(), err)
	s.T().Logf("got element: %v", e)

	e, err = s.queue.DeQueueFromTail()
	assert.Nil(s.T(), err)
	s.T().Logf("got element from tail: %v", e)

	err = s.queue.EnQueueToFront(6)
	assert.Nil(s.T(), err)

	e, err = s.queue.DeQueue()
	assert.Nil(s.T(), err)
	s.T().Logf("got element: %v", e)
}

func (s *cacheTestSuit) TestBlockQueue() {

	go func() {
		for {
			e, err := s.blockQueue.DeQueue()
			assert.Nil(s.T(), err)
			s.T().Logf("got elment from block queue %v", e)

			if e.(int) == 4 {
				return
			}
		}
	}()

	assert.Nil(s.T(), s.blockQueue.EnQueue(0))
	assert.Nil(s.T(), s.blockQueue.EnQueue(1))
	assert.Nil(s.T(), s.blockQueue.EnQueue(2))
	assert.Nil(s.T(), s.blockQueue.EnQueue(3))
	assert.Nil(s.T(), s.blockQueue.EnQueue(4))
	assert.Nil(s.T(), s.blockQueue.EnQueue(5))
	assert.Nil(s.T(), s.blockQueue.EnQueue(6))
	assert.Nil(s.T(), s.blockQueue.EnQueue(7))
	assert.Nil(s.T(), s.blockQueue.EnQueue(8))
}
