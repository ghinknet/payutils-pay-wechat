package paywechat

import (
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"

	"go.gh.ink/payutils/v3/model"
)

type Client struct {
	Params model.PayDriverClientParam
	Client *wechat.ClientV3
}

type Driver struct{}

func (d Driver) NewClient(params model.PayDriverClientParam) (model.PayClient, error) {
	// Create wechat-pay v3 client
	client, err := wechat.NewClientV3(
		params.Credential[MerchantID],
		params.Credential[MerchantCertSerialNumber],
		params.Credential[MerchantAPIv3Key],
		params.Credential[MerchantPrivateKey],
	)
	if err != nil {
		return nil, err
	}

	// Auto verify sign by public key
	if err = client.AutoVerifySignByPublicKey(
		[]byte(params.Credential[PublicKey]),
		params.Credential[PublicKeyID],
	); err != nil {
		return nil, err
	}

	// Debug switch
	if params.Debug {
		client.DebugSwitch = gopay.DebugOn
	} else {
		client.DebugSwitch = gopay.DebugOff
	}

	return Client{
		Params: params,
		Client: client,
	}, nil
}
