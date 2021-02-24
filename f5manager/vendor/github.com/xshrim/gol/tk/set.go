package tk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Set map[interface{}]struct{}

// An orderedPair represents a 2-tuple of values.
// type orderedPair struct {
// 	First  interface{}
// 	Second interface{}
// }

type iterator struct {
	C    <-chan interface{}
	stop chan struct{}
}

func (i *iterator) Stop() {
	defer func() {
		_ = recover()
	}()

	close(i.stop)

	for range i.C {
	}
}

func newIterator() (*iterator, chan<- interface{}, <-chan struct{}) {
	itemChan := make(chan interface{})
	stopChan := make(chan struct{})
	return &iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}

// build a set
func NewSet(elts ...interface{}) Set {
	s := make(Set)
	if len(elts) == 1 {
		tp := reflect.TypeOf(elts[0]).Kind()
		if tp == reflect.Slice || tp == reflect.Array {
			v := reflect.ValueOf(elts[0])
			for i := 0; i < v.Len(); i++ {
				s.Add(v.Index(i).Interface())
			}
		}
	} else {
		for _, e := range elts {
			s.Add(e)
		}
	}
	return s
}

// func (pair *orderedPair) Equal(other orderedPair) bool {
// 	if pair.First == other.First &&
// 		pair.Second == other.Second {
// 		return true
// 	}

// 	return false
// }

func (s Set) Add(i interface{}) bool {
	_, found := s[i]
	if found {
		return false
	}

	s[i] = struct{}{}
	return true
}

func (s Set) Contain(i ...interface{}) bool {
	for _, val := range i {
		if _, ok := s[val]; !ok {
			return false
		}
	}
	return true
}

func (s Set) IsSubs(other Set) bool {
	if s.Len() > other.Len() {
		return false
	}
	for elem := range s {
		if !other.Contain(elem) {
			return false
		}
	}
	return true
}

func (s Set) IsProperSubs(other Set) bool {
	return s.IsSubs(other) && !s.Equal(other)
}

func (s Set) IsSupers(other Set) bool {
	return other.IsSubs(s)
}

func (s Set) IsProperSupers(other Set) bool {
	return s.IsSupers(other) && !s.Equal(other)
}

func (s Set) Union(other Set) Set {
	unionedSet := NewSet()

	for elem := range s {
		unionedSet.Add(elem)
	}
	for elem := range other {
		unionedSet.Add(elem)
	}
	return unionedSet
}

func (s Set) Intersect(other Set) Set {
	intersection := NewSet()
	if s.Len() < other.Len() {
		for elem := range s {
			if other.Contain(elem) {
				intersection.Add(elem)
			}
		}
	} else {
		for elem := range other {
			if s.Contain(elem) {
				intersection.Add(elem)
			}
		}
	}
	return intersection
}

func (s Set) Diff(other Set) Set {
	difference := NewSet()
	for elem := range s {
		if !other.Contain(elem) {
			difference.Add(elem)
		}
	}
	return difference
}

func (s Set) SymmetricDiff(other Set) Set {
	aDiff := s.Diff(other)
	bDiff := other.Diff(s)
	return aDiff.Union(bDiff)
}

func (s Set) Clear() {
	for elem := range s {
		delete(s, elem)
	}
}

func (s Set) Remove(i interface{}) {
	delete(s, i)
}

func (s Set) Sort() {
	sl := s.ToSlice()
	quickSortInterface(sl, 0, len(sl)-1)
	s.Clear()
	for _, elem := range sl {
		s.Add(elem)
	}
}

func (s Set) Len() int {
	return len(s)
}

func (s Set) Each(fn func(interface{}) bool) {
	for elem := range s {
		if fn(elem) {
			break
		}
	}
}

func (s Set) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for elem := range s {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (s Set) iterator() *iterator {
	iterator, ch, stopCh := newIterator()

	go func() {
	L:
		for elem := range s {
			select {
			case <-stopCh:
				break L
			case ch <- elem:
			}
		}
		close(ch)
	}()

	return iterator
}

func (s Set) Equal(other Set) bool {
	if s.Len() != other.Len() {
		return false
	}
	for elem := range s {
		if !other.Contain(elem) {
			return false
		}
	}
	return true
}

func (s Set) Clone() Set {
	cs := NewSet()
	for elem := range s {
		cs.Add(elem)
	}
	return cs
}

func (s Set) String() string {
	items := make([]string, 0, len(s))

	for elem := range s {
		items = append(items, fmt.Sprintf("%v", elem))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}

// func (pair orderedPair) String() string {
// 	return fmt.Sprintf("(%v, %v)", pair.First, pair.Second)
// }

func (s Set) Pop() interface{} {
	for item := range s {
		delete(s, item)
		return item
	}
	return nil
}

func (s Set) PowerSet() Set {
	ps := NewSet()
	ns := NewSet()
	ps.Add(&ns)

	for es := range s {
		u := NewSet()
		j := ps.Iter()
		for er := range j {
			p := NewSet()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(Set)
				for ek := range k {
					p.Add(ek)
				}
			} else {
				p.Add(er)
			}
			p.Add(es)
			u.Add(&p)
		}

		ps = ps.Union(u)
	}

	return ps
}

// func (s Set) CartesianProduct(other Set) Set {
// 	cartProduct := NewSet()

// 	for i := range s {
// 		for j := range other {
// 			elem := orderedPair{First: i, Second: j}
// 			cartProduct.Add(elem)
// 		}
// 	}

// 	return cartProduct
// }

func (s Set) ToSlice() []interface{} {
	keys := make([]interface{}, 0, s.Len())
	for elem := range s {
		keys = append(keys, elem)
	}

	return keys
}

// Marshal creates a JSON array from the set, it marshals all elements
func (s Set) Marshal() ([]byte, error) {
	items := make([]string, 0, s.Len())

	for elem := range s {
		b, err := json.Marshal(elem)
		if err != nil {
			return nil, err
		}

		items = append(items, string(b))
	}

	return []byte(fmt.Sprintf("[%s]", strings.Join(items, ","))), nil
}

// Unmarshal recreates a set from a JSON array, it only decodes primitive types. Numbers are decoded as json.Number.
func (s Set) Unmarshal(b []byte) error {
	var i []interface{}

	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	err := d.Decode(&i)
	if err != nil {
		return err
	}

	for _, v := range i {
		switch t := v.(type) {
		case []interface{}, map[string]interface{}:
			continue
		default:
			s.Add(t)
		}
	}

	return nil
}
