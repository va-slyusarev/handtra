package cmd

type PrepareCommand struct {
	ExtensionOpts
}

func (pc *PrepareCommand) Execute(args []string) error {
	return process(pc.InputFile, pc.OutputFile, "", pc.Expr, true)
}
