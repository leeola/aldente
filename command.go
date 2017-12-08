package motley

type CommandOutput struct {
	Output       string
	CommandName  string
	MachineName  string
	ProviderName string
	ProviderType string
	Error        error
}
