package data

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/duke-git/lancet/v2/cryptor"
)

// DefaultSponsorAESKeyHex 与 main.checkDir 在 BuildKey 为空时的回退值一致，
// 供 ai-assistant-web 等独立进程解密本地配置中的赞助码。
const DefaultSponsorAESKeyHex = ""

// SponsorDecryptKeyHex 由主程序在启动时同步为 ldflags 注入的 BuildKey；为空则使用 DefaultSponsorAESKeyHex。
var SponsorDecryptKeyHex string

// SafeDecryptSponsorCode 解密赞助码（hex 解码 + AES-ECB）。
// 预校验密钥长度（16/24/32 字节）与密文块长度（16 字节整数倍），并 defer recover 兜底
// lancet AesEcbDecrypt 对非法输入（如密钥不匹配导致 PKCS#7 填充非法）的直接 panic；
// 任何失败以 error 返回，绝不使调用方进程崩溃。
func SafeDecryptSponsorCode(sponsorCode, keyHex string) (raw []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			raw, err = nil, fmt.Errorf("赞助码解密异常: %v", r)
		}
	}()
	sponsorCode = strings.TrimSpace(sponsorCode)
	if sponsorCode == "" {
		return nil, fmt.Errorf("赞助码为空")
	}
	encrypted, err := hex.DecodeString(sponsorCode)
	if err != nil {
		return nil, fmt.Errorf("赞助码 hex 解码失败: %w", err)
	}
	key, err := hex.DecodeString(strings.TrimSpace(keyHex))
	if err != nil {
		return nil, fmt.Errorf("赞助码密钥 hex 解码失败: %w", err)
	}
	if l := len(key); l != 16 && l != 24 && l != 32 {
		return nil, fmt.Errorf("赞助码密钥长度非法: %d 字节", l)
	}
	if len(encrypted) == 0 || len(encrypted)%16 != 0 {
		return nil, fmt.Errorf("赞助码密文长度非法: %d 字节（须为 16 字节整数倍）", len(encrypted))
	}
	return cryptor.AesEcbDecrypt(encrypted, key), nil
}

// EffectiveSponsorVipLevel 始终返回最高VIP等级，无需验证。
func EffectiveSponsorVipLevel() (level int, active bool) {
	return 999, true
}
