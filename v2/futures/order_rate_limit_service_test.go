package futures

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type orderRateLimitServiceTestSuite struct {
	baseTestSuite
}

func TestOrderRateLimitService(t *testing.T) {
	suite.Run(t, new(orderRateLimitServiceTestSuite))
}

func (s *orderRateLimitServiceTestSuite) TestOrderRateLimit() {
	data := []byte(`
		[
			{
				"rateLimitType": "ORDERS",
				"interval": "SECOND",
				"intervalNum": 10,
				"limit": 10000
			},
			{
				"rateLimitType": "ORDERS",
				"interval": "MINUTE",
				"intervalNum": 1,
				"limit": 20000
			}
		]
	`)
	s.mockDo(data, nil)
	defer s.assertDo()

	s.assertReq(func(r *request) {
		e := newSignedRequest()
		s.assertRequestEqual(e, r)
	})

	res, _, err := s.client.NewOrderRateLimitService().Do(newContext())
	s.r().NoError(err)
	s.r().Len(res, 2)
	s.r().Equal(&RateLimit{RateLimitType: "ORDERS", Interval: "SECOND", IntervalNum: 10, Limit: 10000}, res[0])
	s.r().Equal(&RateLimit{RateLimitType: "ORDERS", Interval: "MINUTE", IntervalNum: 1, Limit: 20000}, res[1])
}

func (s *orderRateLimitServiceTestSuite) TestOrderRateLimitWithSymbol() {
	data := []byte(`
		[
			{
				"rateLimitType": "ORDERS",
				"interval": "MINUTE",
				"intervalNum": 1,
				"limit": 1200
			}
		]
	`)
	s.mockDo(data, nil)
	defer s.assertDo()

	symbol := "BTCUSDT"
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParam("symbol", symbol)
		s.assertRequestEqual(e, r)
	})

	res, _, err := s.client.NewOrderRateLimitService().Symbol(symbol).Do(newContext())
	s.r().NoError(err)
	s.r().Len(res, 1)
	s.r().Equal(int64(1200), res[0].Limit)
}
