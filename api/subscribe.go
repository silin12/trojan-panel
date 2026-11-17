package api

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"trojan-panel/model/constant"
	"trojan-panel/model/dto"
	"trojan-panel/model/vo"
	"trojan-panel/service"
	"trojan-panel/util"
)

// subscribeShortLink 生成账号对应的订阅短链
func subscribeShortLink(accountId uint, password string) string {
	return fmt.Sprintf("/api/s/%s", util.EncodeShortCode(accountId, password))
}

// ClashSubscribe 获取Clash订阅地址
func ClashSubscribe(c *gin.Context) {
	accountVo := service.GetCurrentAccount(c)
	password, err := service.SelectConnectPassword(&accountVo.Id, &accountVo.Username)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(subscribeShortLink(accountVo.Id, password), c)
}

// ClashSubscribeForSb 获取指定人的Clash订阅地址
func ClashSubscribeForSb(c *gin.Context) {
	var accountRequiredIdDto dto.RequiredIdDto
	_ = c.ShouldBindQuery(&accountRequiredIdDto)
	if err := validate.Struct(&accountRequiredIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	password, err := service.SelectConnectPassword(accountRequiredIdDto.Id, nil)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(subscribeShortLink(*accountRequiredIdDto.Id, password), c)
}

// Subscribe 订阅（兼容旧链接 /api/auth/subscribe/{base64(pass)}）
func Subscribe(c *gin.Context) {
	token := c.Param("token")
	tokenDecode, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	writeSubscribe(c, string(tokenDecode), detectSubscribeFormat(c))
}

// SubscribeShort 订阅短链 /s/:code，按 UA 或 format 参数返回不同客户端格式
func SubscribeShort(c *gin.Context) {
	code := c.Param("code")
	accountId, _, ok := util.DecodeShortCode(code)
	if !ok {
		vo.Fail(constant.SysError, c)
		return
	}
	password, err := service.SelectConnectPassword(&accountId, nil)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	if !util.CheckShortCode(code, accountId, password) {
		vo.Fail(constant.UnauthorizedError, c)
		return
	}
	writeSubscribe(c, password, detectSubscribeFormat(c))
}

// detectSubscribeFormat 按 format 参数或 User-Agent 判断输出格式
func detectSubscribeFormat(c *gin.Context) string {
	if format := strings.ToLower(strings.TrimSpace(c.Query("format"))); format != "" {
		switch format {
		case "clash", "clashmeta", "mihomo", "meta":
			return "clash"
		case "v2ray", "v2rayn", "v2rayng", "base64", "shadowrocket", "ssr", "nekobox", "nekoray":
			return "v2ray"
		case "singbox", "sing-box", "sfa", "sfi", "sft":
			return "singbox"
		case "info", "page", "html":
			return "info"
		}
	}

	userAgent := strings.ToLower(c.Request.Header.Get("User-Agent"))
	switch {
	case strings.Contains(userAgent, "clash"), strings.Contains(userAgent, "mihomo"),
		strings.Contains(userAgent, "stash"), strings.Contains(userAgent, "surfboard"),
		strings.Contains(userAgent, "clashx"), strings.Contains(userAgent, "flclash"),
		strings.Contains(userAgent, "clash verge"), strings.Contains(userAgent, "clash-verge"):
		return "clash"
	case strings.Contains(userAgent, "sing-box"), strings.Contains(userAgent, "singbox"),
		strings.Contains(userAgent, "sfa/"), strings.Contains(userAgent, "sfi/"),
		strings.Contains(userAgent, "sft/"), strings.Contains(userAgent, "hiddify"):
		return "singbox"
	case strings.Contains(userAgent, "v2rayn"), strings.Contains(userAgent, "v2rayng"),
		strings.Contains(userAgent, "v2rayu"), strings.Contains(userAgent, "shadowrocket"),
		strings.Contains(userAgent, "nekobox"), strings.Contains(userAgent, "nekoray"),
		strings.Contains(userAgent, "v2ray"), strings.Contains(userAgent, "quantumult"),
		strings.Contains(userAgent, "kitsunebi"), strings.Contains(userAgent, "ssr"),
		strings.Contains(userAgent, "surge"), strings.Contains(userAgent, "loon"),
		strings.Contains(userAgent, "passwall"), strings.Contains(userAgent, "openwrt"),
		strings.Contains(userAgent, "matsuri"), strings.Contains(userAgent, "clash.meta"),
		strings.Contains(userAgent, "v2raytun"), strings.Contains(userAgent, "happ"),
		strings.Contains(userAgent, "streisand"), strings.Contains(userAgent, "karing"),
		strings.Contains(userAgent, "shadowsocks"), strings.Contains(userAgent, "winxray"):
		return "v2ray"
	}

	// 没匹配到已知客户端：只有浏览器会主动要求 text/html，其余一律按 V2Ray(base64) 输出，
	// 这样 v2rayN 之类 UA 不规范的客户端也能直接用同一条短链。
	if strings.Contains(strings.ToLower(c.Request.Header.Get("Accept")), "text/html") {
		return "info"
	}
	return "v2ray"
}

// writeSubscribe 按格式输出订阅内容
func writeSubscribe(c *gin.Context, pass string, format string) {
	switch format {
	case "clash":
		account, userInfo, clashConfigYaml, systemConfig, err := service.SubscribeClash(pass)
		if err != nil {
			vo.Fail(err.Error(), c)
			return
		}
		result := fmt.Sprintf(`%s
%s`, string(clashConfigYaml), systemConfig.ClashRule)
		c.Header("content-disposition", fmt.Sprintf("attachment; filename=%s.yaml", *account.Username))
		c.Header("profile-update-interval", "12")
		c.Header("subscription-userinfo", userInfo)
		c.String(200, result)
	case "v2ray":
		account, userInfo, content, err := service.SubscribeV2Ray(pass)
		if err != nil {
			vo.Fail(err.Error(), c)
			return
		}
		c.Header("content-disposition", fmt.Sprintf("attachment; filename=%s.txt", *account.Username))
		c.Header("profile-update-interval", "12")
		c.Header("subscription-userinfo", userInfo)
		c.String(200, content)
	case "singbox":
		account, userInfo, configByte, err := service.SubscribeSingBox(pass)
		if err != nil {
			vo.Fail(err.Error(), c)
			return
		}
		c.Header("content-disposition", fmt.Sprintf("attachment; filename=%s.json", *account.Username))
		c.Header("profile-update-interval", "12")
		c.Header("subscription-userinfo", userInfo)
		c.Data(200, "application/json; charset=utf-8", configByte)
	default:
		account, _, _, systemConfig, err := service.SubscribeClash(pass)
		if err != nil {
			vo.Fail(err.Error(), c)
			return
		}
		// 浏览器访问：展示订阅信息页
		// 用账号自身的 id + pass 现算短码，这样旧长链接打开时页面里的四个格式链接同样可用
		code := util.EncodeShortCode(*account.Id, pass)
		page := service.BuildSubscribeInfoPage(account, code, systemConfig)
		c.Data(200, "text/html; charset=utf-8", []byte(page))
	}
}
