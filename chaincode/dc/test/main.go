package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type User struct {
	Uid     string `json:"uid"`     // 用户唯一标识(数字证书ski)
	Mspid   string `json:"mspid"`   // 用户组织ID
	Role    string `json:"role"`    // 用户角色
	Balance int64  `json:"balance"` // 用户余额
}

func GetFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		//log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

func GetTagName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		//log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		tagName := t.Field(i).Name
		tags := strings.Split(string(t.Field(i).Tag), "\"")
		if len(tags) > 1 {
			tagName = tags[1]
		}
		result = append(result, tagName)
	}
	return result
}

func firstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	rs := []rune(str)
	for i := range rs {
		if i == 0 {
			if rs[i] >= 97 && rs[i] <= 122 {
				rs[i] -= 32
			}
		} else {
			if rs[i] >= 65 && rs[i] <= 90 {
				rs[i] += 32
			}
		}
	}

	return string(rs)
}

func main() {
	str := "mspid=aa, role=ok go, balance=111"

	user := &User{
		Uid:     "123",
		Mspid:   "org1",
		Role:    "member",
		Balance: int64(100),
	}
	kvs := strings.Split(str, ",")

	pp := reflect.ValueOf(user)
	//pt := reflect.TypeOf(user)
	pt := pp.Type()
	var fields []string
	if pt.Kind() == reflect.Ptr {
		pt = pt.Elem()
	}
	if pt.Kind() == reflect.Struct {
		for i := 0; i < pt.NumField(); i++ {
			fields = append(fields, pt.Field(i).Name)
		}
	}

	for _, kv := range kvs {
		kv := strings.Split(kv, "=")
		if len(kv) < 2 {
			return
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		k = strings.ToLower(k)
		fmt.Println(k, v)
		for _, fname := range fields {
			if k == strings.ToLower(fname) {
				k = fname
				break
			}
		}
		field := pp.Elem().FieldByName(k)
		if field.CanSet() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(v)
			case reflect.Int64:
				iv, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return
				}
				field.SetInt(iv)
			}

		} else {
			return
		}

	}

	fmt.Println(GetTagName(user))
	fmt.Println(user)
	mp := make(map[string]string)
	mp["a"] = "A"
	mp["b"] = "B"
	fmt.Println(len(mp))

	fmt.Println(firstToUpper("aaDda"))

	ty := []byte{0x00}
	fmt.Println(bytes.Compare(ty, []byte{0x01}))
	ty = []byte("OK")
	fmt.Println(ty)
	/*
			targetValue := reflect.ValueOf(target)
			switch reflect.TypeOf(target).Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < targetValue.Len(); i++ {
					if targetValue.Index(i).Interface() == obj {
						return true
					}
				}
			case reflect.Map:
				if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
					return true
				}
		    }
	*/
}
