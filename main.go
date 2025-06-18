package main

import (
	"github.com/OscillatingBlock/url_shortner_go/app"
)

var GlobalConfig app.Configuration = app.InitConfig()

func main() {
	app := app.Default(GlobalConfig)
	app.Run(GlobalConfig.Port)
}
