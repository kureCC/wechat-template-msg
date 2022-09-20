package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"wechat-template-msg/cmd"
)

func main() {
	var (
		ctx = gctx.New()
	)

	// 加载配置
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("config.yaml")
	configName, _ := g.Cfg().Get(ctx, "server.config")
	g.Log().Info(ctx, "环境:", configName.String())

	command, err := gcmd.NewFromObject(cmd.Main)
	if err != nil {
		panic(err)
	}
	err = command.AddObject(
		cmd.Consumer,
		cmd.Producer,
	)
	if err != nil {
		panic(err)
	}
	command.Run(ctx)
}
