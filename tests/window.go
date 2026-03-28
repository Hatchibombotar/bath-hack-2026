package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

var (
	user32             = syscall.MustLoadDLL("user32.dll")
	procEnumWindows    = user32.MustFindProc("EnumWindows")
	procGetWindowTextW = user32.MustFindProc("GetWindowTextW")
	procGetWindowRect  = user32.MustFindProc("GetWindowRect")
)

type RECT struct {
	left   int32 // or Left, Top, etc. if this type is to be exported
	top    int32
	right  int32
	bottom int32
}

func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
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

func FindWindow(title string) (syscall.Handle, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		titleBuffer := make([]uint16, 200)
		rectBuffer := RECT{}
		_, err1 := GetWindowText(h, &titleBuffer[0], int32(len(titleBuffer)))
		_, err2 := GetWindowRect(h, &rectBuffer, 16)
		if err1 != nil || err2 != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		fmt.Println(syscall.UTF16ToString(titleBuffer))
		if syscall.UTF16ToString(titleBuffer) == title {
			// note the window
			hwnd = h
			fmt.Println(rectBuffer.bottom, rectBuffer.right, rectBuffer.left, rectBuffer.top)
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	EnumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("No window with title '%s' found", title)
	}
	return hwnd, nil
}

func main() {
	const title = "GitHub Desktop"
	h, err := FindWindow(title)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found '%s' window: handle=0x%x\n", title, h)
}
