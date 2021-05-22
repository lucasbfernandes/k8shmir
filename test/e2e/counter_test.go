package e2e_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	testUtils "k8s-smr/test/e2e/utilities"
	"testing"
)

type CounterE2ETestSuite struct {
	suite.Suite
}

func (suite *CounterE2ETestSuite) SetupSuite() {
	err := testUtils.DoResetCounter()
	if err != nil {
		suite.FailNow(fmt.Sprintf("failed to setup test suite: %s", err))
	}
}

func (suite *CounterE2ETestSuite) TearDownSuite() {
}

func (suite *CounterE2ETestSuite) TearDownTest() {
	err := testUtils.DoResetCounter()
	if err != nil {
		suite.FailNow(fmt.Sprintf("failed to tear down test: %s", err))
	}
}

func (suite *CounterE2ETestSuite) TestRequestShouldPropagateToReplicasCorrectly() {
	err := testUtils.DoPostCounterRequest(testUtils.Counter1URL, testUtils.CounterIncOP, 2)
	assert.Nil(suite.T(), err, "request error should be nil")

	err = testUtils.DoPostCounterRequest(testUtils.Counter1URL, testUtils.CounterDecOP, 1)
	assert.Nil(suite.T(), err, "request error should be nil")

	err = testUtils.DoPostCounterRequest(testUtils.Counter1URL, testUtils.CounterIncOP, 3)
	assert.Nil(suite.T(), err, "request error should be nil")

	resp, err := testUtils.DoGetCounterRequest(testUtils.Counter2URL)
	assert.Nil(suite.T(), err, "request error should be nil")
	assert.Equal(suite.T(), 4, resp.Value, "replicas value should be equal")
}

func (suite *CounterE2ETestSuite) TestParallelRequestsShouldConfirmSequentialConsistency() {
	counter1Chan := make(chan struct{})
	counter2Chan := make(chan struct{})

	doRequestFunc := func(index int, url string, incVal int, decVal int) {
		if index % 2 == 0 {
			err := testUtils.DoPostCounterRequest(url, testUtils.CounterIncOP, incVal)
			assert.Nil(suite.T(), err, "request error should be nil")
		} else {
			err := testUtils.DoPostCounterRequest(url, testUtils.CounterDecOP, decVal)
			assert.Nil(suite.T(), err, "request error should be nil")
		}
	}

	go func() {
		for i := 0; i < 250; i++ {
			doRequestFunc(i, testUtils.Counter1URL, 3, 2)
		}
		close(counter1Chan)
	}()
	go func() {
		for i := 0; i < 250; i++ {
			doRequestFunc(i, testUtils.Counter2URL, 4, 7)
		}
		close(counter2Chan)
	}()

	<-counter1Chan
	<-counter2Chan

	counter1Resp, err := testUtils.DoGetCounterRequest(testUtils.Counter1URL)
	assert.Nil(suite.T(), err, "request error should be nil")

	counter2Resp, err := testUtils.DoGetCounterRequest(testUtils.Counter2URL)
	assert.Nil(suite.T(), err, "request error should be nil")

	assert.Equal(suite.T(), counter1Resp.Value, counter2Resp.Value, "replicas value should be equal")
}

func TestCounterE2ETest(t *testing.T) {
	suite.Run(t, new(CounterE2ETestSuite))
}
