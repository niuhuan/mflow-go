//go:build !windows

package sysutil

type WindowBounds struct {
	X      int
	Y      int
	Width  int
	Height int
}

func GetMainWindowBounds(title string) (WindowBounds, error) {
	return WindowBounds{}, nil
}

func SetMainWindowBounds(title string, bounds WindowBounds) error {
	return nil
}

func ClampWindowBounds(bounds WindowBounds, fallbackWidth, fallbackHeight int) WindowBounds {
	if bounds.Width <= 0 {
		bounds.Width = fallbackWidth
	}
	if bounds.Height <= 0 {
		bounds.Height = fallbackHeight
	}
	return bounds
}
