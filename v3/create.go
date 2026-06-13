package paywechat

import (
	"context"
	"strings"
	"time"

	"github.com/go-pay/gopay"

	"go.gh.ink/payutils/v3/errors"
	"go.gh.ink/payutils/v3/router"
)

// Create implements model.PayClient. It only claims params of type CreateParam.
func (c Client) Create(ctx context.Context, param any) (bool, any, error) {
	// Assert ownership of the param.
	p, ok := param.(CreateParam)
	if !ok {
		return false, nil, nil
	}

	// Check time: there must be enough room before expiry to open a payment.
	if time.Until(p.Detail.Expiry) < c.Params.NoNewPaymentWindows {
		return true, nil, errors.ErrNoEnoughTimeToPay.WithUpstreamName(Name)
	}

	// WeChat Pay expects an RFC3339 timestamp. Pull back by the safety margin.
	expire := p.Detail.Expiry.Add(-c.Params.SafetyMargin).Format(time.RFC3339)

	outTradeNo := strings.Join([]string{c.Params.TradeIDPrefix, p.TradeID, c.Params.TradeIDSuffix}, "")
	notifyURL := router.Notify(c.Params.Endpoint, Name)

	// Prepare params
	bm := make(gopay.BodyMap)
	bm.Set("appid", c.Params.Credential[AppID]).
		Set("mchid", c.Params.Credential[MerchantID]).
		Set("description", p.Detail.Subject).
		Set("out_trade_no", outTradeNo).
		Set("time_expire", expire).
		Set("notify_url", notifyURL).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", p.Detail.Price).
				Set("currency", p.Detail.Currency)
		})

	switch p.Platform {
	case PlatformPC:
		// NATIVE transaction (QR code)
		wxRsp, err := c.Client.V3TransactionNative(ctx, bm)
		if err != nil {
			return true, nil, err
		}
		if wxRsp.Code != 0 {
			return true, nil, ErrWeChatPayRespCodeInvalid.
				WithUpstreamName(Name).
				WithUpstreamCode(wxRsp.ErrResponse.Code).
				WithUpstreamMessage(wxRsp.ErrResponse.Message).
				WithUpstreamResponse(wxRsp)
		}
		return true, map[string]string{"payUrl": wxRsp.Response.CodeUrl}, nil

	case PlatformMobile, PlatformWeChat:
		// JSAPI transaction (requires OpenID)
		if p.OpenID == "" {
			return true, nil, ErrWeChatOpenIDIsRequired.WithUpstreamName(Name)
		}
		bm.SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", p.OpenID)
		})

		wxRsp, err := c.Client.V3TransactionJsapi(ctx, bm)
		if err != nil {
			return true, nil, err
		}
		if wxRsp.Code != 0 {
			return true, nil, ErrWeChatPayRespCodeInvalid.
				WithUpstreamName(Name).
				WithUpstreamCode(wxRsp.ErrResponse.Code).
				WithUpstreamMessage(wxRsp.ErrResponse.Message).
				WithUpstreamResponse(wxRsp)
		}

		// Build the JSAPI sign object for the front-end.
		jsapi, err := c.Client.PaySignOfJSAPI(c.Params.Credential[AppID], wxRsp.Response.PrepayId)
		if err != nil {
			return true, nil, err
		}
		return true, jsapi, nil

	default:
		return true, nil, errors.ErrUnsupportedPlatform.WithUpstreamName(Name)
	}
}
