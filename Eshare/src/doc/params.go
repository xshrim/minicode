package doc

type Document struct {
	ID       string `json:"id"`       // 文件标识, 可用hash
	SliceID  string `json:"sliceid"`  // 切片id
	Owner    string `json:"owner"`    // 上传者
	Title    string `json:"title"`    // 标题
	Name     string `json:"name"`     // 文件名
	Catalog  string `json:"catalog"`  // 类目
	Class    string `json:"class"`    // 类别
	SubClass string `json:"subclass"` // 子类
	Tag      string `json:"tag"`      // 标签
	Desc     string `json:"desc"`     // 描述
	PageSize int64  `json:"pagesize"` // 页数
	Vcnt     int64  `json:"vcnt"`     // 浏览次数
	Dcnt     int64  `json:"dcnt"`     // 下载次数
	Score    int64  `json:"score"`    // 评分(5分制)
	RaterNum int64  `json:"raternum"` // 评分人数
	Price    int64  `json:"price"`    // 售价
	Date     int64  `json:"date"`     // 上传时间
	Status   int64  `json:"status"`   // 状态 (0正常, 1转换失败, 2删除)
}
