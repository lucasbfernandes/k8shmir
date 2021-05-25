package e2e_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"k8s-smr/test/e2e/utilities/counter"
	"testing"
)

const (
	restartCounter2AppScriptPath = "utilities/counter/restart-counter2-app.sh"

	restartCounter2ProxyScriptPath = "utilities/counter/restart-counter2-proxy.sh"
)

type CounterE2ETestSuite struct {
	suite.Suite
}

func (suite *CounterE2ETestSuite) SetupSuite() {
	err := counter.DoResetCounter()
	if err != nil {
		suite.FailNow(fmt.Sprintf("failed to setup test suite: %s", err))
	}
}

func (suite *CounterE2ETestSuite) TearDownSuite() {
}

func (suite *CounterE2ETestSuite) TearDownTest() {
	err := counter.DoResetCounter()
	if err != nil {
		suite.FailNow(fmt.Sprintf("failed to tear down test: %s", err))
	}
}

func (suite *CounterE2ETestSuite) TestRequestShouldPropagateToReplicasCorrectly() {
	err := counter.DoPostCounterRequest(counter.URL1, counter.IncOP, 2)
	assert.Nil(suite.T(), err, "request error should be nil")

	err = counter.DoPostCounterRequest(counter.URL1, counter.DecOP, 1)
	assert.Nil(suite.T(), err, "request error should be nil")

	err = counter.DoPostCounterRequest(counter.URL1, counter.IncOP, 3)
	assert.Nil(suite.T(), err, "request error should be nil")

	resp, err := counter.DoGetCounterRequest(counter.URL2)
	assert.Nil(suite.T(), err, "request error should be nil")
	assert.Equal(suite.T(), 4, resp.Value, "replicas value should be equal")
}

func (suite *CounterE2ETestSuite) TestParallelRequestsShouldConfirmSequentialConsistency() {
	counter1Chan := make(chan struct{})
	counter2Chan := make(chan struct{})

	go func() {
		for i := 0; i < 250; i++ {
			err := counter.DoAlternateRequest(i, counter.URL1, 3, 2)
			assert.Nil(suite.T(), err, "request error should be nil")
		}
		close(counter1Chan)
	}()
	go func() {
		for i := 0; i < 250; i++ {
			err := counter.DoAlternateRequest(i, counter.URL2, 4, 7)
			assert.Nil(suite.T(), err, "request error should be nil")
		}
		close(counter2Chan)
	}()

	<-counter1Chan
	<-counter2Chan

	counter1Resp, err := counter.DoGetCounterRequest(counter.URL1)
	assert.Nil(suite.T(), err, "request error should be nil")

	counter2Resp, err := counter.DoGetCounterRequest(counter.URL2)
	assert.Nil(suite.T(), err, "request error should be nil")

	assert.Equal(suite.T(), counter1Resp.Value, counter2Resp.Value, "replicas value should be equal")
}

func (suite *CounterE2ETestSuite) TestAppContainerFailureRestartShouldCatchUpStateCorrectlyBeforeReady() {
	counter1Chan := make(chan struct{}, 1)

	go func() {
		for i := 0; i < 250; i++ {
			err := counter.DoAlternateRequest(i, counter.URL1, 7, 11)
			assert.Nil(suite.T(), err, "request error should be nil")
		}
		counter1Chan <- struct{}{}
	}()

	err := counter.ExecuteAndWaitScriptFile(restartCounter2AppScriptPath)
	assert.Nil(suite.T(), err, "restart error should be nil")

	<-counter1Chan
	close(counter1Chan)

	counter1Resp, err := counter.DoGetCounterRequest(counter.URL1)
	assert.Nil(suite.T(), err, "request error should be nil")

	counter2Resp, err := counter.DoGetCounterRequest(counter.URL2)
	assert.Nil(suite.T(), err, "request error should be nil")

	assert.Equal(suite.T(), counter1Resp.Value, counter2Resp.Value, "replicas value should be equal")
}

func (suite *CounterE2ETestSuite) TestProxyContainerFailureRestartShouldCatchUpStateCorrectlyBeforeReady() {
	counter1Chan := make(chan struct{}, 1)

	go func() {
		for i := 0; i < 250; i++ {
			err := counter.DoAlternateRequest(i, counter.URL1, 7, 11)
			assert.Nil(suite.T(), err, "request error should be nil")
		}
		counter1Chan <- struct{}{}
	}()

	err := counter.ExecuteAndWaitScriptFile(restartCounter2ProxyScriptPath)
	assert.Nil(suite.T(), err, "restart error should be nil")

	<-counter1Chan
	close(counter1Chan)

	counter1Resp, err := counter.DoGetCounterRequest(counter.URL1)
	assert.Nil(suite.T(), err, "request error should be nil")

	counter2Resp, err := counter.DoGetCounterRequest(counter.URL2)
	assert.Nil(suite.T(), err, "request error should be nil")

	assert.Equal(suite.T(), counter1Resp.Value, counter2Resp.Value, "replicas value should be equal")
}

func TestCounterE2ETest(t *testing.T) {
	suite.Run(t, new(CounterE2ETestSuite))
}
