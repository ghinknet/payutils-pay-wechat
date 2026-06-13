# payutils-pay-wechat

WeChat Pay pay driver for [payutils](https://go.gh.ink/payutils/v3).

Wraps [go-pay/gopay](https://github.com/go-pay/gopay)'s WeChat Pay v3 client and
implements the `payutils` `PayClient` interface: trade creation (Native &
JSAPI), asynchronous notify (callback) handling, and status / close / refund
actions.

This driver covers **payment only**. 

## Module

```
go.gh.ink/payutils/pay/wechat/v3
```

```bash
go get go.gh.ink/payutils/pay/wechat/v3
```

## Registration

Registers itself under the name `wechat` via `init()`. Blank import to enable:

```go
import _ "go.gh.ink/payutils/pay/wechat/v3"
```

Import normally to use its constants and `CreateParam`:

```go
import payWeChat "go.gh.ink/payutils/pay/wechat/v3"
```

## Credentials

Provide these keys under `Credentials["wechat"]` in `model.Config`. Use the
exported constants instead of raw strings.

| Constant                   | Key                | Meaning |
|----------------------------|--------------------|---------|
| `AppID`                    | `appID`            | Application (公众号/小程序/APP) ID |
| `MerchantID`               | `mchID`            | Merchant ID |
| `MerchantAPIv3Key`         | `apiV3Key`         | Merchant APIv3 key (used to decrypt notifications) |
| `MerchantCertSerialNumber` | `certSerialNumber` | Merchant certificate serial number |
| `MerchantPrivateKey`       | `privateKey`       | Merchant private key |
| `PublicKey`                | `publicKey`        | WeChat Pay public key (for signature verification) |
| `PublicKeyID`              | `publicKeyID`      | WeChat Pay public key ID |

```go
Credentials: model.C{
    payWeChat.Name: {
        payWeChat.AppID:                    "...",
        payWeChat.MerchantID:               "...",
        payWeChat.MerchantAPIv3Key:         "...",
        payWeChat.MerchantCertSerialNumber: "...",
        payWeChat.MerchantPrivateKey:       "...",
        payWeChat.PublicKey:                "...",
        payWeChat.PublicKeyID:              "...",
    },
},
```

## Creating a trade

Pass `payWeChat.CreateParam` to `(*client.Client).Create`:

```go
result, err := c.Create(ctx, payWeChat.CreateParam{
    TradeID:  "order-123",
    Platform: payWeChat.PlatformWeChat,
    OpenID:   "user-openid",          // required for JSAPI platforms
    Detail: model.TradeDetail{
        Subject:  "A nice product",
        Price:    1990,               // cents
        Currency: "CNY",
        Expiry:   time.Now().Add(time.Hour),
    },
})
```

### Platforms

| Constant         | Value    | WeChat API             | Result | OpenID |
|------------------|----------|------------------------|--------|--------|
| `PlatformPC`     | `pc`     | `V3TransactionNative`  | `map[string]string{"payUrl": code_url}` (QR) | not needed |
| `PlatformMobile` | `mobile` | `V3TransactionJsapi`   | JSAPI sign object | **required** |
| `PlatformWeChat` | `wechat` | `V3TransactionJsapi`   | JSAPI sign object | **required** |

### Result

`result` is `any`:

- Native (`PlatformPC`) → `map[string]string{"payUrl": "<code_url>"}`.
- JSAPI (`PlatformMobile` / `PlatformWeChat`) → the gopay JSAPI sign object,
  ready to hand to `wx.requestPayment` on the front-end.

Marshal it directly, or type-assert when you need the concrete value.

### Validation errors

`Create` returns these (matchable with `errors.Is`) before any network call:

- `errors.ErrNoEnoughTimeToPay` — too little time remains before `Expiry`.
- `ErrWeChatOpenIDIsRequired` — JSAPI platform without an `OpenID`.
- `errors.ErrUnsupportedPlatform` — unknown `Platform`.

## Callback

With an HTTP driver configured, payutils auto-registers
`POST /wechat/callback`. The driver:

1. parses the v3 notification,
2. verifies the signature against the merchant public-key map,
3. decrypts the ciphertext with `apiV3Key`,
4. pushes the mapped status to your `Contract.StatusUpdater`,
5. replies with the WeChat `SUCCESS` acknowledgement.

To handle it manually instead, call `c.Callback("wechat", w, r)`.

## State mapping

| WeChat state | payutils state |
|--------------|----------------|
| `NOTPAY`     | `PENDING` |
| `SUCCESS`    | `SUCCESS` |
| `REFUND`     | `CLOSED` |
| `CLOSED`     | `CLOSED` |
| (anything else) | `UNKNOWN` |

## License

See [LICENSE](LICENSE).
