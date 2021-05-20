package e2e_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	testUtils "k8s-smr/test/e2e/utilities"
	"testing"
)

type CounterE2ETestSuite struct {
	suite.Suite
}

func (suite *CounterE2ETestSuite) SetupSuite() {
}

func (suite *CounterE2ETestSuite) TearDownSuite() {
}

func (suite *CounterE2ETestSuite) TestRequestShouldPropagateToReplicasCorrectly() {
	err := testUtils.DoPostCounterRequest(testUtils.Counter1URL, "INC", 2)
	assert.Nil(suite.T(), err, "request error should be nil")

	err = testUtils.DoPostCounterRequest(testUtils.Counter1URL, "DEC", 1)
	assert.Nil(suite.T(), err, "request error should be nil")

	err = testUtils.DoPostCounterRequest(testUtils.Counter1URL, "INC", 3)
	assert.Nil(suite.T(), err, "request error should be nil")

	resp, err := testUtils.DoGetCounterRequest(testUtils.Counter2URL)
	assert.Nil(suite.T(), err, "request error should be nil")
	assert.Equal(suite.T(), 4, resp.Value, "replicas value should be equal")
}

func TestCounterE2ETest(t *testing.T) {
	suite.Run(t, new(CounterE2ETestSuite))
}
