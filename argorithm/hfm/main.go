package main

import (
	"fmt"
	"sort"
)

type Node struct {
	Key   string // Key值非空表示叶节点
	Value int    // 权值(key出现次数)
	Left  *Node
	Right *Node
}

/*
模拟实现哈夫曼压缩算法
*/

// 模拟显示字符串的二进制表示
func StringToBin(s string) (binString string) {
	for _, c := range s {
		binString = fmt.Sprintf("%s%b", binString, c)
	}
	return
}

// 字符排序
func SortNode(charmap map[rune]int) []*Node {

	var nodes []*Node
	//根据map构建slice
	for k, v := range charmap {
		node := &Node{
			Key:   string(k),
			Value: v,
			Left:  nil,
			Right: nil,
		}
		nodes = append(nodes, node)
	}

	// 排序
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Value > nodes[j].Value // 降序, 降序有利于结合插入排序构建哈夫曼树
	})

	return nodes
}

// 根据原始字符串生成按字符频率排序的切片
func CreateSlice(str string) []*Node {

	charmap := make(map[rune]int)

	for _, c := range str {
		if _, ok := charmap[c]; ok {
			charmap[c] += 1
		} else {
			charmap[c] = 1
		}
	}

	return SortNode(charmap) // 生成降序排列的切片
}

// 生成树
func CreateTree(str string) *Node {

	nodes := CreateSlice(str) // 生成字符有序切片

	fmt.Println("有序字符切片: ")
	for _, node := range nodes {
		fmt.Println(fmt.Sprintf("%4s: %d", node.Key, node.Value))
	}

	// 生成树
	if len(nodes) == 0 {
		return nil
	}

	if len(nodes) == 1 {
		return nodes[0]
	}

	for len(nodes) > 1 { // 每次循环都将切片最后两个元素删除, 将其合并为一个元素并按元素大小插入到切片合适位置
		length := len(nodes)
		pnode := &Node{
			Key:   nodes[length-2].Key + nodes[length-1].Key,
			Value: nodes[length-2].Value + nodes[length-1].Value,
			Left:  nodes[length-2],
			Right: nodes[length-1],
		}

		if length == 2 {
			nodes[0] = pnode
		} else {
			for i := length - 3; i >= 0; i-- {
				if nodes[i].Value < pnode.Value {
					nodes[i+1] = nodes[i]
					if i == 0 {
						nodes[i] = pnode
					}
				} else {
					nodes[i+1] = pnode
					break
				}
			}
		}

		nodes = nodes[:length-1]

		/*
					fmt.Println("===================================")
					for _, node := range nodes {
						fmt.Println(node.Key, node.Value)
			        }
		*/
	}
	return nodes[0]
}

// 遍历树, 生成编码对照表
func TravelTree(cmap map[string]string, tree *Node, str string) { // 深度优先遍历

	if tree.Left == nil && tree.Right == nil {
		//fmt.Println(tree.Key, str)
		if str == "" { // 单字符的字符串编码
			cmap[tree.Key] = "0"
		} else {
			cmap[tree.Key] = str
		}
	}
	if tree.Left != nil {
		TravelTree(cmap, tree.Left, str+"0")
	}
	if tree.Right != nil {
		TravelTree(cmap, tree.Right, str+"1")
	}
}

// 根据哈夫曼树对原始字符串进行编码
func Encode(str string, tree *Node) string {
	var res string

	if tree == nil {
		return str
	}

	cmap := make(map[string]string) // 保存字符与编码映射关系（编码对照表）

	TravelTree(cmap, tree, "") // 遍历哈夫曼树路径获取编码对照表

	fmt.Println("编码对照表:")
	for k, v := range cmap {
		fmt.Println(fmt.Sprintf("%4s: %s", k, v))
	}

	for _, c := range str { // 根据编码对照表生成编码字符串
		res += cmap[string(c)]
	}

	return res
}

// 根据哈夫曼树对编码字符串进行解码
func Decode(str string, tree *Node) string {
	var res string

	if tree == nil {
		return str
	}

	for i := 0; i < len(str); i++ {
		subtree := tree
		if subtree.Left == nil && subtree.Right == nil { // 哈夫曼树只有一层时
			res += subtree.Key
			continue
		}

		for {
			if string(str[i]) == "0" {
				subtree = subtree.Left
			} else {
				subtree = subtree.Right
			}

			if subtree.Left == nil && subtree.Right == nil { // 达到叶节点表示一段编码字符串可以解码为该叶节点的key(即原始字符)
				res += subtree.Key
				break
			}
			i++
		}
	}

	return res
}

func main() {
	originStr := "AABBBCCCCDEEFFF字符串哈字符串对对对"
	//originStr = "ab"
	/*
			for _, node := range nodes {
				fmt.Println(node.Key, node.Value)
		    }
	*/
	fmt.Println(fmt.Sprintf("原始字符串: %s", originStr))
	fmt.Println(originStr)

	fmt.Println(StringToBin(originStr))

	tree := CreateTree(originStr) // 生成哈夫曼树

	encodedStr := Encode(originStr, tree) // 根据哈夫曼树编码原始字符串

	fmt.Println("编码字符串:")
	fmt.Println(encodedStr)

	decodedStr := Decode(encodedStr, tree) // 根据哈夫曼树解码编码字符串

	fmt.Println(fmt.Sprintf("解码字符串: %s", decodedStr))
	fmt.Println(StringToBin(decodedStr))

	fmt.Println(fmt.Sprintf("压缩率: %.4f", float64(len(encodedStr))/float64(len(StringToBin(originStr)))))
}
