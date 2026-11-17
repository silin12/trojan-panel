package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"trojan-panel/model"
	"trojan-panel/model/bo"
	"trojan-panel/model/constant"
	"trojan-panel/model/vo"
	"trojan-panel/util"
)

func uintToStr(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}

// BuildV2RaySubscription 生成 v2rayN / Shadowrocket / NekoBox 等客户端通用的分享链接
func BuildV2RaySubscription(proxies []interface{}) []string {
	links := make([]string, 0, len(proxies))
	for _, item := range proxies {
		if link, ok := BuildV2RayLink(item); ok {
			links = append(links, link)
		}
	}
	return links
}

// BuildV2RayLink 把单个代理转换成分享链接
func BuildV2RayLink(proxy interface{}) (string, bool) {
	switch p := proxy.(type) {
	case bo.Vless:
		params := url.Values{}
		params.Set("encryption", "none")
		security := "none"
		if p.RealityOpts.PublicKey != "" {
			security = "reality"
		} else if p.Tls {
			security = "tls"
		}
		params.Set("security", security)
		if p.ServerName != "" {
			params.Set("sni", p.ServerName)
		}
		if p.ClientFingerprint != "" {
			params.Set("fp", p.ClientFingerprint)
		}
		if p.Flow != "" {
			params.Set("flow", p.Flow)
		}
		if p.RealityOpts.PublicKey != "" {
			params.Set("pbk", p.RealityOpts.PublicKey)
			if p.RealityOpts.ShortId != "" {
				params.Set("sid", p.RealityOpts.ShortId)
			}
		}
		if p.SkipCertVerify {
			params.Set("allowInsecure", "1")
		}
		network := p.Network
		if network == "" {
			network = "tcp"
		}
		params.Set("type", network)
		if network == "ws" {
			if p.WsOpts.Path != "" {
				params.Set("path", p.WsOpts.Path)
			}
			if p.WsOpts.Headers.Host != "" {
				params.Set("host", p.WsOpts.Headers.Host)
			}
		}
		return fmt.Sprintf("vless://%s@%s:%s?%s#%s", p.Uuid, p.Server, uintToStr(p.Port),
			params.Encode(), url.PathEscape(p.Name)), true
	case bo.Vmess:
		network := p.Network
		if network == "" {
			network = "tcp"
		}
		tls := ""
		if p.Tls {
			tls = "tls"
		}
		host := ""
		path := ""
		if network == "ws" {
			host = p.WsOpts.Headers.Host
			path = p.WsOpts.Path
		}
		vmess := map[string]interface{}{
			"v":    "2",
			"ps":   p.Name,
			"add":  p.Server,
			"port": uintToStr(p.Port),
			"id":   p.Uuid,
			"aid":  uintToStr(p.AlterId),
			"scy":  p.Cipher,
			"net":  network,
			"type": "none",
			"host": host,
			"path": path,
			"tls":  tls,
			"sni":  p.ServerName,
			"fp":   p.ClientFingerprint,
		}
		vmessByte, err := json.Marshal(vmess)
		if err != nil {
			return "", false
		}
		return "vmess://" + base64.StdEncoding.EncodeToString(vmessByte), true
	case bo.Trojan:
		params := url.Values{}
		if p.Sni != "" {
			params.Set("sni", p.Sni)
		}
		if p.ClientFingerprint != "" {
			params.Set("fp", p.ClientFingerprint)
		}
		if len(p.Alpn) > 0 {
			params.Set("alpn", strings.Join(p.Alpn, ","))
		}
		if p.SkipCertVerify {
			params.Set("allowInsecure", "1")
		}
		network := "tcp"
		if p.WsOpts.Path != "" || p.WsOpts.Headers.Host != "" {
			network = "ws"
			if p.WsOpts.Path != "" {
				params.Set("path", p.WsOpts.Path)
			}
			if p.WsOpts.Headers.Host != "" {
				params.Set("host", p.WsOpts.Headers.Host)
			}
		}
		params.Set("type", network)
		return fmt.Sprintf("trojan://%s@%s:%s?%s#%s", p.Password, p.Server, uintToStr(p.Port),
			params.Encode(), url.PathEscape(p.Name)), true
	case bo.TrojanGo:
		params := url.Values{}
		if p.SNI != "" {
			params.Set("sni", p.SNI)
		}
		network := "tcp"
		if p.Network == "ws" {
			network = "ws"
			if p.WsOpts.Path != "" {
				params.Set("path", p.WsOpts.Path)
			}
			if p.WsOpts.Headers.Host != "" {
				params.Set("host", p.WsOpts.Headers.Host)
			}
		}
		params.Set("type", network)
		return fmt.Sprintf("trojan://%s@%s:%s?%s#%s", p.Password, p.Server, uintToStr(p.Port),
			params.Encode(), url.PathEscape(p.Name)), true
	case bo.Shadowsocks:
		userInfo := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", p.Cipher, p.Password)))
		return fmt.Sprintf("ss://%s@%s:%s#%s", userInfo, p.Server, uintToStr(p.Port), url.PathEscape(p.Name)), true
	case bo.Socks:
		userInfo := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", p.Username, p.Password)))
		return fmt.Sprintf("socks://%s@%s:%s#%s", userInfo, p.Server, uintToStr(p.Port), url.PathEscape(p.Name)), true
	case bo.Hysteria:
		params := url.Values{}
		protocol := p.Protocol
		if protocol == "" {
			protocol = "udp"
		}
		params.Set("protocol", protocol)
		if p.AuthStr != "" {
			params.Set("auth", p.AuthStr)
		}
		if p.Sni != "" {
			params.Set("peer", p.Sni)
		}
		if p.SkipCertVerify {
			params.Set("insecure", "1")
		}
		if p.Up > 0 {
			params.Set("upmbps", strconv.Itoa(p.Up))
		}
		if p.Down > 0 {
			params.Set("downmbps", strconv.Itoa(p.Down))
		}
		if p.Obfs != "" {
			params.Set("obfs", p.Obfs)
		}
		if p.FastOpen {
			params.Set("fastopen", "1")
		}
		return fmt.Sprintf("hysteria://%s:%s?%s#%s", p.Server, uintToStr(p.Port),
			params.Encode(), url.PathEscape(p.Name)), true
	case bo.Hysteria2:
		params := url.Values{}
		if p.Sni != "" {
			params.Set("sni", p.Sni)
		}
		if p.SkipCertVerify {
			params.Set("insecure", "1")
		}
		if p.Obfs != "" {
			params.Set("obfs", p.Obfs)
			if p.ObfsPassword != "" {
				params.Set("obfs-password", p.ObfsPassword)
			}
		}
		return fmt.Sprintf("hysteria2://%s@%s:%s/?%s#%s", p.Password, p.Server, uintToStr(p.Port),
			params.Encode(), url.PathEscape(p.Name)), true
	}
	return "", false
}

