//go:build !windows

package sysutil

type WindowBounds struct {
	X      int
	Y      int
	Width  int
	Height int
}

func IsMainWindowMinimized(title string) (bool, error) {
	return false, nil
}

func GetMainWindowBounds(title string) (WindowBounds, error) {
	return WindowBounds{}, nil
}

func SetMainWindowBounds(title string, bounds WindowBounds) error {
	return nil
}

func ValidWindowBounds(bounds WindowBounds) bool {
	return bounds.Width >= 320 && bounds.Height >= 240
}

func ClampWindowBounds(bounds WindowBounds, fallbackWidth, fallbackHeight int) WindowBounds {
	if bounds.Width < 320 {
		bounds.Width = fallbackWidth
	}
	if bounds.Height < 240 {
		bounds.Height = fallbackHeight
	}
	return bounds
}
