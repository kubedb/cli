package sh

import (
	"bytes"
	"sync"
)

type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

// Write appends data to the buffer in a thread-safe manner.
func (sb *safeBuffer) Write(p []byte) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.Write(p)
}

// Bytes returns a copy of the buffer's contents in a thread-safe manner.
func (sb *safeBuffer) Bytes() []byte {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return bytes.Clone(sb.buf.Bytes())
}

func (sb *safeBuffer) Reset() {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.buf.Reset()
}
