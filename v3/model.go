package paywechat

import (
	"time"

	"go.gh.ink/payutils/v3/model"
)

const (
	// TradeStateNotPay Trade created and waiting for buyer to pay
	// -> TradeStatePending
	TradeStateNotPay string = "NOTPAY"
	// TradeStateSuccess Trade successes
	// -> TradeStateSuccess
	TradeStateSuccess string = "SUCCESS"
	// TradeStateRefund Trade refunded
	// -> TradeStateClosed
	TradeStateRefund string = "REFUND"
	// TradeStateClosed Trade closed due to time out
	// -> TradeStateClosed
	TradeStateClosed string = "CLOSED"
)

// TradeStateMap maps WeChat trade state to the internal trade state.
var TradeStateMap = map[string]model.TradeState{
	TradeStateNotPay:  model.TradeStatePending,
	TradeStateSuccess: model.TradeStateSuccess,
	TradeStateRefund:  model.TradeStateClosed,
	TradeStateClosed:  model.TradeStateClosed,
}

// MapState maps a WeChat trade state (string) to the internal trade state.
func MapState(state string) model.TradeState {
	internalState, ok := TradeStateMap[state]
	if !ok {
		return model.TradeStateUnknown
	}
	return internalState
}

// FormatTime parses a WeChat RFC3339 time string into time.Time.
func FormatTime(timeStr string) time.Time {
	if timeStr == "" {
		return time.Time{}
	}
	timeObj, err := time.ParseInLocation(time.RFC3339, timeStr, time.Local)
	if err != nil {
		return time.Time{}
	}
	return timeObj
}

// WeChatPayCallback is the decrypted body of an async payment notification.
type WeChatPayCallback struct {
	TransactionID   string             `json:"transaction_id"`
	Amount          AmountInfo         `json:"amount"`
	MchID           string             `json:"mchid"`
	TradeState      string             `json:"trade_state"`
	BankType        string             `json:"bank_type"`
	PromotionDetail []*PromotionDetail `json:"promotion_detail,omitempty"`
	SuccessTime     string             `json:"success_time"`
	Payer           PayerInfo          `json:"payer"`
	OutTradeNo      string             `json:"out_trade_no"`
	AppID           string             `json:"appid"`
	TradeStateDesc  string             `json:"trade_state_desc"`
	TradeType       string             `json:"trade_type"`
	Attach          string             `json:"attach,omitempty"`
	SceneInfo       SceneInfo          `json:"scene_info,omitempty"`
}

type AmountInfo struct {
	PayerTotal    int    `json:"payer_total"`
	Total         int    `json:"total"`
	Currency      string `json:"currency"`
	PayerCurrency string `json:"payer_currency"`
}

type PromotionDetail struct {
	Amount              int            `json:"amount"`
	WeChatPayContribute int            `json:"wechatpay_contribute"`
	CouponID            string         `json:"coupon_id"`
	Scope               string         `json:"scope"`
	MerchantContribute  int            `json:"merchant_contribute"`
	Name                string         `json:"name"`
	OtherContribute     int            `json:"other_contribute"`
	Currency            string         `json:"currency"`
	StockID             string         `json:"stock_id"`
	GoodsDetail         []*GoodsDetail `json:"goods_detail,omitempty"`
}

type GoodsDetail struct {
	GoodsRemark    string `json:"goods_remark"`
	Quantity       int    `json:"quantity"`
	DiscountAmount int    `json:"discount_amount"`
	GoodsID        string `json:"goods_id"`
	UnitPrice      int    `json:"unit_price"`
}

type PayerInfo struct {
	OpenID string `json:"openid"`
}

type SceneInfo struct {
	DeviceID string `json:"device_id"`
}
