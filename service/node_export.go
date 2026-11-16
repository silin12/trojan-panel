package service

import (
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/sirupsen/logrus"

	"trojan-panel/dao"
	"trojan-panel/model"
	"trojan-panel/model/constant"
	"trojan-panel/model/dto"
	"trojan-panel/model/vo"
)

// nodeTypeName 节点类型名称，仅用于导出文件可读性
func nodeTypeName(nodeTypeId uint) string {
	switch nodeTypeId {
	case constant.Xray:
		return "xray"
	case constant.TrojanGo:
		return "trojan-go"
	case constant.Hysteria:
		return "hysteria"
	case constant.NaiveProxy:
		return "naiveproxy"
	case constant.Hysteria2:
		return "hysteria2"
	default:
		return fmt.Sprintf("%d", nodeTypeId)
	}
}

// nodeTypeIdByName 导入时按类型名称兜底解析类型 id
func nodeTypeIdByName(name string) uint {
	switch name {
	case "xray":
		return constant.Xray
	case "trojan-go":
		return constant.TrojanGo
	case "hysteria":
		return constant.Hysteria
	case "naiveproxy":
		return constant.NaiveProxy
	case "hysteria2":
		return constant.Hysteria2
	default:
		return 0
	}
}

// ExportNode 导出节点（含各协议子表配置），生成 json 文件并写入文件任务表
func ExportNode(accountId uint, accountUsername string) error {
	nodes, err := dao.SelectNodeExportList()
	if err != nil {
		return err
	}

	nodeExportVos := make([]vo.NodeExportVo, 0, len(nodes))
	for _, node := range nodes {
		nodeExportVo := vo.NodeExportVo{
			Id:                 *node.Id,
			Name:               *node.Name,
			NodeServerIp:       *node.NodeServerIp,
			NodeServerGrpcPort: *node.NodeServerGrpcPort,
			NodeTypeId:         *node.NodeTypeId,
			NodeTypeName:       nodeTypeName(*node.NodeTypeId),
			Domain:             *node.Domain,
			Port:               *node.Port,
			Priority:           *node.Priority,
		}

		switch *node.NodeTypeId {
		case constant.Xray:
			nodeXray, err := dao.SelectNodeXrayById(node.NodeSubId)
			if err != nil {
				logrus.Errorf("ExportNode SelectNodeXrayById err: %v", err)
				continue
			}
			nodeExportVo.XrayProtocol = valueOfString(nodeXray.Protocol)
			nodeExportVo.XrayFlow = valueOfString(nodeXray.XrayFlow)
			nodeExportVo.XraySSMethod = valueOfString(nodeXray.XraySSMethod)
			nodeExportVo.RealityPbk = valueOfString(nodeXray.RealityPbk)
			nodeExportVo.XraySettings = valueOfString(nodeXray.Settings)
			nodeExportVo.XrayStreamSettings = valueOfString(nodeXray.StreamSettings)
			nodeExportVo.XrayTag = valueOfString(nodeXray.Tag)
			nodeExportVo.XraySniffing = valueOfString(nodeXray.Sniffing)
			nodeExportVo.XrayAllocate = valueOfString(nodeXray.Allocate)
		case constant.TrojanGo:
			nodeTrojanGo, err := dao.SelectNodeTrojanGoById(node.NodeSubId)
			if err != nil {
				logrus.Errorf("ExportNode SelectNodeTrojanGoById err: %v", err)
				continue
			}
			nodeExportVo.TrojanGoSni = valueOfString(nodeTrojanGo.Sni)
			nodeExportVo.TrojanGoMuxEnable = valueOfUint(nodeTrojanGo.MuxEnable)
			nodeExportVo.TrojanGoWebsocketEnable = valueOfUint(nodeTrojanGo.WebsocketEnable)
			nodeExportVo.TrojanGoWebsocketPath = valueOfString(nodeTrojanGo.WebsocketPath)
			nodeExportVo.TrojanGoWebsocketHost = valueOfString(nodeTrojanGo.WebsocketHost)
			nodeExportVo.TrojanGoSsEnable = valueOfUint(nodeTrojanGo.SsEnable)
			nodeExportVo.TrojanGoSsMethod = valueOfString(nodeTrojanGo.SsMethod)
			nodeExportVo.TrojanGoSsPassword = valueOfString(nodeTrojanGo.SsPassword)
		case constant.Hysteria:
			nodeHysteria, err := dao.SelectNodeHysteriaById(node.NodeSubId)
			if err != nil {
				logrus.Errorf("ExportNode SelectNodeHysteriaById err: %v", err)
				continue
			}
			nodeExportVo.HysteriaProtocol = valueOfString(nodeHysteria.Protocol)
			nodeExportVo.HysteriaObfs = valueOfString(nodeHysteria.Obfs)
			nodeExportVo.HysteriaUpMbps = valueOfInt(nodeHysteria.UpMbps)
			nodeExportVo.HysteriaDownMbps = valueOfInt(nodeHysteria.DownMbps)
			nodeExportVo.HysteriaServerName = valueOfString(nodeHysteria.ServerName)
			nodeExportVo.HysteriaInsecure = valueOfUint(nodeHysteria.Insecure)
			nodeExportVo.HysteriaFastOpen = valueOfUint(nodeHysteria.FastOpen)
		case constant.Hysteria2:
			nodeHysteria2, err := dao.SelectNodeHysteria2ById(node.NodeSubId)
			if err != nil {
				logrus.Errorf("ExportNode SelectNodeHysteria2ById err: %v", err)
				continue
			}
			nodeExportVo.Hysteria2ObfsPassword = valueOfString(nodeHysteria2.ObfsPassword)
			nodeExportVo.Hysteria2UpMbps = valueOfInt(nodeHysteria2.UpMbps)
			nodeExportVo.Hysteria2DownMbps = valueOfInt(nodeHysteria2.DownMbps)
			nodeExportVo.Hysteria2ServerName = valueOfString(nodeHysteria2.ServerName)
			nodeExportVo.Hysteria2Insecure = valueOfUint(nodeHysteria2.Insecure)
		}
		nodeExportVos = append(nodeExportVos, nodeExportVo)
	}

	return ExportTaskJson(accountId, accountUsername, constant.TaskTypeNodeExport, "nodeExport", nodeExportVos)
}

