package mathx

// Number is a type constraint that matches all numeric types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// ABS returns the absolute value of a number.
func Abs[T Number](x T) T {
	return AbsDiff(x, 0)
}

// AbsDiff returns the absolute difference between two numbers.
func AbsDiff[T Number](x, y T) T {
	if x < y {
		return y - x
	}
	return x - y
}
