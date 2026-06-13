package paywechat

const Name = "wechat"

// Credential keys.
const (
	AppID                    = "appID"
	MerchantID               = "mchID"
	MerchantAPIv3Key         = "apiV3Key"
	MerchantCertSerialNumber = "certSerialNumber"
	MerchantPrivateKey       = "privateKey"
	PublicKey                = "publicKey"
	PublicKeyID              = "publicKeyID"
)

// Platform selects the WeChat Pay product.
const (
	// PlatformPC creates a NATIVE transaction (returns a code_url for QR code).
	PlatformPC = "pc"
	// PlatformMobile creates a JSAPI transaction (requires OpenID).
	PlatformMobile = "mobile"
	// PlatformWeChat creates a JSAPI transaction (requires OpenID).
	PlatformWeChat = "wechat"
)