// ImportNode 导入节点
// cover=1 时同名节点会被覆盖（先删后建），cover=0 时同名节点跳过
func ImportNode(token string, cover uint, file *multipart.FileHeader, accountId uint, accountUsername string) error {
	fileName := file.Filename

	var fileTaskType uint = constant.TaskTypeNodeImport
	var fileTaskStatus = constant.TaskDoing
	fileTask := model.FileTask{
		Name:            &fileName,
		Path:            nil,
		Type:            &fileTaskType,
		Status:          &fileTaskStatus,
		AccountId:       &accountId,
		AccountUsername: &accountUsername,
	}
	fileTaskId, err := dao.CreateFileTask(&fileTask)
	if err != nil {
		return err
	}

	go func(fileTaskId uint) {
		var fail = constant.TaskFail
		var success = constant.TaskSuccess
		task := model.FileTask{
			Id:     &fileTaskId,
			Status: &fail,
		}

		src, err := file.Open()
		if err != nil {
			logrus.Errorf("ImportNode file Open err: %v", err)
			fileUploadError := constant.FileUploadError
			task.ErrMsg = &fileUploadError
			if err = dao.UpdateFileTaskById(&task); err != nil {
				logrus.Errorf("ImportNode UpdateFileTaskById err: %v", err)
			}
			return
		}
		defer src.Close()

		var nodeExportVos []vo.NodeExportVo
		decoder := json.NewDecoder(src)
		if err = decoder.Decode(&nodeExportVos); err != nil {
			logrus.Errorf("ImportNode decoder Decode err: %v", err)
			rowNotEnough := constant.RowNotEnough
			task.ErrMsg = &rowNotEnough
			if err = dao.UpdateFileTaskById(&task); err != nil {
				logrus.Errorf("ImportNode UpdateFileTaskById err: %v", err)
			}
			return
		}
		if len(nodeExportVos) == 0 {
			logrus.Errorf("ImportNode err: %s", constant.RowNotEnough)
			rowNotEnough := constant.RowNotEnough
			task.ErrMsg = &rowNotEnough
			if err = dao.UpdateFileTaskById(&task); err != nil {
				logrus.Errorf("ImportNode UpdateFileTaskById err: %v", err)
			}
			return
		}

		successCount := 0
		failCount := 0
		for _, item := range nodeExportVos {
			if item.Name == "" || item.Domain == "" || item.Port == 0 {
				failCount++
				continue
			}
			if item.Port <= 100 || item.Port >= 30000 {
				logrus.Errorf("ImportNode port out of range, node: %s port: %d", item.Name, item.Port)
				failCount++
				continue
			}

			nodeTypeId := item.NodeTypeId
			if nodeTypeId == 0 {
				nodeTypeId = nodeTypeIdByName(item.NodeTypeName)
			}
			if nodeTypeId == 0 {
				logrus.Errorf("ImportNode unknown node type, node: %s", item.Name)
				failCount++
				continue
			}

			grpcPort := item.NodeServerGrpcPort
			if grpcPort == 0 {
				grpcPort = 8100
			}

			// 先按 ip 找服务器，没有就建一条，方便换面板机时一键搬节点
			nodeServer, err := dao.SelectNodeServer(map[string]interface{}{"ip": item.NodeServerIp})
			if err != nil || nodeServer == nil {
				ip := item.NodeServerIp
				name := item.NodeServerIp
				nodeServer = &model.NodeServer{
					Ip:       &ip,
					Name:     &name,
					GrpcPort: &grpcPort,
				}
				if err = dao.CreateNodeServer(nodeServer); err != nil {
					logrus.Errorf("ImportNode CreateNodeServer err: %v", err)
					failCount++
					continue
				}
				if nodeServer, err = dao.SelectNodeServer(map[string]interface{}{"ip": item.NodeServerIp}); err != nil {
					logrus.Errorf("ImportNode SelectNodeServer err: %v", err)
					failCount++
					continue
				}
			}

			existingNode, err := dao.SelectNodeByNameAndNodeServerId(item.Name, *nodeServer.Id)
			if err != nil {
				failCount++
				continue
			}
			if existingNode != nil {
				if cover == 0 {
					// 不覆盖，跳过同名节点
					successCount++
					continue
				}
				if err = DeleteNodeById(token, existingNode.Id); err != nil {
					logrus.Errorf("ImportNode DeleteNodeById err: %v", err)
					failCount++
					continue
				}
			}

			nodeCreateDto := buildNodeCreateDto(item, nodeTypeId, *nodeServer.Id)
			if err = CreateNode(token, nodeCreateDto); err != nil {
				logrus.Errorf("ImportNode CreateNode err node: %s err: %v", item.Name, err)
				failCount++
				continue
			}
			successCount++
		}

		if successCount == 0 && failCount > 0 {
			rowNotEnough := constant.RowNotEnough
			task.ErrMsg = &rowNotEnough
		} else {
			task.Status = &success
		}
		if err = dao.UpdateFileTaskById(&task); err != nil {
			logrus.Errorf("ImportNode UpdateFileTaskById err: %v", err)
		}
	}(fileTaskId)
	return nil
}

