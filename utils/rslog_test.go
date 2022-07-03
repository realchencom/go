package utils

import "testing"

func TestLog(t *testing.T) {
	defer Logger.Sync()
	Logger.Debug("年后")
}
