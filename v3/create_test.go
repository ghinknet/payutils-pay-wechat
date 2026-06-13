package paywechat

import (
	"context"
	"testing"
	"time"

	stderrors "errors"

	"go.gh.ink/payutils/v3/errors"
	"go.gh.ink/payutils/v3/model"
)

func newTestClient() Client {
	return Client{
		Params: model.PayDriverClientParam{
			NoNewPaymentWindows: time.Minute,
			SafetyMargin:        5 * time.Second,
			Endpoint:            "https://gh.ink",
			Credential:          map[string]string{AppID: "app", MerchantID: "mch"},
		},
		// Client (gopay) intentionally nil: all cases below return before any
		// network call is made.
	}
}

func TestCreate_DoesNotClaimForeignParam(t *testing.T) {
	c := newTestClient()
	claimed, result, err := c.Create(context.Background(), struct{ X int }{X: 1})
	if claimed {
		t.Error("claimed = true for a foreign param, want false")
	}
	if result != nil || err != nil {
		t.Errorf("got (%v, %v), want (nil, nil)", result, err)
	}
}

func TestCreate_RejectsInsufficientTime(t *testing.T) {
	c := newTestClient()
	claimed, _, err := c.Create(context.Background(), CreateParam{
		TradeID:  "T1",
		Platform: PlatformPC,
		Detail: model.TradeDetail{
			Currency: "CNY",
			Expiry:   time.Now().Add(time.Second),
		},
	})
	if !claimed {
		t.Fatal("claimed = false, want true")
	}
	if !stderrors.Is(err, errors.ErrNoEnoughTimeToPay) {
		t.Errorf("err = %v, want ErrNoEnoughTimeToPay", err)
	}
}

func TestCreate_JSAPIRequiresOpenID(t *testing.T) {
	c := newTestClient()
	for _, platform := range []string{PlatformMobile, PlatformWeChat} {
		claimed, _, err := c.Create(context.Background(), CreateParam{
			TradeID:  "T1",
			Platform: platform,
			OpenID:   "", // missing
			Detail: model.TradeDetail{
				Currency: "CNY",
				Expiry:   time.Now().Add(time.Hour),
			},
		})
		if !claimed {
			t.Fatalf("platform %s: claimed = false, want true", platform)
		}
		if !stderrors.Is(err, ErrWeChatOpenIDIsRequired) {
			t.Errorf("platform %s: err = %v, want ErrWeChatOpenIDIsRequired", platform, err)
		}
	}
}

func TestCreate_RejectsUnknownPlatform(t *testing.T) {
	c := newTestClient()
	claimed, _, err := c.Create(context.Background(), CreateParam{
		TradeID:  "T1",
		Platform: "telepathy",
		Detail: model.TradeDetail{
			Currency: "CNY",
			Expiry:   time.Now().Add(time.Hour),
		},
	})
	if !claimed {
		t.Fatal("claimed = false, want true")
	}
	if !stderrors.Is(err, errors.ErrUnsupportedPlatform) {
		t.Errorf("err = %v, want ErrUnsupportedPlatform", err)
	}
}