// buildNodeCreateDto 把导出对象还原成创建节点所需的 dto（指针字段全部补齐，避免空指针）
func buildNodeCreateDto(item vo.NodeExportVo, nodeTypeId uint, nodeServerId uint) dto.NodeCreateDto {
	name := item.Name
	domain := item.Domain
	port := item.Port
	priority := item.Priority
	if priority == 0 {
		priority = 100
	}
	nodeServer := nodeServerId
	nodeType := nodeTypeId

	xrayProtocol := item.XrayProtocol
	xrayFlow := item.XrayFlow
	xraySSMethod := item.XraySSMethod
	if xraySSMethod == "" {
		xraySSMethod = "aes-256-gcm"
	}
	realityPbk := item.RealityPbk
	xraySettings := item.XraySettings
	xrayStreamSettings := item.XrayStreamSettings
	xrayTag := item.XrayTag
	xraySniffing := item.XraySniffing
	xrayAllocate := item.XrayAllocate

	trojanGoSni := item.TrojanGoSni
	trojanGoMuxEnable := item.TrojanGoMuxEnable
	if trojanGoMuxEnable == 0 {
		trojanGoMuxEnable = 1
	}
	trojanGoWebsocketEnable := item.TrojanGoWebsocketEnable
	trojanGoWebsocketPath := item.TrojanGoWebsocketPath
	if trojanGoWebsocketPath == "" {
		trojanGoWebsocketPath = "trojan-panel-websocket-path"
	}
	trojanGoWebsocketHost := item.TrojanGoWebsocketHost
	trojanGoSsEnable := item.TrojanGoSsEnable
	trojanGoSsMethod := item.TrojanGoSsMethod
	if trojanGoSsMethod == "" {
		trojanGoSsMethod = "AES-128-GCM"
	}
	trojanGoSsPassword := item.TrojanGoSsPassword

	hysteriaProtocol := item.HysteriaProtocol
	if hysteriaProtocol == "" {
		hysteriaProtocol = "udp"
	}
	hysteriaObfs := item.HysteriaObfs
	hysteriaUpMbps := item.HysteriaUpMbps
	if hysteriaUpMbps <= 0 {
		hysteriaUpMbps = 100
	}
	hysteriaDownMbps := item.HysteriaDownMbps
	if hysteriaDownMbps <= 0 {
		hysteriaDownMbps = 100
	}
	hysteriaServerName := item.HysteriaServerName
	hysteriaInsecure := item.HysteriaInsecure
	hysteriaFastOpen := item.HysteriaFastOpen

	hysteria2ObfsPassword := item.Hysteria2ObfsPassword
	hysteria2UpMbps := item.Hysteria2UpMbps
	if hysteria2UpMbps <= 0 {
		hysteria2UpMbps = 100
	}
	hysteria2DownMbps := item.Hysteria2DownMbps
	if hysteria2DownMbps <= 0 {
		hysteria2DownMbps = 100
	}
	hysteria2ServerName := item.Hysteria2ServerName
	hysteria2Insecure := item.Hysteria2Insecure

	return dto.NodeCreateDto{
		NodeServerId: &nodeServer,
		NodeTypeId:   &nodeType,
		Name:         &name,
		Domain:       &domain,
		Port:         &port,
		Priority:     &priority,

		XrayProtocol:       &xrayProtocol,
		XrayFlow:           &xrayFlow,
		XraySSMethod:       &xraySSMethod,
		RealityPbk:         &realityPbk,
		XraySettings:       &xraySettings,
		XrayStreamSettings: &xrayStreamSettings,
		XrayTag:            &xrayTag,
		XraySniffing:       &xraySniffing,
		XrayAllocate:       &xrayAllocate,

		TrojanGoSni:             &trojanGoSni,
		TrojanGoMuxEnable:       &trojanGoMuxEnable,
		TrojanGoWebsocketEnable: &trojanGoWebsocketEnable,
		TrojanGoWebsocketPath:   &trojanGoWebsocketPath,
		TrojanGoWebsocketHost:   &trojanGoWebsocketHost,
		TrojanGoSsEnable:        &trojanGoSsEnable,
		TrojanGoSsMethod:        &trojanGoSsMethod,
		TrojanGoSsPassword:      &trojanGoSsPassword,

		HysteriaProtocol:   &hysteriaProtocol,
		HysteriaObfs:       &hysteriaObfs,
		HysteriaUpMbps:     &hysteriaUpMbps,
		HysteriaDownMbps:   &hysteriaDownMbps,
		HysteriaServerName: &hysteriaServerName,
		HysteriaInsecure:   &hysteriaInsecure,
		HysteriaFastOpen:   &hysteriaFastOpen,

		Hysteria2ObfsPassword: &hysteria2ObfsPassword,
		Hysteria2UpMbps:       &hysteria2UpMbps,
		Hysteria2DownMbps:     &hysteria2DownMbps,
		Hysteria2ServerName:   &hysteria2ServerName,
		Hysteria2Insecure:     &hysteria2Insecure,
	}
}

func valueOfString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func valueOfUint(v *uint) uint {
	if v == nil {
		return 0
	}
	return *v
}

func valueOfInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}
