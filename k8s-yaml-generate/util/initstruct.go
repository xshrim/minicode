package util

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

//初始化一个struct,递归初始化内部属性
//过滤掉不可导出属性，过滤不能json序列化的属性
func InitializeStruct(object interface{}) reflect.Value {
	t := reflect.TypeOf(object)
	v := reflect.New(t)
	maps := make(map[string]string)
	maps["kind"] = t.Name()
	value := v.Elem()
	SetValue(t, value, maps)
	return value
}

//设置属性值
//string 默认赋值 string-value
//int 默认赋值 44
func SetValue(t reflect.Type, v reflect.Value, paramMap map[string]string) {
	if !v.CanSet() {
		fmt.Println("can not set", t, v)
		return
	} else {
		//fmt.Println("v.Interface():",v)
	}

	switch t.Kind() {
	case reflect.String:
		fmt.Println("string--:", t, v)
		if "kind" == t.String() {
			v.SetString(paramMap["kind"])
		} else {
			v.SetString("string-value")
		}
		return
	case reflect.Int8:
		v.SetInt(44)
		return
	case reflect.Uint:
		v.SetUint(44)
		return
	case reflect.Uint8:
		v.SetUint(44)
		return
	case reflect.Int64:
		//fmt.Println(v)
		v.SetInt(44)
		return
	}

	var ft reflect.StructField
	var fv reflect.Value
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic-:", err, t, v, ft.Type.Kind(), ft, fv)
			panic(err)
		}
	}()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft = t.Field(i)
		//json Marshal bug '-'
		if strings.Contains(fmt.Sprint(ft.Tag), `json:"-"`) {
			fmt.Println("tag:", ft.Tag)
			continue
		}
		switch ft.Type.Kind() {
		case reflect.String:
			if f.CanSet() {
				//fmt.Println("string:",ft.Tag.Get("json"))
				if strings.Contains(ft.Tag.Get("json"), "kind,") {
					f.SetString(paramMap["kind"])
				} else {
					f.SetString("string-value")
				}

			}
		case reflect.Int:
			f.SetInt(44)
		case reflect.Int8:
			f.SetInt(44)
		case reflect.Int16:
			f.SetInt(44)
		case reflect.Int32:
			f.SetInt(44)
		case reflect.Int64:
			//fmt.Println("int64:",f,ft.Type)
			if fmt.Sprint(ft.Type) == "intstr.Type" {
				f.SetInt(1)
			} else {
				f.SetInt(44)
			}

		case reflect.Bool:
			f.SetBool(false)
		case reflect.Map:
			fmt.Println("map:", ft, ft.Type.Key(), ft.Type.Elem())
			maps := reflect.MakeMapWithSize(ft.Type, 1)
			key := reflect.New(ft.Type.Key())
			SetValue(key.Elem().Type(), key.Elem(), paramMap)

			value := reflect.New(ft.Type.Elem())
			if value.Elem().Type().Kind() == reflect.Slice {
				sliceValue := reflect.MakeSlice(value.Elem().Type(), 1, 1)
				v := reflect.New(sliceValue.Type().Elem())
				SetValue(v.Elem().Type(), v.Elem(), paramMap)
				sliceValue.Index(0).Set(v.Elem())
				//fmt.Println("slice1:", sliceValue)
				if value.Elem().Type().Kind() != reflect.Ptr {
				} else {
					fmt.Println("slice prt:", t)
				}
			} else {
				SetValue(value.Elem().Type(), value.Elem(), paramMap)
				fmt.Println(key.Elem(), value.Elem())
				maps.SetMapIndex(key.Elem(), value.Elem())
				f.Set(maps)
			}

		case reflect.Slice:
			sliceValue := reflect.MakeSlice(ft.Type, 1, 1)
			if ft.Type.Elem().Kind() != reflect.Ptr {
				v := reflect.New(ft.Type.Elem())
				fmt.Println("slice1:", ft.Type.Elem())
				SetValue(ft.Type.Elem(), v.Elem(), paramMap)
				sliceValue.Index(0).Set(v.Elem())
				f.Set(sliceValue)
			} else {
				fmt.Println("slice prt:", ft)
			}

		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			if f.Type().String() == "time.Time" {
				f.Set(reflect.ValueOf(time.Now()))
				continue
			}
			//fmt.Println("struct--:",ft.Type,ft.Tag,f)
			SetValue(ft.Type, f, paramMap)
		case reflect.Ptr:
			fv = reflect.New(ft.Type.Elem())
			//fmt.Println("Ptr Type:", ft.Type.Elem(), ft.Type.Elem().Kind())
			switch ft.Type.Elem().Kind() {
			case reflect.String:
				stringValue := "string-value"
				if ft.Tag.Get("json") == "kind" {
					stringValue = paramMap["kind"]
				}
				//fmt.Println("Ptr string:",ft.Type)
				fv = reflect.NewAt(ft.Type.Elem(), unsafe.Pointer(&stringValue))
				f.Set(fv)
				break
			case reflect.Int:
				var intValue int = 44
				fv = reflect.NewAt(ft.Type.Elem(), unsafe.Pointer(&intValue))
				f.Set(fv)
				break
			case reflect.Int8:
				var intValue int8 = 44
				fv = reflect.NewAt(ft.Type.Elem(), unsafe.Pointer(&intValue))
				f.Set(fv)
				break
			case reflect.Int16:
				var intValue int16 = 44
				fv = reflect.NewAt(ft.Type.Elem(), unsafe.Pointer(&intValue))
				f.Set(fv)
				break
			case reflect.Int32:
				var intValue int32 = 44
				fv = reflect.NewAt(ft.Type.Elem(), unsafe.Pointer(&intValue))
				fmt.Println("*int32")
				f.Set(fv)
				break
			case reflect.Int64:
				var intValue int64 = 44
				fv = reflect.NewAt(ft.Type.Elem(), unsafe.Pointer(&intValue))
				f.Set(fv)
				break
			case reflect.Bool:
				boolValue := false
				fv = reflect.NewAt(ft.Type.Elem(), unsafe.Pointer(&boolValue))
				f.Set(fv)
				break
			case reflect.Map:
				f.Set(reflect.MakeMap(ft.Type))
				break
			case reflect.Slice:
				f.Set(reflect.MakeSlice(ft.Type, 0, 0))
				break
			case reflect.Chan:
				f.Set(reflect.MakeChan(ft.Type, 0))
				break
			case reflect.Struct:
				//fmt.Println("f.ptr.struct:",f.Kind(),f.Type())
				SetValue(fv.Elem().Type(), fv.Elem(), paramMap)
				//fmt.Println("Prt fv struct:", f,fv)
				f.Set(fv)
				break
			}

		}
	}
}
