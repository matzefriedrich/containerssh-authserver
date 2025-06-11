package types

import "fmt"

type Quota struct {
	fmt.Stringer
	size       int
	unitString string
}

func Gb(size int) Quota {
	return Quota{
		size:       size,
		unitString: "G",
	}
}

func Mb(size int) *Quota {
	return &Quota{
		size:       size,
		unitString: "M",
	}
}

func (q Quota) String() string {
	return fmt.Sprintf("%d%s", q.size, q.unitString)
}
