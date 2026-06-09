package paywechat

import (
	"net/http"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
)

// Callback implements model.PayClient. It handles WeChat Pay's async notify.
func (c Client) Callback(w http.ResponseWriter, r *http.Request) {
	// Parse notify params
	notifyReq, err := wechat.V3ParseNotify(r)
	if err != nil {
		c.handleErr(r, w, err)
		return
	}

	// Verify sign with the merchant public key map
	certMap := c.Client.WxPublicKeyMap()
	if err = notifyReq.VerifySignByPKMap(certMap); err != nil {
		c.handleErr(r, w, err)
		return
	}

	// Decrypt the cipher text into our struct
	callback := new(WeChatPayCallback)
	if err = notifyReq.DecryptCipherTextToStruct(c.Params.Credential[MerchantAPIv3Key], callback); err != nil {
		c.handleErr(r, w, err)
		return
	}

	// Push status. The framework-provided StatusUpdater trims prefix/suffix and
	// forwards to the user contract (no-op when no contract is configured).
	if c.Params.StatusUpdater != nil {
		if err = c.Params.StatusUpdater(
			r.Context(), r, Name,
			callback.OutTradeNo,
			MapState(callback.TradeState),
			FormatTime(callback.SuccessTime),
		); err != nil {
			c.handleErr(r, w, err)
			return
		}
	}

	// Acknowledge success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if body, mErr := c.Params.Marshal(&wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"}); mErr == nil {
		_, _ = w.Write(body)
	}
}

// handleErr routes an error through the user ErrorHandler when present and
// always returns a 500 so WeChat Pay retries the notification.
func (c Client) handleErr(r *http.Request, w http.ResponseWriter, err error) {
	if c.Params.ErrorHandler != nil {
		_ = c.Params.ErrorHandler(r.Context(), r, err)
	}
	w.WriteHeader(http.StatusInternalServerError)
}
