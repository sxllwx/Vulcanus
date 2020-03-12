package cache

import "github.com/stretchr/testify/assert"

func (s *cacheTestSuit) TestSet() {

	e, err := s.set.Get()
	assert.NotNil(s.T(), err)
	s.T().Logf("the set response pop %v", err)

	assert.Nil(s.T(), s.set.Put(0))
	assert.Nil(s.T(), s.set.Put(1))
	assert.Nil(s.T(), s.set.Put(2))
	assert.Nil(s.T(), s.set.Put(3))
	assert.Nil(s.T(), s.set.Put(4))
	assert.Nil(s.T(), s.set.Put(5))

	e, err = s.set.Get()
	assert.Nil(s.T(), err)
	s.T().Logf("got element: %v", e)

	e, err = s.set.Get()
	assert.Nil(s.T(), err)
	s.T().Logf("got element: %v", e)

	e, err = s.set.Get()
	assert.Nil(s.T(), err)
	s.T().Logf("got element: %v", e)
}

func (s *cacheTestSuit) TestBlockSet() {

	go func() {

		i := 0
		for {
			e, err := s.blockSet.Get()
			assert.Nil(s.T(), err)
			s.T().Logf("got elment from block set %v", e)

			i++

			if i == 4 {
				return
			}
		}
	}()

	go s.blockSet.Put(0)
	go s.blockSet.Put(1)
	go s.blockSet.Put(2)
	go s.blockSet.Put(3)
	go s.blockSet.Put(4)
	go s.blockSet.Put(5)
	go s.blockSet.Put(6)
	go s.blockSet.Put(7)
	go s.blockSet.Put(8)
}
