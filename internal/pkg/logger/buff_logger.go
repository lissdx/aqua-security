package logger

import (
	"bytes"
	"fmt"
	standardLogger "log"
	"os"
)

// BufferLogger used for test
type BufferLogger struct {
	loggerBuff    *standardLogger.Logger
	logger        *standardLogger.Logger
	Out           *bytes.Buffer
	printToIO     bool
	printToBuffer bool
}

func NewBufferLogger(printToBuffer bool, printToIO bool) TOLogger {
	outBuff := new(bytes.Buffer)
	stdLogg := standardLogger.New(os.Stdout, "", 0)

	return &BufferLogger{
		loggerBuff:    standardLogger.New(outBuff, "", 0),
		logger:        stdLogg,
		Out:           outBuff,
		printToIO:     printToIO,
		printToBuffer: printToBuffer,
	}
}

func (bl *BufferLogger) Debug(args ...interface{}) {
	fmtStr, a := getArgs(args...)

	logStr := func() string {
		if len(a) > 0 {
			return fmt.Sprintf(fmtStr, a...)
		}
		return fmtStr
	}()

	if bl.printToBuffer {
		bl.loggerBuff.Println("DEBUG ", logStr)
	}
	if bl.printToIO {
		bl.logger.Println("DEBUG ", logStr)
	}
}

func (bl *BufferLogger) Info(args ...interface{}) {
	fmtStr, a := getArgs(args...)

	logStr := func() string {
		if len(a) > 0 {
			return fmt.Sprintf(fmtStr, a...)
		}
		return fmtStr
	}()

	if bl.printToBuffer {
		bl.loggerBuff.Println("INFO ", logStr)
	}
	if bl.printToIO {
		bl.logger.Println("INFO ", logStr)
	}
}

func (bl *BufferLogger) Warn(args ...interface{}) {
	fmtStr, a := getArgs(args...)

	logStr := func() string {
		if len(a) > 0 {
			return fmt.Sprintf(fmtStr, a...)
		}
		return fmtStr
	}()

	if bl.printToBuffer {
		bl.loggerBuff.Println("WARN ", logStr)
	}
	if bl.printToIO {
		bl.logger.Println("WARN ", logStr)
	}
}

func (bl *BufferLogger) Error(args ...interface{}) {
	fmtStr, a := getArgs(args...)

	logStr := func() string {
		if len(a) > 0 {
			return fmt.Sprintf(fmtStr, a...)
		}
		return fmtStr
	}()

	if bl.printToBuffer {
		bl.loggerBuff.Println("ERROR ", logStr)
	}
	if bl.printToIO {
		bl.logger.Println("ERROR ", logStr)
	}
}

func (bl *BufferLogger) Panic(args ...interface{}) {
	fmtStr, a := getArgs(args...)

	logStr := func() string {
		if len(a) > 0 {
			return fmt.Sprintf(fmtStr, a...)
		}
		return fmtStr
	}()

	if bl.printToBuffer {
		bl.loggerBuff.Println("PANIC ", logStr)
	}
	if bl.printToIO {
		bl.logger.Println("PANIC ", logStr)
	}
}

func (bl *BufferLogger) Fatal(args ...interface{}) {
	fmtStr, a := getArgs(args...)

	logStr := func() string {
		if len(a) > 0 {
			return fmt.Sprintf(fmtStr, a...)
		}
		return fmtStr
	}()

	if bl.printToBuffer {
		bl.loggerBuff.Println("FATAL ", logStr)
	}
	if bl.printToIO {
		bl.logger.Println("FATAL ", logStr)
	}
}

func (bl *BufferLogger) GetBufferedString() string {
	return bl.Out.String()
}
