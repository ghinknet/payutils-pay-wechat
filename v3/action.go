package paywechat

import (
	"context"
	"strings"

	"github.com/go-pay/gopay"

	"go.gh.ink/payutils/v3/errors"
	"go.gh.ink/payutils/v3/model"
)

func (c Client) Status(tradeID string) (model.ReturnStatus, error) {
	outTradeNo := strings.Join([]string{c.Params.TradeIDPrefix, tradeID, c.Params.TradeIDSuffix}, "")

	// Query order by out_trade_no (idType = 2)
	wxRsp, err := c.Client.V3TransactionQueryOrder(context.Background(), 2, outTradeNo)
	if err != nil {
		return model.ReturnStatus{}, err
	}

	// Check return status (404 means trade not exist)
	if wxRsp.Code != 0 && wxRsp.Code != 404 {
		return model.ReturnStatus{}, ErrWeChatPayRespCodeInvalid.
			WithUpstreamName(Name).
			WithUpstreamCode(wxRsp.ErrResponse.Code).
			WithUpstreamMessage(wxRsp.ErrResponse.Message).
			WithUpstreamResponse(wxRsp)
	}
	if wxRsp.Code == 404 {
		return model.ReturnStatus{}, errors.ErrTradeNotExist.
			WithUpstreamName(Name).
			WithUpstreamCode(wxRsp.ErrResponse.Code).
			WithUpstreamMessage(wxRsp.ErrResponse.Message).
			WithUpstreamResponse(wxRsp)
	}

	return model.ReturnStatus{
		TradeStatus: MapState(wxRsp.Response.TradeState),
		Upstream:    Name,
		Time:        FormatTime(wxRsp.Response.SuccessTime),
	}, nil
}

func (c Client) Close(tradeID string) error {
	outTradeNo := strings.Join([]string{c.Params.TradeIDPrefix, tradeID, c.Params.TradeIDSuffix}, "")

	wxRsp, err := c.Client.V3TransactionCloseOrder(context.Background(), outTradeNo)
	if err != nil {
		return err
	}

	if wxRsp.Code != 0 && wxRsp.Code != 404 {
		return ErrWeChatPayRespCodeInvalid.
			WithUpstreamName(Name).
			WithUpstreamCode(wxRsp.ErrResponse.Code).
			WithUpstreamMessage(wxRsp.ErrResponse.Message).
			WithUpstreamResponse(wxRsp)
	}
	if wxRsp.Code == 404 {
		return errors.ErrTradeNotExist.
			WithUpstreamName(Name).
			WithUpstreamCode(wxRsp.ErrResponse.Code).
			WithUpstreamMessage(wxRsp.ErrResponse.Message).
			WithUpstreamResponse(wxRsp)
	}

	return nil
}

func (c Client) Refund(tradeID string, curr string, refundID string, totalAmount int64, refundAmount int64, reason string) error {
	outTradeNo := strings.Join([]string{c.Params.TradeIDPrefix, tradeID, c.Params.TradeIDSuffix}, "")
	outRefundNo := strings.Join([]string{c.Params.TradeIDPrefix, refundID, c.Params.TradeIDSuffix}, "")
	notifyURL := strings.Join([]string{c.Params.Endpoint, "/", Name, "/callback"}, "")

	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", outTradeNo).
		Set("out_refund_no", outRefundNo).
		Set("reason", reason).
		Set("notify_url", notifyURL).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", totalAmount).
				Set("refund", refundAmount).
				Set("currency", curr)
		})

	wxRsp, err := c.Client.V3Refund(context.Background(), bm)
	if err != nil {
		return err
	}

	if wxRsp.Code != 0 && wxRsp.Code != 404 {
		return ErrWeChatPayRespCodeInvalid.
			WithUpstreamName(Name).
			WithUpstreamCode(wxRsp.ErrResponse.Code).
			WithUpstreamMessage(wxRsp.ErrResponse.Message).
			WithUpstreamResponse(wxRsp)
	}
	if wxRsp.Code == 404 {
		return errors.ErrTradeNotExist.
			WithUpstreamName(Name).
			WithUpstreamCode(wxRsp.ErrResponse.Code).
			WithUpstreamMessage(wxRsp.ErrResponse.Message).
			WithUpstreamResponse(wxRsp)
	}

	return nil
}
