package constant

const (
	TgStatusInit             = 0 // 未开始
	TgStatusCodeSent         = 1 // 已发送验证码
	TgStatusPasswordRequired = 2 // 需要二步验证
	TgStatusReady            = 3 // 登录完成，可用
)
