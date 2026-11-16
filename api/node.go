package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"trojan-panel/model/constant"
	"trojan-panel/model/dto"
	"trojan-panel/model/vo"
	"trojan-panel/service"
	"trojan-panel/util"
)

func SelectNodeById(c *gin.Context) {
	var nodeRequireIdDto dto.RequiredIdDto
	_ = c.ShouldBindQuery(&nodeRequireIdDto)
	if err := validate.Struct(&nodeRequireIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	nodeOneVo, err := service.SelectNodeById(nodeRequireIdDto.Id)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nodeOneVo, c)
}

func SelectNodeInfo(c *gin.Context) {
	var nodeRequireIdDto dto.RequiredIdDto
	_ = c.ShouldBindQuery(&nodeRequireIdDto)
	if err := validate.Struct(&nodeRequireIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	nodeOneVo, err := service.SelectNodeInfo(nodeRequireIdDto.Id, c)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nodeOneVo, c)
}

func CreateNode(c *gin.Context) {
	var nodeCreateDto dto.NodeCreateDto
	_ = c.ShouldBindJSON(&nodeCreateDto)
	if err := validate.Struct(&nodeCreateDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.CreateNode(util.GetToken(c), nodeCreateDto); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func SelectNodePage(c *gin.Context) {
	var nodePageDto dto.NodePageDto
	_ = c.ShouldBindQuery(&nodePageDto)
	if err := validate.Struct(&nodePageDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	nodePageVo, err := service.SelectNodePage(nodePageDto.Name,
		nodePageDto.NodeServerId,
		nodePageDto.PageNum,
		nodePageDto.PageSize, c)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nodePageVo, c)
}

func DeleteNodeById(c *gin.Context) {
	var nodeRequireIdDto dto.RequiredIdDto
	_ = c.ShouldBindJSON(&nodeRequireIdDto)
	if err := validate.Struct(&nodeRequireIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.DeleteNodeById(util.GetToken(c), nodeRequireIdDto.Id); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func UpdateNodeById(c *gin.Context) {
	var nodeUpdateDto dto.NodeUpdateDto
	_ = c.ShouldBindJSON(&nodeUpdateDto)
	if err := validate.Struct(&nodeUpdateDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	if err := service.UpdateNodeById(util.GetToken(c), &nodeUpdateDto); err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nil, c)
}

func NodeQRCode(c *gin.Context) {
	var requiredIdDto dto.RequiredIdDto
	_ = c.ShouldBindJSON(&requiredIdDto)
	if err := validate.Struct(&requiredIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	accountInfo, err := service.GetAccountInfo(c)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	qrCode, err := service.NodeQRCode(&accountInfo.Id, &accountInfo.Username, requiredIdDto.Id)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(qrCode, c)
}

func NodeURL(c *gin.Context) {
	var requiredIdDto dto.RequiredIdDto
	_ = c.ShouldBindJSON(&requiredIdDto)
	if err := validate.Struct(&requiredIdDto); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		fmt.Println(err.Error())
		return
	}
	accountInfo, err := service.GetAccountInfo(c)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	url, _, err := service.NodeURL(&accountInfo.Id, &accountInfo.Username, requiredIdDto.Id)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(url, c)
}

func NodeDefault(c *gin.Context) {
	nodeDefault, err := service.NodeDefault()
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(nodeDefault, c)
}

// ExportNode 导出节点
func ExportNode(c *gin.Context) {
	accountVo := service.GetCurrentAccount(c)
	if err := service.ExportNode(accountVo.Id, accountVo.Username); err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	vo.Success(nil, c)
}

// ImportNode 导入节点
func ImportNode(c *gin.Context) {
	coverStr, b := c.GetPostForm("cover")
	if !b {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	cover, err := strconv.ParseUint(coverStr, 10, 32)
	if err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	// 文件大小 10MB
	if file.Size > 1024*1024*10 {
		vo.Fail(constant.FileSizeTooBig, c)
		return
	}
	// 文件后缀.json
	if !strings.HasSuffix(file.Filename, ".json") {
		vo.Fail(constant.FileFormatError, c)
		return
	}
	account := service.GetCurrentAccount(c)
	if err := service.ImportNode(util.GetToken(c), uint(cover), file, account.Id, account.Username); err != nil {
		vo.Fail(constant.SysError, c)
		return
	}
	vo.Success(nil, c)
}
