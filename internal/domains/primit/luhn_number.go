package primit

import (
	"fmt"
)

type LuhnNumber uint64

var _ fmt.Stringer = (*LuhnNumber)(nil)

func (num LuhnNumber) String() string {
	return fmt.Sprintf("%d", num)
}

func (num LuhnNumber) IsValid() bool {
	return (uint64(num)%10+checksum(uint64(num)/10))%10 == 0
}

func checksum(number uint64) uint64 {
	var luhn uint64

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
