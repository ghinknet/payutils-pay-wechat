package paywechat

import "go.gh.ink/payutils/v3/model"

// CreateParam is the WeChat-Pay-specific create param.
//
// The user prepares the trade detail themselves and passes it to
// (*client.Client).Create. The dispatcher hands the param to every registered
// driver; this driver only claims params of this concrete type.
type CreateParam struct {
	// TradeID is the user-side trade id (without prefix/suffix).
	TradeID string
	// Detail carries the order info: Subject / Price (cents) / Currency / Expiry.
	Detail model.TradeDetail
	// Platform selects the WeChat Pay product:
	//   PlatformPC               -> NATIVE (code_url)
	//   PlatformMobile/WeChat    -> JSAPI  (requires OpenID)
	Platform string
	// OpenID is required for JSAPI (mobile / in-WeChat) transactions.
	OpenID string
}