// SubscribeV2Ray 生成 v2rayN 等客户端使用的 base64 订阅内容
func SubscribeV2Ray(pass string) (*model.Account, string, string, error) {
	account, userInfo, proxies, _, _, err := buildClashProxies(pass)
	if err != nil {
		return nil, "", "", err
	}
	links := BuildV2RaySubscription(proxies)
	if len(links) == 0 {
		return nil, "", "", errors.New(constant.SysError)
	}
	content := strings.Join(links, "\n")
	return account, userInfo, base64.StdEncoding.EncodeToString([]byte(content)), nil
}

func singBoxTLS(serverName, fingerprint string, insecure bool, realityPbk, realityShortId string) map[string]interface{} {
	tlsSettings := map[string]interface{}{"enabled": true}
	if serverName != "" {
		tlsSettings["server_name"] = serverName
	}
	if insecure {
		tlsSettings["insecure"] = true
	}
	if fingerprint != "" {
		tlsSettings["utls"] = map[string]interface{}{"enabled": true, "fingerprint": fingerprint}
	}
	if realityPbk != "" {
		reality := map[string]interface{}{"enabled": true, "public_key": realityPbk}
		if realityShortId != "" {
			reality["short_id"] = realityShortId
		}
		tlsSettings["reality"] = reality
	}
	return tlsSettings
}

func singBoxTransport(path, host string) map[string]interface{} {
	transport := map[string]interface{}{"type": "ws", "path": path}
	if host != "" {
		transport["headers"] = map[string]string{"Host": host}
	}
	return transport
}

