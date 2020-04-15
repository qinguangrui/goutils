package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/axiaoxin-com/logging"

	"go.uber.org/zap"
)

// level 全局变量，便于动态修改，初始化为 Debug 级别
var level zap.AtomicLevel = zap.NewAtomicLevelAt(zap.DebugLevel)

func main() {
	/* change log level on fly */

	// 创建指定Level的logger，并开启http服务
	options := logging.Options{
		Level:           level,
		AtomicLevelAddr: ":2012",
	}
	logger, _ := logging.NewLogger(options)
	logger.Debug("Debug level msg", zap.Any("current level", level.Level()))
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:03:17.799767","logger":"root","caller":"example/atomiclevel.go:main:26","msg":"Debug level msg","pid":6088,"current level":"debug"}

	// 使用SetLevel动态修改logger 日志级别为error
	// 实际应用中可以监听配置文件中日志级别配置项的变化动态调用该函数
	level.SetLevel(zap.ErrorLevel)
	// Info 级别将不会被打印
	logger.Info("Info level msg will not be logged")
	// 只会打印error以上
	logger.Error("Error level msg", zap.Any("current level", level.Level()))
	// Output:
	// {"level":"ERROR","time":"2020-04-15 18:03:17.799999","logger":"root","caller":"example/atomiclevel.go:main:34","msg":"Error level msg","pid":6088,"current level":"error","stacktrace":"main.main\n\t/Users/ashin/go/src/logging/example/atomiclevel.go:34\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:203"}

	// 通过HTTP方式动态修改当前的error level为debug level
	// 查询当前 level
	url := "http://localhost" + options.AtomicLevelAddr
	resp, _ := http.Get(url)
	content, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println("currentlevel:", string(content))
	// Output: currentlevel: {"level":"error"}

	logger.Info("Info level will not be logged")

	// 修改level为debug
	c := &http.Client{}
	req, _ := http.NewRequest("PUT", url, strings.NewReader(`{"level": "debug"}`))
	resp, _ = c.Do(req)
	content, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println("newlevel:", string(content))
	// Output: newlevel: {"level":"debug"}

	logger.Debug("level is changed on fly!")

	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:03:17.805293","logger":"root","caller":"example/atomiclevel.go:main:57","msg":"level is changed on fly!","pid":6088}
}
