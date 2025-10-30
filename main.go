package main

import (
	"fmt"

	"github.com/asqit/open-ping/web"
)

func main() {
	fmt.Println(`
░█▀█░█▀█░█▀▀░█▀█░░░░░█▀█░▀█▀░█▀█░█▀▀
░█░█░█▀▀░█▀▀░█░█░▄▄▄░█▀▀░░█░░█░█░█░█
░▀▀▀░▀░░░▀▀▀░▀░▀░░░░░▀░░░▀▀▀░▀░▀░▀▀▀
====================================`)
	init_db()
	defer close_db()
	web.Start_dashboard(db)

	//cfg, err := read_config()
	//if err != nil {
	//	panic(err)
	//}
	//
	//run_monitor(cfg)
}
