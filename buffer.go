package tracing

import (
	"bytes"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		return allocBuffer()
	}}

// Allocate a new line buffer with an initial capacity of 1KiB
func allocBuffer() *bytes.Buffer {
	const initialCapacity = 1024
	b := make([]byte, 0, initialCapacity)
	return bytes.NewBuffer(b)
}

// Get a line buffer from the buffer pool for use
func GetBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

// Clear a line buffer (while still retaining capacity)
// and potentially return it to the main buffer pool
func CloseBuffer(buf *bytes.Buffer) {
	// Only cache buffers up to 16KiB large
	const maxCapacity = 16 << 10
	if buf.Cap() > maxCapacity {
		return
	}
	buf.Reset()
	bufferPool.Put(buf)
}
