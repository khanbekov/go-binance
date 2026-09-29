package futures

import (
	"context"
	"encoding/json"
	"net/http"
)

// OrderRateLimitService queries the order rate limits of the account
// (GET /fapi/v1/rateLimit/order). Unlike ExchangeInfo.RateLimits, which
// describes the default limits of the venue, the response carries the limits
// that are actually applied to the API key's account (VIP / market maker
// programs raise them).
type OrderRateLimitService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *OrderRateLimitService) Symbol(symbol string) *OrderRateLimitService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *OrderRateLimitService) Do(ctx context.Context, opts ...RequestOption) (res []*RateLimit, rateLimits map[string]string, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/rateLimit/order",
		secType:  secTypeSigned,
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	data, rateLimits, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, nil, err
	}
	res = make([]*RateLimit, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, nil, err
	}
	return res, rateLimits, nil
}
