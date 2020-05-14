package sort

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// 判断数组是否递增
func isIncresing(array []int) bool {
	for i := 0; i < len(array)-1; i++ {
		if array[i] > array[i+1] {
			return false
		}
	}
	return true
}

func timeCost(start time.Time) {
	terminal := time.Since(start)
	fmt.Println(terminal)
}

// 生成随机数组
func RandIntGen(max, num int) []int {
	var res []int
	// fmt.Println(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	for len(res) < num {
		number := rand.Intn(max)
		res = append(res, number)

	}
	return res
}

// 测试冒泡排序
func Test_BubbleSort(t *testing.T) {
	defer timeCost(time.Now())
	array := RandIntGen(100, 10)
	fmt.Println(array)
	BubbleSort(array)
	fmt.Println(array)
	if !isIncresing(array) {
		t.Fail()
	}
}

// 测试选择排序
func Test_SelectSort(t *testing.T) {
	array := RandIntGen(100, 10)
	fmt.Println(array)
	start := time.Now()
	SelectSort(array)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
	fmt.Println(array)
	if !isIncresing(array) {
		t.Fail()
	}
}

// 测试插入排序
func Test_InsertSort(t *testing.T) {
	array := RandIntGen(100, 10)
	fmt.Println(array)
	start := time.Now().UnixNano()
	InsertSort(array)
	end := time.Now().UnixNano()
	fmt.Println(end-start, "ns")
	fmt.Println(array)
	if !isIncresing(array) {
		t.Fail()
	}
}
