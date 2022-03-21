package cmd

import "path"

type TranslateCommand struct {
	TrFile string `short:"l" long:"list" default:"tr.js" required:"true" description:"Файл со списком фраз" long-description:"Формируемый файл, содержащий список фраз для перевода"`
	ExtensionOpts
}

func (tc *TranslateCommand) Execute(args []string) error {
	return process(tc.InputFile, tc.OutputFile, path.Join(".", tc.TrFile), tc.Expr, false)
}
