package global

const (
	SysCnfUserBanDuritonKey   = "userBanDuriton"
	SysCnfPaymentWayKey       = "paymentWay"
	SysRepeatOrderIntervalKey = "repeatOrderInterval"
)

const (
	DefaultUserBanDuriton      = "360"
	DefaultSysCnfPaymentWay    = "1"
	DefaultRepeatOrderInterval = "15"
)

const (
	BanTypeWord    = 1
	BanTypeMem     = 2
	BanTypeForword = 3
)

// telegram 相关设置
const (
	GroupTypeChat    = 1
	GroupTypeChannel = 2

	ButtonTypeKeyBoard = 1
	ButtonTypeInline   = 2
)

const (
	BotReplyCmdType = 1
	BotReplyCnfType = 2

	BotReplyCnfPublish2Channel = "publish_to_channel"
)

const (
	AdRechargeCreate      = 1 // 创建支付
	AdRechargePaid        = 2 // 完成支付
	AdRechargeTimeout     = 3 // 超时
	AdRechargePaidTimeout = 4 // 超时但完成
)
