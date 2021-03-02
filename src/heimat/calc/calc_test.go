package calc_test

import (
	"fmt"
	"heimatcli/src/heimat/calc"
)

func ExampleDuration() {
	start := "10:30"
	end := "11:45"

	dur := calc.Duration(start, end)
	fmt.Printf("%s\n", dur)
	// Output: 1h15m0s
}
