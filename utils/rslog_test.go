package rs

import "testing"

func TestLog(t *testing.T) {
	defer func(Log *Logger) {
		err := Log.Sync()
		if err != nil {

		}
	}(&Log)
	Log.Debug("年后")
}
