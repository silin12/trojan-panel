package vo

// NodeExportVo 节点导出对象
// 主表字段 + 各协议子表字段，用于面板迁移或批量备份
type NodeExportVo struct {
	Id                 uint   `json:"id"`
	Name               string `json:"name"`
	NodeServerIp       string `json:"nodeServerIp"`
	NodeServerGrpcPort uint   `json:"nodeServerGrpcPort"`
	NodeTypeId         uint   `json:"nodeTypeId"`
	NodeTypeName       string `json:"nodeTypeName"`
	Domain             string `json:"domain"`
	Port               uint   `json:"port"`
	Priority           int    `json:"priority"`

	XrayProtocol       string `json:"xrayProtocol"`
	XrayFlow           string `json:"xrayFlow"`
	XraySSMethod       string `json:"xraySSMethod"`
	RealityPbk         string `json:"realityPbk"`
	XraySettings       string `json:"xraySettings"`
	XrayStreamSettings string `json:"xrayStreamSettings"`
	XrayTag            string `json:"xrayTag"`
	XraySniffing       string `json:"xraySniffing"`
	XrayAllocate       string `json:"xrayAllocate"`

	TrojanGoSni             string `json:"trojanGoSni"`
	TrojanGoMuxEnable       uint   `json:"trojanGoMuxEnable"`
	TrojanGoWebsocketEnable uint   `json:"trojanGoWebsocketEnable"`
	TrojanGoWebsocketPath   string `json:"trojanGoWebsocketPath"`
	TrojanGoWebsocketHost   string `json:"trojanGoWebsocketHost"`
	TrojanGoSsEnable        uint   `json:"trojanGoSsEnable"`
	TrojanGoSsMethod        string `json:"trojanGoSsMethod"`
	TrojanGoSsPassword      string `json:"trojanGoSsPassword"`

	HysteriaProtocol   string `json:"hysteriaProtocol"`
	HysteriaObfs       string `json:"hysteriaObfs"`
	HysteriaUpMbps     int    `json:"hysteriaUpMbps"`
	HysteriaDownMbps   int    `json:"hysteriaDownMbps"`
	HysteriaServerName string `json:"hysteriaServerName"`
	HysteriaInsecure   uint   `json:"hysteriaInsecure"`
	HysteriaFastOpen   uint   `json:"hysteriaFastOpen"`

	Hysteria2ObfsPassword string `json:"hysteria2ObfsPassword"`
	Hysteria2UpMbps       int    `json:"hysteria2UpMbps"`
	Hysteria2DownMbps     int    `json:"hysteria2DownMbps"`
	Hysteria2ServerName   string `json:"hysteria2ServerName"`
	Hysteria2Insecure     uint   `json:"hysteria2Insecure"`
}
