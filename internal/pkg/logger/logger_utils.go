package logger

func getArgs(args ...interface{}) (resStr string, resArgs []interface{}) {
	logStr := args[0]
	resStr, ok := logStr.(string)
	if !ok {
		return "", resArgs
	}
	if len(args[1:]) > 0 {
		return resStr, args[1:]
	}

	return
}
