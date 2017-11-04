package collections

import (
	"fmt"
	"strings"
)

const (
	ringBufferMinimumGrow     = 4
	ringBufferShrinkThreshold = 32
	ringBufferGrowFactor      = 200
	ringBufferDefaultCapacity = 4
)

var (
	emptyArray = make([]interface{}, 0)
)

// NewRingBuffer creates a new, empty, RingBuffer.
func NewRingBuffer() *RingBuffer {
	return &RingBuffer{
		array: make([]interface{}, ringBufferDefaultCapacity),
		head:  0,
		tail:  0,
		size:  0,
	}
}

// NewRingBufferWithCapacity creates a new RingBuffer pre-allocated with the given capacity.
func NewRingBufferWithCapacity(capacity int) *RingBuffer {
	return &RingBuffer{
		array: make([]interface{}, capacity),
		head:  0,
		tail:  0,
		size:  0,
	}
}

// NewRingBufferFromSlice createsa  ring buffer out of a slice.
func NewRingBufferFromSlice(values []interface{}) *RingBuffer {
	return &RingBuffer{
		array: values,
		head:  0,
		tail:  len(values) - 1,
		size:  len(values),
	}
}

// RingBuffer is a fifo buffer that is backed by a pre-allocated array, instead of allocating
// a whole new node object for each element (which saves GC churn).
// Enqueue can be O(n), Dequeue can be O(1).
type RingBuffer struct {
	array []interface{}
	head  int
	tail  int
	size  int
}

// Len returns the length of the ring buffer (as it is currently populated).
// Actual memory footprint may be different.
func (rb *RingBuffer) Len() (len int) {
	return rb.size
}

// TotalLen returns the total size of the ring bufffer, including empty elements.
func (rb *RingBuffer) TotalLen() int {
	return len(rb.array)
}

// Clear removes all objects from the RingBuffer.
func (rb *RingBuffer) Clear() {

	if rb.head < rb.tail {
		arrayClear(rb.array, rb.head, rb.size)
	} else {
		arrayClear(rb.array, rb.head, len(rb.array)-rb.head)
		arrayClear(rb.array, 0, rb.tail)
	}

	rb.head = 0
	rb.tail = 0
	rb.size = 0
}

// Enqueue adds an element to the "back" of the RingBuffer.
func (rb *RingBuffer) Enqueue(object interface{}) {
	if rb.size == len(rb.array) {
		newCapacity := int(len(rb.array) * int(ringBufferGrowFactor/100))
		if newCapacity < (len(rb.array) + ringBufferMinimumGrow) {
			newCapacity = len(rb.array) + ringBufferMinimumGrow
		}
		rb.setCapacity(newCapacity)
	}

	rb.array[rb.tail] = object
	rb.tail = (rb.tail + 1) % len(rb.array)
	rb.size++
}

// Dequeue removes the first (oldest) element from the RingBuffer.
func (rb *RingBuffer) Dequeue() interface{} {
	if rb.size == 0 {
		return nil
	}

	removed := rb.array[rb.head]
	rb.head = (rb.head + 1) % len(rb.array)
	rb.size--

	return removed
}

// Peek returns but does not remove the first element.
func (rb *RingBuffer) Peek() interface{} {
	if rb.size == 0 {
		return nil
	}
	return rb.array[rb.head]
}

// PeekBack returns but does not remove the last element.
func (rb *RingBuffer) PeekBack() interface{} {
	if rb.size == 0 {
		return nil
	}
	if rb.tail == 0 {
		return rb.array[len(rb.array)-1]
	}
	return rb.array[rb.tail-1]
}

func (rb *RingBuffer) setCapacity(capacity int) {
	newArray := make([]interface{}, capacity)
	if rb.size > 0 {
		if rb.head < rb.tail {
			arrayCopy(rb.array, rb.head, newArray, 0, rb.size)
		} else {
			arrayCopy(rb.array, rb.head, newArray, 0, len(rb.array)-rb.head)
			arrayCopy(rb.array, 0, newArray, len(rb.array)-rb.head, rb.tail)
		}
	}
	rb.array = newArray
	rb.head = 0
	if rb.size == capacity {
		rb.tail = 0
	} else {
		rb.tail = rb.size
	}
}

