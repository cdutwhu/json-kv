package preprocess

import "log"

// FailOnErr :
func FailOnErr(format string, v ...interface{}) {
	for _, p := range v {
		switch p.(type) {
		case error:
			{
				if p != nil {
					log.Fatalf(format, v...)
				}
			}
		}
	}
}
