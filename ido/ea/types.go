package ea

type Event struct {
	Link  *Link       `json:"link"`
	App   string      `json:"app"`
	Class string      `json:"class"`
	Scope string      `json:"scope"`
	Item  []string    `json:"item"`
	Op    []string    `json:"op"`
	Msg   chan []byte `json:"msg"`
}

type Action struct {
	Link    *Link             `json:"link"`
	App     string            `json:"app"`
	Owner   string            `json:"owner"`
	Class   string            `json:"class"`
	Scope   string            `json:"scope"`
	Mapping map[string]string `json:"mapping"`
	Op      map[string]string `json:"op"`
	Status  string            `json:"status"`
}

type Link struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
	Event     *Event    `json:"event"`
	Actions   []*Action `json:"actions"`
	Timestamp int64     `json:"timestamp"`
	Status    string    `json:"status"`
	Stop      chan int  `json:"stop"`
}
