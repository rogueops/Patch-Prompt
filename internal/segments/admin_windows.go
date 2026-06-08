//go:build windows

package segments

import (
	"syscall"
	"unsafe"
)

// Admin reports whether the current process token is elevated (Windows).
func Admin(_ Context) (string, bool) {
	advapi32 := syscall.NewLazyDLL("advapi32.dll")
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	openProcessToken := advapi32.NewProc("OpenProcessToken")
	getTokenInformation := advapi32.NewProc("GetTokenInformation")
	getCurrentProcess := kernel32.NewProc("GetCurrentProcess")
	closeHandle := kernel32.NewProc("CloseHandle")

	const tokenQuery = 0x0008
	const tokenElevation = 20 // TokenElevation class

	proc, _, _ := getCurrentProcess.Call()
	var token syscall.Handle
	ret, _, _ := openProcessToken.Call(proc, tokenQuery, uintptr(unsafe.Pointer(&token)))
	if ret == 0 {
		return "", false
	}
	defer closeHandle.Call(uintptr(token))

	var elevated uint32
	var retLen uint32
	ret, _, _ = getTokenInformation.Call(
		uintptr(token),
		uintptr(tokenElevation),
		uintptr(unsafe.Pointer(&elevated)),
		unsafe.Sizeof(elevated),
		uintptr(unsafe.Pointer(&retLen)),
	)
	if ret == 0 || elevated == 0 {
		return "", false
	}
	return "ADMIN", true
}
