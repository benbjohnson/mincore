package mincore_test

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"
	"time"
	"unsafe"

	"github.com/benbjohnson/mincore"
)

func TestMincore(t *testing.T) {
	// Create temporary file.
	t0 := time.Now()
	f, err := ioutil.TempFile("", "mincore-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()
	t.Logf("FILE=%s", f.Name())
	t.Logf("temp file time: %v", time.Since(t0))

	// Truncate it to 10 pages
	t0 = time.Now()
	sz := int64(10 * os.Getpagesize())
	if err := f.Truncate(sz); err != nil {
		t.Fatal(err)
	}
	t.Logf("truncate time: %v", time.Since(t0))

	// Memory map the file.
	t0 = time.Now()
	b, err := syscall.Mmap(int(f.Fd()), 0, int(sz), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := syscall.Syscall(syscall.SYS_MADVISE, uintptr(unsafe.Pointer(&b[0])), uintptr(syscall.MADV_RANDOM), 0); err != 0 {
		t.Fatal(err)
	}
	t.Logf("mmap time: %v", time.Since(t0))
	defer syscall.Munmap(b)

	// Read from the first page & the 5th page.
	t.Log("Page 0", b[0])
	t.Log("Page 9", b[9*os.Getpagesize()])

	// Create a byte slice to hold page values.
	t0 = time.Now()
	vec := make([]byte, int(sz)/os.Getpagesize())
	if err := mincore.Mincore(unsafe.Pointer(&b[0]), uint64(sz), vec); err != nil {
		t.Fatal(err)
	}
	t.Logf("mincore time: %v", time.Since(t0))

	t.Logf("VEC %x", vec)
}