// BuildSingBoxOutbound 生成 sing-box 出站配置
func BuildSingBoxOutbound(proxy interface{}) (map[string]interface{}, bool) {
	switch p := proxy.(type) {
	case bo.Vless:
		outbound := map[string]interface{}{
			"type":        "vless",
			"tag":         p.Name,
			"server":      p.Server,
			"server_port": p.Port,
			"uuid":        p.Uuid,
		}
		if p.Flow != "" {
			outbound["flow"] = p.Flow
		}
		if p.Tls || p.RealityOpts.PublicKey != "" {
			outbound["tls"] = singBoxTLS(p.ServerName, p.ClientFingerprint, p.SkipCertVerify, p.RealityOpts.PublicKey, p.RealityOpts.ShortId)
		}
		if p.Network == "ws" {
			outbound["transport"] = singBoxTransport(p.WsOpts.Path, p.WsOpts.Headers.Host)
		}
		return outbound, true
	case bo.Vmess:
		outbound := map[string]interface{}{
			"type":        "vmess",
			"tag":         p.Name,
			"server":      p.Server,
			"server_port": p.Port,
			"uuid":        p.Uuid,
			"alter_id":    p.AlterId,
			"security":    "auto",
		}
		if p.Cipher == "none" {
			outbound["security"] = "none"
		}
		if p.Tls {
			outbound["tls"] = singBoxTLS(p.ServerName, p.ClientFingerprint, p.SkipCertVerify, "", "")
		}
		if p.Network == "ws" {
			outbound["transport"] = singBoxTransport(p.WsOpts.Path, p.WsOpts.Headers.Host)
		}
		return outbound, true
	case bo.Trojan:
		outbound := map[string]interface{}{
			"type":        "trojan",
			"tag":         p.Name,
			"server":      p.Server,
			"server_port": p.Port,
			"password":    p.Password,
			"tls":         singBoxTLS(p.Sni, p.ClientFingerprint, p.SkipCertVerify, "", ""),
		}
		if p.WsOpts.Path != "" || p.WsOpts.Headers.Host != "" {
			outbound["transport"] = singBoxTransport(p.WsOpts.Path, p.WsOpts.Headers.Host)
		}
		return outbound, true
	case bo.TrojanGo:
		outbound := map[string]interface{}{
			"type":        "trojan",
			"tag":         p.Name,
			"server":      p.Server,
			"server_port": p.Port,
			"password":    p.Password,
			"tls":         singBoxTLS(p.SNI, "", false, "", ""),
		}
		if p.Network == "ws" {
			outbound["transport"] = singBoxTransport(p.WsOpts.Path, p.WsOpts.Headers.Host)
		}
		return outbound, true
	case bo.Shadowsocks:
		return map[string]interface{}{
			"type":        "shadowsocks",
			"tag":         p.Name,
			"server":      p.Server,
			"server_port": p.Port,
			"method":      p.Cipher,
			"password":    p.Password,
		}, true
	case bo.Socks:
		outbound := map[string]interface{}{
			"type":        "socks",
			"tag":         p.Name,
			"server":      p.Server,
			"server_port": p.Port,
			"version":     "5",
		}
		if p.Username != "" {
			outbound["username"] = p.Username
		}
		if p.Password != "" {
			outbound["password"] = p.Password
		}
		return outbound, true
	case bo.Hysteria:
		outbound := map[string]interface{}{
			"type":        "hysteria",
			"tag":         p.Name,
			"server":      p.Server,
			"server_port": p.Port,
			"tls":         singBoxTLS(p.Sni, "", p.SkipCertVerify, "", ""),
		}
		if p.AuthStr != "" {
			outbound["auth_str"] = p.AuthStr
		}
		if p.Up > 0 {
			outbound["up_mbps"] = p.Up
		}
		if p.Down > 0 {
			outbound["down_mbps"] = p.Down
		}
		if p.Obfs != "" {
			outbound["obfs"] = p.Obfs
		}
		return outbound, true
	case bo.Hysteria2:
		outbound := map[string]interface{}{
			"type":        "hysteria2",
			"tag":         p.Name,
			"server":      p.Server,
			"server_port": p.Port,
			"password":    p.Password,
			"tls":         singBoxTLS(p.Sni, "", p.SkipCertVerify, "", ""),
		}
		if p.Obfs != "" && p.ObfsPassword != "" {
			outbound["obfs"] = map[string]interface{}{"type": p.Obfs, "password": p.ObfsPassword}
		}
		return outbound, true
	}
	return nil, false
}

