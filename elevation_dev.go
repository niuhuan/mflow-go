//go:build dev

package main

// maybeElevate 在 dev 模式下不做提权，直接返回 false。
func maybeElevate() bool {
	return false
}
