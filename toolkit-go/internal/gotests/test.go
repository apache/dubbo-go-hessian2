package gotests

import "context"

func AbsSum(ctx context.Context, a int, b int) int {
	if a+b < 0 {
		return -(a + b)
	}
	return a + b
}

func Incr(i *int) {
	*i = *i + 1
}

func SomeCal(ctx context.Context, a int, b int) int {
	sum := AbsSum(context.Background(), a, b)
	i := 0
	if i++; b < 0 {
		return AbsSum(context.Background(), a, -b)
	}
	for ; i+a <= sum; Incr(&i) {
		if AbsSum(context.Background(), a, i) == 0 {
			return i
		}
	}
	return i
}
