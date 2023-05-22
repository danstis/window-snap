package main

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/JamesHovious/w32"
	"golang.org/x/sys/windows"
)

func main() {
	err := RunExeAndSnapWindow("c:\\windows\\system32\\notepad.exe", 1, 1)
	if err != nil {
		fmt.Println(err)
	}
}

func startExecutable(path string) (syscall.Handle, error) {
	// Set up startup info and process information structs
	var si syscall.StartupInfo
	var pi syscall.ProcessInformation

	// Start the process
	w32.CreateProcessA(nil, path, nil, nil, true, 0, nil, nil, &si, &pi)

	// Return the process handle
	return pi.Process, nil
}

func RunExeAndSnapWindow(filePath string, quadrant int, monitor int) error {
	h, err := startExecutable(filePath)
	if err != nil {
		return err
	}

	// Determine the screen bounds of the specified monitor
	screen, err := getMonitorBounds(monitor - 1)
	if err != nil {
		return err
	}

	// Calculate the window size and position for the specified quadrant
	width := int(screen.Right-screen.Left) / 2
	height := int(screen.Bottom-screen.Top) / 2
	var x, y int
	switch quadrant {
	case 1:
		x = int(screen.Right - int32(width))
		y = int(screen.Top)
	case 2:
		x = int(screen.Right - int32(width))
		y = int(screen.Bottom - int32(height))
	case 3:
		x = int(screen.Left)
		y = int(screen.Bottom - int32(height))
	case 4:
		x = int(screen.Left)
		y = int(screen.Top)
	}

	// Move and resize the window
	w32.MoveWindow(w32.HWND(h), x, y, width, height, true)

	// Return the process object
	return nil
}

// findWindow finds the window handle for the specified process ID.
func findWindow(pid int) (windows.Handle, error) {
	var window windows.Handle
	callback := windows.NewCallback(func(h windows.Handle, p uintptr) uintptr {
		var processId uint32
		windows.GetWindowThreadProcessId(windows.HWND(h), &processId)
		if int(processId) == pid {
			window = h
			return 0
		}
		return 1
	})
	windows.EnumWindows(callback, nil)
	if window == 0 {
		return 0, errors.New("window not found")
	}
	return window, nil
}

// Get the bounds of the specified monitor.
func getMonitorBounds(monitor int) (w32.RECT, error) {
	var mi w32.MONITORINFO
	mi.CbSize = uint32(unsafe.Sizeof(mi))
	if state := w32.GetMonitorInfo(w32.HMONITOR(monitor), &mi); state == 0 {
		return w32.RECT{}, errors.New("failed to get monitor info")
	}
	return mi.RcMonitor, nil
}
