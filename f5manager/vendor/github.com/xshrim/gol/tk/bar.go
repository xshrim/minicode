package tk

import (
	"fmt"
	"math"
	"time"
)

type Bar struct {
	total   int
	cur     int
	c       chan int
	verbose bool
}

func NewBar(count ...int) *Bar {
	pb := &Bar{
		total: 100,
		// c: make(chan int),
	}

	if len(count) > 0 {
		pb.total = count[0]
	}

	return pb
}

func (pb *Bar) Start(verbose ...bool) {
	if pb.c != nil {
		return
	}

	pb.cur = 0
	pb.c = make(chan int)

	if len(verbose) > 0 && verbose[0] {
		pb.verbose = true
	}

	startTime := time.Now()

	length := 50

	go func() {
		if pb.verbose {
			fmt.Printf("\r%d/%d [%s] %s", 0, pb.total, fmt.Sprintf("%s %d%%", stringRepeat("-", length), 0), time.Duration(time.Since(startTime).Milliseconds())*time.Millisecond)
		} else {
			fmt.Printf("\r[%s]", stringRepeat("-", length))
		}

		for {
			//time.Sleep(10 * time.Millisecond)
			if value, open := <-pb.c; open {
				if value < 0 {
					value = 0
				} else if value > pb.total {
					value = pb.total
				}

				percent := int(math.Round(float64(length * value / pb.total)))
				//fmt.Printf("\033[%dA\033[K", 1)
				if pb.verbose {
					fmt.Printf("\r%d/%d [%s] %s", value, pb.total, fmt.Sprintf("%s%s %d%%", stringRepeat("=", percent), stringRepeat("-", length-percent), int(math.Round(float64(100*percent/length)))), time.Duration(time.Since(startTime).Milliseconds())*time.Millisecond)
				} else {
					fmt.Printf("\r[%s]", fmt.Sprintf("%s%s", stringRepeat("=", percent), stringRepeat("-", length-percent)))
				}
				if value >= pb.total {
					fmt.Printf("\n")
					break
				}
			} else {
				break
			}
		}
		if pb.c != nil {
			close(pb.c)
			pb.c = nil
		}
	}()
}

func (pb *Bar) Stop() {
	if pb.c != nil {
		close(pb.c)
		pb.c = nil
	}
}

func (pb *Bar) Done() {
	pb.Incr(pb.total)
}

func (pb *Bar) Over() bool {
	return pb.c == nil
}

func (pb *Bar) Wait() {
	for !pb.Over() {
	}
}

func (pb *Bar) Set(percent int) {
	defer func() {
		_ = recover()
	}()

	if pb.c == nil {
		return
	}

	pb.c <- int(math.Round(float64(pb.total * (percent % 101) / 100)))
}

func (pb *Bar) Incr(value ...int) {
	defer func() {
		_ = recover()
	}()

	if pb.c == nil {
		return
	}

	v := 1
	if len(value) > 0 {
		v = value[0]
	}

	pb.cur += v
	if pb.cur > pb.total {
		pb.cur = pb.total
	} else {
		pb.c <- pb.cur
	}
}