// SubscribeSingBox 生成 sing-box 配置文件
func SubscribeSingBox(pass string) (*model.Account, string, []byte, error) {
	account, userInfo, proxies, _, _, err := buildClashProxies(pass)
	if err != nil {
		return nil, "", nil, err
	}

	tags := make([]string, 0, len(proxies))
	outbounds := make([]map[string]interface{}, 0, len(proxies)+4)
	for _, item := range proxies {
		outbound, ok := BuildSingBoxOutbound(item)
		if !ok {
			continue
		}
		tags = append(tags, outbound["tag"].(string))
		outbounds = append(outbounds, outbound)
	}
	if len(tags) == 0 {
		return nil, "", nil, errors.New(constant.SysError)
	}

	selectorTargets := make([]string, 0, len(tags)+1)
	selectorTargets = append(selectorTargets, "auto")
	selectorTargets = append(selectorTargets, tags...)
	outbounds = append(outbounds,
		map[string]interface{}{
			"type":      "selector",
			"tag":       "PROXY",
			"outbounds": selectorTargets,
			"default":   "auto",
		},
		map[string]interface{}{
			"type":                        "urltest",
			"tag":                         "auto",
			"outbounds":                   tags,
			"url":                         "https://www.gstatic.com/generate_204",
			"interval":                    "300s",
			"tolerance":                   50,
			"interrupt_exist_connections": false,
		},
		map[string]interface{}{"type": "direct", "tag": "direct"},
		map[string]interface{}{"type": "dns", "tag": "dns-out"},
	)

	config := map[string]interface{}{
		"log": map[string]interface{}{"level": "info", "timestamp": true},
		"dns": map[string]interface{}{
			"servers": []map[string]interface{}{
				{"tag": "remote", "address": "https://1.1.1.1/dns-query", "detour": "PROXY"},
				{"tag": "local", "address": "223.5.5.5", "detour": "direct"},
			},
			"final": "remote",
		},
		"inbounds": []map[string]interface{}{
			{
				"type":                       "mixed",
				"tag":                        "mixed-in",
				"listen":                     "127.0.0.1",
				"listen_port":                2080,
				"sniff":                      true,
				"sniff_override_destination": false,
			},
		},
		"outbounds": outbounds,
		"route": map[string]interface{}{
			"rules": []map[string]interface{}{
				{"protocol": "dns", "outbound": "dns-out"},
			},
			"final":                 "PROXY",
			"auto_detect_interface": true,
		},
		"experimental": map[string]interface{}{
			"cache_file": map[string]interface{}{"enabled": true},
		},
	}

	configByte, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, "", nil, errors.New(constant.SysError)
	}
	return account, userInfo, configByte, nil
}

// BuildSubscribeInfoPage 浏览器打开订阅链接时展示的信息页
func BuildSubscribeInfoPage(account *model.Account, code string, systemConfig vo.SystemVo) string {
	var quota, download, upload, expire int64
	if account.Quota != nil {
		quota = int64(*account.Quota)
	}
	if account.Download != nil {
		download = int64(*account.Download)
	}
	if account.Upload != nil {
		upload = int64(*account.Upload)
	}
	if account.ExpireTime != nil {
		expire = int64(*account.ExpireTime)
	}

	used := download + upload
	percent := 0
	if quota > 0 {
		percent = int(used * 100 / quota)
	}
	usedText := util.FormatByte(used)
	totalText := util.FormatByte(quota)
	remainText := util.FormatByte(quota - used)
	expireText := "不限"
	if expire > 0 {
		expireText = formatExpireTimeUnix(expire)
	}

	page := constant.SubscribePageTemplate
	page = strings.ReplaceAll(page, "{{SYSTEM_NAME}}", systemConfig.SystemName)
	page = strings.ReplaceAll(page, "{{USERNAME}}", *account.Username)
	page = strings.ReplaceAll(page, "{{PERCENT}}", strconv.Itoa(percent))
	page = strings.ReplaceAll(page, "{{USED}}", usedText)
	page = strings.ReplaceAll(page, "{{TOTAL}}", totalText)
	page = strings.ReplaceAll(page, "{{REMAIN}}", remainText)
	page = strings.ReplaceAll(page, "{{EXPIRE}}", expireText)
	page = strings.ReplaceAll(page, "{{CODE}}", code)
	return page
}
