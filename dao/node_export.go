package dao

import (
	"errors"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/sirupsen/logrus"
	"trojan-panel/model"
	"trojan-panel/model/constant"
)

// SelectNodeExportList 查询全部节点（导出用，含主表全部字段）
func SelectNodeExportList() ([]model.Node, error) {
	var nodes []model.Node

	buildSelect, values, err := builder.BuildSelect("node", nil, []string{
		"id", "node_server_id", "node_sub_id", "node_type_id", "`name`",
		"node_server_ip", "node_server_grpc_port", "domain", "port", "priority"})
	if err != nil {
		logrus.Errorln(err.Error())
		return nodes, errors.New(constant.SysError)
	}
	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nodes, errors.New(constant.SysError)
	}
	defer rows.Close()

	if err = scanner.Scan(rows, &nodes); err != nil {
		logrus.Errorln(err.Error())
		return nodes, errors.New(constant.SysError)
	}
	return nodes, nil
}

// SelectNodeByNameAndNodeServerId 按名称和服务器查询节点（导入时判断是否重名）
func SelectNodeByNameAndNodeServerId(name string, nodeServerId uint) (*model.Node, error) {
	var node model.Node

	where := map[string]interface{}{
		"`name`":         name,
		"node_server_id": nodeServerId,
	}
	buildSelect, values, err := builder.BuildSelect("node", where, []string{
		"id", "node_server_id", "node_sub_id", "node_type_id", "`name`",
		"node_server_ip", "node_server_grpc_port", "domain", "port", "priority"})
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	defer rows.Close()

	if err = scanner.Scan(rows, &node); err == scanner.ErrEmptyResult {
		return nil, nil
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	return &node, nil
}
