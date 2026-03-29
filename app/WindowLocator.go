package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.MustLoadDLL("user32.dll")
	procGetWindowTextW      = user32.MustFindProc("GetWindowTextW")
	procGetWindowRect       = user32.MustFindProc("GetWindowRect")
	procGetForegroundWindow = user32.MustFindProc("GetForegroundWindow")
)

type RECT struct {
	left   int32 // or Left, Top, etc. if this type is to be exported
	top    int32
	right  int32
	bottom int32
}

func GetWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetWindowRect(hwnd syscall.Handle, rect *RECT, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowRect.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(rect)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetForegroundWindowInfo() (RECT, string) {
	active, _, _ := procGetForegroundWindow.Call()
	fmt.Println(active)
	titleBuffer := make([]uint16, 200)
	rectBuffer := RECT{}
	//GetWindowText(syscall.Handle(active), &titleBuffer[0], int32(len(titleBuffer)))
	GetWindowRect(syscall.Handle(active), &rectBuffer, 16)
	fmt.Println(syscall.UTF16ToString(titleBuffer))
	fmt.Println(rectBuffer.bottom, rectBuffer.right, rectBuffer.left, rectBuffer.top)
	return rectBuffer, syscall.UTF16ToString(titleBuffer)
}