// TrimExcess resizes the buffer to better fit the contents.
func (rb *RingBuffer) TrimExcess() {
	threshold := float64(len(rb.array)) * 0.9
	if rb.size < int(threshold) {
		rb.setCapacity(rb.size)
	}
}

// AsSlice returns the ring buffer, in order, as a slice.
func (rb *RingBuffer) AsSlice() []interface{} {
	newArray := make([]interface{}, rb.size)

	if rb.size == 0 {
		return newArray
	}

	if rb.head < rb.tail {
		arrayCopy(rb.array, rb.head, newArray, 0, rb.size)
	} else {
		arrayCopy(rb.array, rb.head, newArray, 0, len(rb.array)-rb.head)
		arrayCopy(rb.array, 0, newArray, len(rb.array)-rb.head, rb.tail)
	}

	return newArray
}

// Each calls the consumer for each element in the buffer.
func (rb *RingBuffer) Each(consumer func(value interface{})) {
	if rb.size == 0 {
		return
	}

	if rb.head < rb.tail {
		for cursor := rb.head; cursor < rb.tail; cursor++ {
			consumer(rb.array[cursor])
		}
	} else {
		for cursor := rb.head; cursor < len(rb.array); cursor++ {
			consumer(rb.array[cursor])
		}
		for cursor := 0; cursor < rb.tail; cursor++ {
			consumer(rb.array[cursor])
		}
	}
}

// Drain calls the consumer for each element in the buffer, while also dequeueing that entry.
func (rb *RingBuffer) Drain(consumer func(value interface{})) {
	if rb.size == 0 {
		return
	}

	len := rb.Len()
	for i := 0; i < len; i++ {
		consumer(rb.Dequeue())
	}
}

// EachUntil calls the consumer for each element in the buffer with a stopping condition in head=>tail order.
func (rb *RingBuffer) EachUntil(consumer func(value interface{}) bool) {
	if rb.size == 0 {
		return
	}

	if rb.head < rb.tail {
		for cursor := rb.head; cursor < rb.tail; cursor++ {
			if !consumer(rb.array[cursor]) {
				return
			}
		}
	} else {
		for cursor := rb.head; cursor < len(rb.array); cursor++ {
			if !consumer(rb.array[cursor]) {
				return
			}
		}
		for cursor := 0; cursor < rb.tail; cursor++ {
			if !consumer(rb.array[cursor]) {
				return
			}
		}
	}
}

// ReverseEachUntil calls the consumer for each element in the buffer with a stopping condition in tail=>head order.
func (rb *RingBuffer) ReverseEachUntil(consumer func(value interface{}) bool) {
	if rb.size == 0 {
		return
	}

	if rb.head < rb.tail {
		for cursor := rb.tail - 1; cursor >= rb.head; cursor-- {
			if !consumer(rb.array[cursor]) {
				return
			}
		}
	} else {
		for cursor := rb.tail; cursor > 0; cursor-- {
			if !consumer(rb.array[cursor]) {
				return
			}
		}
		for cursor := len(rb.array) - 1; cursor >= rb.head; cursor-- {
			if !consumer(rb.array[cursor]) {
				return
			}
		}

	}
}

func (rb *RingBuffer) String() string {
	var values []string
	for _, elem := range rb.AsSlice() {
		values = append(values, fmt.Sprintf("%v", elem))
	}
	return strings.Join(values, " <= ")
}

func arrayClear(source []interface{}, index, length int) {
	for x := 0; x < length; x++ {
		absoluteIndex := x + index
		source[absoluteIndex] = nil
	}
}

func arrayCopy(source []interface{}, sourceIndex int, destination []interface{}, destinationIndex, length int) {
	for x := 0; x < length; x++ {
		from := sourceIndex + x
		to := destinationIndex + x

		destination[to] = source[from]
	}
}
