package main

import (
	"github.com/tikaflow/app-go"
	"memo/assets"
)

var (
	appName = "西夏备忘录"
)

func main() {
	app.New(appName, func() {
		app.InitTray(appName)
		app.SetTrayIcon(assets.R["/img/app-icon.ico"])
		tray := app.GetTray()

		m := tray.AddItem(nil, "打开主界面")
		m.OnClick(func(item *app.TrayItem) {
			app.OpenWindow(1200, 800)
		})
		tray.AddSeparator()
		q := tray.AddItem(nil, "退出")
		q.OnClick(func(item *app.TrayItem) {
			app.Notify("已退出")
			app.Quit()
		})
	})

	app.OnReady(func() {
		initDB()
		app.OpenWindow(1200, 800)
	})
	app.SetWebRoot(assets.GetStatic())
	app.SetNoArgTask(task)
	app.Run()
}
