package core

type Pool struct {
	Name   string
	Weight int
}

type Pools []Pool

func (ps Pools) Len() int {
	return len(ps)
}

func (ps Pools) Less(i, j int) bool {
	return ps[i].Weight > ps[j].Weight
}

func (ps Pools) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps Pools) Accumulate(i int) int {
	sum := 0
	for idx, p := range ps {
		if idx <= i {
			sum += p.Weight
		} else {
			break
		}
	}

	return sum
}
