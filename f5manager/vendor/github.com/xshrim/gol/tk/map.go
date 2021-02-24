package tk

import (
	"encoding/json"
	"reflect"
)

type Link struct {
	Value interface{}
	Next  *Link
	Prev  *Link
}

type Map struct {
	m    map[interface{}]interface{}
	head *Link
	tail *Link
}

// build an ordered map
func NewMap(elts ...interface{}) *Map {
	m := &Map{
		m:    make(map[interface{}]interface{}),
		head: nil,
		tail: nil,
	}

	if len(elts) == 1 && reflect.TypeOf(elts[0]).Kind() == reflect.Map {
		v := reflect.ValueOf(elts[0])
		for _, e := range v.MapKeys() {
			m.Set(e.Interface(), v.MapIndex(e).Interface())
		}
	}

	return m
}

func (om *Map) Set(key, value interface{}) {
	om.m[key] = value
	var prev *Link
	cur := om.head
	for {
		if cur == nil {
			break
		}
		i := compareInterface(cur.Value, key)
		if i >= 0 {
			break
		} else {
			prev = cur
			cur = cur.Next
		}
	}

	var node *Link
	if prev == nil {
		node = &Link{Value: key, Next: cur}
		om.head = node
	} else {
		node = &Link{Value: key, Next: cur, Prev: prev}
		prev.Next = node
	}

	if cur == nil {
		om.tail = node
	} else {
		cur.Prev = node
	}
}

func (om *Map) Remove(key interface{}) {
	delete(om.m, key)
	var prev *Link
	cur := om.head
	for {
		if cur == nil {
			return
		}
		i := compareInterface(cur.Value, key)
		if i == 0 {
			break
		} else if i > 0 {
			return
		} else {
			prev = cur
			cur = cur.Next
		}
	}

	if prev == nil {
		om.head = cur.Next
	} else {
		prev.Next = cur.Next
	}

	if cur.Next != nil {
		cur.Next.Prev = prev
	} else {
		om.tail = cur.Prev
	}
}

func (om *Map) Clear() {
	om.m = make(map[interface{}]interface{})
	om.head = nil
	om.tail = nil
}

func (om *Map) Keys(reverse ...bool) []interface{} {
	keys := []interface{}{}
	if len(reverse) > 0 && reverse[0] {
		for p := om.tail; p != nil; p = p.Prev {
			keys = append(keys, p.Value)
		}
	} else {
		for p := om.head; p != nil; p = p.Next {
			keys = append(keys, p.Value)
		}
	}

	return keys
}

func (om *Map) Values() []interface{} {
	values := []interface{}{}
	for _, v := range om.m {
		values = append(values, v)
	}

	return values
}

func (om *Map) Contain(k interface{}) bool {
	_, ok := om.m[k]
	return ok
}

func (om *Map) Get(k interface{}) interface{} {
	return om.m[k]
}

func (om *Map) Len() int {
	return len(om.m)
}

func (om *Map) Pop(dequeue ...bool) (interface{}, interface{}) {
	var p *Link
	if len(dequeue) > 0 && dequeue[0] {
		p = om.tail
	} else {
		p = om.head
	}

	if p == nil {
		return nil, nil
	}

	key := p.Value
	value := om.m[key]

	om.Remove(key)

	return key, value
}

func (om *Map) Iter(reverse ...bool) <-chan [2]interface{} {
	r := false
	if len(reverse) > 0 && reverse[0] {
		r = true
	}

	c := make(chan [2]interface{})
	p := om.head
	if r {
		p = om.tail
	}

	go func() {
		for p != nil {
			c <- [2]interface{}{p.Value, om.m[p.Value]}
			if r {
				p = p.Prev
			} else {
				p = p.Next
			}
		}
		close(c)
	}()

	return c
}

func (om *Map) Clone() *Map {
	nm := NewMap()
	for kv := range om.Iter() {
		nm.Set(kv[0], kv[1])
	}

	return nm
}

func (om *Map) Equal(other *Map) bool {
	if om.Len() != other.Len() {
		return false
	}

	return reflect.DeepEqual(om.m, other.m)

	// for k, v := range om.m {
	// 	if mv, ok := mp.m[k]; !ok || compareInterface(v, mv) != 0 {
	// 		return false
	// 	}
	// }

	// return true
}

func (om *Map) Join(mp *Map) {
	for kv := range mp.Iter() {
		om.Set(kv[0], kv[1])
	}
}

func (om *Map) Unmarshal(js string) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(js), &m); err == nil {
		om.Clear()
		for k, v := range m {
			om.Set(k, v)
		}
	}
}

func (om *Map) Marshal() string {
	m := make(map[string]interface{})
	for k, v := range om.m {
		if key, ok := k.(string); ok {
			m[key] = v
		}
	}
	js, _ := json.Marshal(m)
	return string(js)
}
