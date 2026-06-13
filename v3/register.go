package paywechat

import "go.gh.ink/payutils/v3/driver"

func init() {
	driver.RegisterPay(Name, Driver{})
}
