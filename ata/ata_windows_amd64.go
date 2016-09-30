// 30 september 2016
package ata

// because encoding/binary will refuse to take the size of uintptr, by design (https://github.com/golang/go/issues/8381)
// TODO ask anyway?
type uint64 _ULONG_PTR

// if that ever changes:
// notes on ULONG_PTR and uintptr:
// golang.org/x/sys/windows provides a wrapper for FlushViewOfFile()
// the second parameter of that function is uintptr
// in the real API, that's SIZE_T
// https://msdn.microsoft.com/en-us/library/windows/desktop/aa383751(v=vs.85).aspx says SIZE_T is a typedef for ULONG_PTR
// therefore, ULONG_PTR == uintptr
