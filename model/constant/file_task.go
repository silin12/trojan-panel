package constant

const (
	TaskFail    = -1
	TaskWait    = 0
	TaskDoing   = 1
	TaskSuccess = 2

	TaskTypeAccountExport    = 1
	TaskTypeNodeServerExport = 2
	TaskTypeAccountImport    = 3
	TaskTypeNodeServerImport = 4
	// TaskTypeNodeExport 节点导出 / TaskTypeNodeImport 节点导入（自定义扩展）
	TaskTypeNodeExport = 5
	TaskTypeNodeImport = 6
)
