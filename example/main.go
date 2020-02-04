package main

import (
	"github.com/suifengpiao14/crud-mysql/cmd"
	"github.com/suifengpiao14/crud-mysql/config"
)

func main() {
	config.Load()
	cmd.Execute()
}
