package paywechat

import "go.gh.ink/payutils/v3/errors"

var ErrWeChatPayRespCodeInvalid = errors.New("wechat pay resp code invalid")
var ErrWeChatPayNotifyVerifyFailed = errors.New("wechat pay notify verify failed")
var ErrWeChatOpenIDIsRequired = errors.New("wechat open id is required")
