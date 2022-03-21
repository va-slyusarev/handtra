// Copyright © 2022 Valentin Slyusarev <va.slyusarev@gmail.com>
package main

import (
	"os"
	"path"

	"github.com/jessevdk/go-flags"

	"github.com/va-slyusarev/handtra/cmd"
)

type Opts struct {
	PrepareCmd   cmd.PrepareCommand   `command:"prepare"     description:"Подготовка списка" long-description:"Подготовить список фраз для перевода"`
	TranslateCmd cmd.TranslateCommand `command:"translate"   description:"Перевод" long-description:"Перевод с использованием списка фраз"`
	VersionCmd   cmd.VersionCommand   `command:"version" alias:"v" description:"Версия приложения"`

	InputFile  string `short:"i" long:"in" required:"true" description:"Входящий файл" long-description:"Входящий файл для перевода"`
	OutputFile string `short:"o" long:"out" default:"tr.js" required:"true" description:"Исходящий файл" long-description:"Исходящий файл. В зависимости от команды содержит список фраз для перевода, либо итоговый перевод"`
	Expr       string `short:"e" long:"expr" default:".*" required:"true" description:"Регулярное выражение" long-description:"Регулярное выражение, которое определяет фразу для перевода в считываемом построчно файле"`
}

var revision = "develop"

func main() {
	var opts Opts
	p := flags.NewParser(&opts, flags.Default)
	p.CommandHandler = func(command flags.Commander, args []string) error {
		if ext, ok := command.(cmd.ExtCommander); ok {
			ext.SetExt(cmd.ExtensionOpts{
				AppRevision: revision,
				InputFile:   path.Join(".", opts.InputFile),
				OutputFile:  path.Join(".", opts.OutputFile),
				Expr:        opts.Expr,
			})
			return ext.Execute(args)
		}
		return command.Execute(args)
	}

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
