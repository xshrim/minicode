package ea

import (
	"encoding/json"
	"fmt"
	"ido/consts"
	"time"
)

var Links map[string]*Link

func init() {
	Links = make(map[string]*Link)
}

// 事件触发器生成信息发送到msg通道
func (e *Event) Run(stop chan int) {
	go e.Notify()
	ticker := time.NewTicker(time.Duration(time.Second * consts.EVENTINTERVAL))
	defer close(e.Msg)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			// TODO 生成真实数据
			e.Msg <- []byte("Event " + e.App)
		}
	}
}

// 事件触发器从msg通道取出信息发送到各动作
func (e *Event) Notify() {
	for info := range e.Msg {
		for _, action := range e.Link.Actions {
			if action.Status == "run" {
				action.Do(info)
			}
		}
	}
}

// 动作处理接收到的信息
func (a *Action) Do(data []byte) {
	// TODO 处理真实数据
	fmt.Println("Action " + a.App + " Received Message From " + string(data))
}

func New(body []byte) (*Link, error) {
	var link Link

	if err := json.Unmarshal(body, &link); err != nil {
		return nil, err
	}

	// TODO link id 也可以根据部分字段进行哈希生成
	link.Event.Link = &link
	link.Event.Msg = make(chan []byte, consts.EVENTLENGTH)

	for _, action := range link.Actions {
		action.Link = &link
	}
	key := link.Owner + consts.STRCONNECTOR + link.Id
	if _, ok := Links[key]; ok {
		return nil, fmt.Errorf("Link exists already")
	}

	Links[key] = &link

	return &link, nil
}

func (l *Link) RunLink() {
	l.Status = "run"
	go l.Event.Run(l.Stop)
}

func (l *Link) AddAction(actions ...*Action) {
	l.Actions = append(l.Actions, actions...)
}

func (l *Link) DelAction(idx int) error {
	length := len(l.Actions)
	if idx >= length || idx < 0 {
		return fmt.Errorf("index out of range")
	} else if idx == length-1 {
		l.Actions = l.Actions[:idx]
	} else {
		l.Actions = append(l.Actions[:idx], l.Actions[idx+1:]...)
	}
	return nil
}

// TODO
func (l *Link) SetEvent(attr, val string) error {
	return nil
}

// TODO
func (l *Link) SetAction(idx int, attr, val string) error {
	return nil
}

func (l *Link) StopLink() {
	l.Stop <- 1
	l.Status = "stop"
}
