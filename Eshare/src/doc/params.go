package doc

type Document struct {
	ID string `json:"id"` // 文件标识, 可用hash
	//SliceID  string `json:"sliceid"`  // 切片id
	Ext      string   `json:"ext"`      // 类型(后缀)
	Owner    string   `json:"owner"`    // 上传者
	Title    string   `json:"title"`    // 标题
	Name     string   `json:"name"`     // 文件名
	Catalog  string   `json:"catalog"`  // 类目
	Class    string   `json:"class"`    // 类别
	SubClass string   `json:"subclass"` // 子类
	Tag      []string `json:"tag"`      // 标签
	Desc     string   `json:"desc"`     // 描述
	Size     int64    `json:"size"`     // 大小(字节)
	Pagenum  int64    `json:"pagenum"`  // 页数
	Vcnt     int64    `json:"vcnt"`     // 浏览次数
	Dcnt     int64    `json:"dcnt"`     // 下载次数
	Score    int64    `json:"score"`    // 评分(5分制)
	Ccnt     int64    `json:"ccnt"`     // 评论人数
	Rcnt     int64    `json:"rcnt"`     // 评分人数
	Perm     [5]int64 `json:"perm"`     // 权限
	Prenum   int64    `json:"prenum"`   // 预览页数
	Date     int64    `json:"date"`     // 上传时间
	Endorse  string   `json:"endorse"`  // 背书组织
	Block    int64    `json:"block"`    // 区块编号
	Txhash   string   `json:"txhash"`   // 交易哈希
	Status   int64    `json:"status"`   // 状态 (0正常, 1转换失败, 2删除)
}

type Page struct {
	ID      string `json:"id"`      // 文档标识
	Name    string `json:"name"`    // 文档名
	Pagenum int64 `json:"pagenum"` // 文档页数
	Prenum  int64  `json:"prenum"`  // 预览页数
	Number  int64  `json:"number"`  // 页号
	Content []byte `json:"content"` // 页面内容
}
