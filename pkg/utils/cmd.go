package utils

type ResponseExec struct {
	ExitCode int    `json:"exit_code"`
	Err      error  `json:"err"`
	Stdout   string `json:"stdout"`
}
