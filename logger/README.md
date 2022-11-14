# Logger

### The `logger` package implements the core logging functionality.

You can set your logger which implements the `Logger` interface or use the `standard logger`.

By default, the `standard logger` is initialized:

- with the time format `2006/01/02 15:04:05.000`;
- with logging format `FormatText`;
- at the highest logging level `LevelTrace`;
- printing the name of the calling function `true`.

You can set the required time format, level, log message format and additional output, as in the following examples:

- at initialization

```go
package main

import (
	"os"

	"github.com/gromey/proto-rest/logger"
)

func main() {
	file, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	logger.Info("Hello World!")
	// 2022/11/01 15:04:05.678 [INFO] /.../proto-rest/main.go:21 Func: main() Hello World!

	lCfg := &logger.Config{
		TimeFormat:    "2006-01-02 15:04:05",
		Format:        logger.FormatJSON,
		AdditionalOut: file,
		Level:         logger.LevelInfo,
		FuncName:      true,
	}

	logger.SetLogger(logger.New(lCfg))

	logger.Info("Hello World!")
	// {"time":"2022-11-01 15:04:05","level":"INFO","func":"main()","message":"Hello World!"}
}
```

- after initialization

```go
package main

import (
	"os"

	"github.com/gromey/proto-rest/logger"
)

func main() {
	l := logger.New(nil)
	logger.SetLogger(l)

	file, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	l.SetTimeFormat("2006-01-02 15:04:05")
	l.SetFormatter(logger.FormatJSON)
	l.SetAdditionalOut(file)
	l.SetLevel(logger.LevelInfo)
	l.SetFuncNamePrinting(false)

	logger.Info("Hello World!")
	// {"time":"2022-11-01 15:04:05","level":"INFO","message":"Hello World!"}
}
```

To reduce unnecessary memory allocations, it is recommended to check the logging level before calling the log function.

Example:

```go
	if logger.InLevel(logger.LevelInfo) {
		logger.Info("Hello World!")
	}
```