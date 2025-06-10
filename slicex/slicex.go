package slicex

// Rotate rotates a slice by the specified number of elements.
// A positive num rotates the slice so that the last num elements move to the beginning.
// A negative num rotates the slice so that the last |num| elements move to the beginning.
func Rotate[T any](slice []T, num int) []T {
	if len(slice) == 0 {
		return slice
	}

	// Get the length of the slice
	length := len(slice)

	var index int
	if num < 0 {
		index = length - ((-num) % length)
	} else {
		index = num % length
	}
	// No rotation needed
	if index == 0 {
		return slice
	}

	// Move the last 'num' elements to the front
	return append(slice[length-index:], slice[:length-index]...)
}
