package data

// SponsorDecryptKeyHex 由主程序在启动时同步为 ldflags 注入的 BuildKey。
var SponsorDecryptKeyHex string

// EffectiveSponsorVipLevel 始终返回最高VIP等级，无需验证。
func EffectiveSponsorVipLevel() (level int, active bool) {
	return 999, true
}
