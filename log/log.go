package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log ....
var Log = logrus.New()

func init() {
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		Log.Out = file
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}
}
