//go:build windows

package sysutil

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type WindowBounds struct {
	X      int
	Y      int
	Width  int
	Height int
}

type rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type monitorInfo struct {
	CbSize    uint32
	RcMonitor rect
	RcWork    rect
	DwFlags   uint32
}

var (
	user32                  = windows.NewLazySystemDLL("user32.dll")
	procFindWindowW         = user32.NewProc("FindWindowW")
	procGetWindowRect       = user32.NewProc("GetWindowRect")
	procSetWindowPos        = user32.NewProc("SetWindowPos")
	procMonitorFromRect     = user32.NewProc("MonitorFromRect")
	procGetMonitorInfoW     = user32.NewProc("GetMonitorInfoW")
	procIsIconic            = user32.NewProc("IsIconic")
	procGetWindowPlacement  = user32.NewProc("GetWindowPlacement")
)

type point struct {
	X int32
	Y int32
}

type windowPlacement struct {
	Length           uint32
	Flags            uint32
	ShowCmd          uint32
	PtMinPosition    point
	PtMaxPosition    point
	RcNormalPosition rect
}

const (
	monitorDefaultToNearest = 2
	swpNoZOrder             = 0x0004
)

// IsMainWindowMinimized 判断主窗口是否处于最小化状态。
func IsMainWindowMinimized(title string) (bool, error) {
	hwnd, err := findWindow(title)
	if err != nil {
		return false, err
	}
	ret, _, _ := procIsIconic.Call(uintptr(hwnd))
	return ret != 0, nil
}

// GetMainWindowBounds 读取主窗口的绝对屏幕坐标与尺寸。
func GetMainWindowBounds(title string) (WindowBounds, error) {
	hwnd, err := findWindow(title)
	if err != nil {
		return WindowBounds{}, err
	}
	if minimized, err := isIconic(hwnd); err == nil && minimized {
		return normalBoundsFromPlacement(hwnd)
	}
	return boundsFromWindowRect(hwnd)
}

func boundsFromWindowRect(hwnd windows.Handle) (WindowBounds, error) {
	var windowRect rect
	ret, _, callErr := procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&windowRect)))
	if ret == 0 {
		return WindowBounds{}, syscallError("GetWindowRect", callErr)
	}
	return rectToBounds(windowRect), nil
}

func normalBoundsFromPlacement(hwnd windows.Handle) (WindowBounds, error) {
	placement := windowPlacement{Length: uint32(unsafe.Sizeof(windowPlacement{}))}
	ret, _, callErr := procGetWindowPlacement.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&placement)),
	)
	if ret == 0 {
		return WindowBounds{}, syscallError("GetWindowPlacement", callErr)
	}
	return rectToBounds(placement.RcNormalPosition), nil
}

func rectToBounds(windowRect rect) WindowBounds {
	return WindowBounds{
		X:      int(windowRect.Left),
		Y:      int(windowRect.Top),
		Width:  int(windowRect.Right - windowRect.Left),
		Height: int(windowRect.Bottom - windowRect.Top),
	}
}

func isIconic(hwnd windows.Handle) (bool, error) {
	ret, _, callErr := procIsIconic.Call(uintptr(hwnd))
	if ret == 0 && callErr != nil && callErr != windows.ERROR_SUCCESS {
		return false, syscallError("IsIconic", callErr)
	}
	return ret != 0, nil
}

// SetMainWindowBounds 以绝对屏幕坐标设置主窗口尺寸与位置。
func SetMainWindowBounds(title string, bounds WindowBounds) error {
	hwnd, err := findWindow(title)
	if err != nil {
		return err
	}
	ret, _, callErr := procSetWindowPos.Call(
		uintptr(hwnd),
		0,
		uintptr(bounds.X),
		uintptr(bounds.Y),
		uintptr(bounds.Width),
		uintptr(bounds.Height),
		uintptr(swpNoZOrder),
	)
	if ret == 0 {
		return syscallError("SetWindowPos", callErr)
	}
	return nil
}

const minSavedWindowWidth = 320
const minSavedWindowHeight = 240

// ValidWindowBounds 判断窗口尺寸是否可用于保存/恢复。
func ValidWindowBounds(bounds WindowBounds) bool {
	return bounds.Width >= minSavedWindowWidth && bounds.Height >= minSavedWindowHeight
}

// ClampWindowBounds 修正窗口尺寸与位置，确保窗口完整落在可用工作区内。
func ClampWindowBounds(bounds WindowBounds, fallbackWidth, fallbackHeight int) WindowBounds {
	if bounds.Width < minSavedWindowWidth {
		bounds.Width = fallbackWidth
	}
	if bounds.Height < minSavedWindowHeight {
		bounds.Height = fallbackHeight
	}

	workArea, err := monitorWorkAreaForBounds(bounds)
	if err != nil {
		return bounds
	}

	maxWidth := int(workArea.Right - workArea.Left)
	maxHeight := int(workArea.Bottom - workArea.Top)
	if maxWidth > 0 && bounds.Width > maxWidth {
		bounds.Width = maxWidth
	}
	if maxHeight > 0 && bounds.Height > maxHeight {
		bounds.Height = maxHeight
	}

	minX := int(workArea.Left)
	minY := int(workArea.Top)
	maxX := int(workArea.Right) - bounds.Width
	maxY := int(workArea.Bottom) - bounds.Height
	if maxX < minX {
		maxX = minX
	}
	if maxY < minY {
		maxY = minY
	}

	bounds.X = clampInt(bounds.X, minX, maxX)
	bounds.Y = clampInt(bounds.Y, minY, maxY)
	return bounds
}

func clampInt(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func monitorWorkAreaForBounds(bounds WindowBounds) (rect, error) {
	windowRect := rect{
		Left:   int32(bounds.X),
		Top:    int32(bounds.Y),
		Right:  int32(bounds.X + bounds.Width),
		Bottom: int32(bounds.Y + bounds.Height),
	}
	ret, _, callErr := procMonitorFromRect.Call(
		uintptr(unsafe.Pointer(&windowRect)),
		uintptr(monitorDefaultToNearest),
	)
	if ret == 0 {
		return rect{}, syscallError("MonitorFromRect", callErr)
	}
	info := monitorInfo{CbSize: uint32(unsafe.Sizeof(monitorInfo{}))}
	ret, _, callErr = procGetMonitorInfoW.Call(ret, uintptr(unsafe.Pointer(&info)))
	if ret == 0 {
		return rect{}, syscallError("GetMonitorInfoW", callErr)
	}
	return info.RcWork, nil
}

func findWindow(title string) (windows.Handle, error) {
	windowTitle, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return 0, err
	}
	ret, _, callErr := procFindWindowW.Call(0, uintptr(unsafe.Pointer(windowTitle)))
	if ret == 0 {
		return 0, syscallError("FindWindowW", callErr)
	}
	return windows.Handle(ret), nil
}

func syscallError(name string, err error) error {
	if err == nil || err == windows.ERROR_SUCCESS {
		return fmt.Errorf("%s failed", name)
	}
	return fmt.Errorf("%s failed: %w", name, err)
}
