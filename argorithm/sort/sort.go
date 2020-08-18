package sort

func BubbleSort(array []int) {
	for t := 0; t < len(array)-1; t++ {
		for i := 0; i < len(array)-1-t; i++ {
			if array[i] > array[i+1] {
				array[i], array[i+1] = array[i+1], array[i]
			}
		}
	}
}

func SelectSort(array []int) {
	for t := 0; t < len(array); t++ {
		midx := t
		for i := t; i < len(array); i++ {
			if array[midx] > array[i] {
				midx = i
			}
		}
		array[t], array[midx] = array[midx], array[t]
	}
}

func InsertSort(array []int) {
	for t := 1; t < len(array); t++ {
		tmp := array[t]
		for i := t - 1; i >= 0; i-- {
			if tmp < array[i] {
				array[i+1] = array[i]
				if i == 0 {
					array[0] = tmp
				}
			} else {
				array[i+1] = tmp
				break
			}
		}
	}
}
