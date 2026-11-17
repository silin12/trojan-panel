package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"strings"
)

const (
	shortCodeAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	shortCodeSecret   = "trojan-panel-short-link-2024"
)

func base62Encode(n uint) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0)
	for n > 0 {
		buf = append(buf, shortCodeAlphabet[n%62])
		n /= 62
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

func base62Decode(s string) (uint, bool) {
	var n uint
	if s == "" {
		return 0, false
	}
	for _, c := range s {
		idx := strings.IndexRune(shortCodeAlphabet, c)
		if idx < 0 {
			return 0, false
		}
		n = n*62 + uint(idx)
	}
	return n, true
}

// shortCodeSig 用 pass 做密钥派生，防止别人猜到 id 就能拉到订阅
func shortCodeSig(id uint, pass string) string {
	mac := hmac.New(sha1.New, []byte(shortCodeSecret))
	mac.Write([]byte(strconv.FormatUint(uint64(id), 10)))
	mac.Write([]byte("|"))
	mac.Write([]byte(pass))
	return hex.EncodeToString(mac.Sum(nil))[:8]
}

// EncodeShortCode 生成订阅短链的 code：base62(id) + "-" + 8位签名
func EncodeShortCode(id uint, pass string) string {
	return base62Encode(id) + "-" + shortCodeSig(id, pass)
}

// DecodeShortCode 解析短链 code，返回 accountId 和签名
func DecodeShortCode(code string) (uint, string, bool) {
	parts := strings.Split(code, "-")
	if len(parts) != 2 {
		return 0, "", false
	}
	id, ok := base62Decode(parts[0])
	if !ok {
		return 0, "", false
	}
	return id, parts[1], true
}

// CheckShortCode 校验 code 是否与账户匹配
func CheckShortCode(code string, id uint, pass string) bool {
	return code == EncodeShortCode(id, pass)
}
