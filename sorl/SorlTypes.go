package main

// SorlMap for Vars
type SorlMap map[string]string

// SorlConfig for Sorl Config
type SorlConfig map[string]string

//SorlConfigProperty for Sorl Configs
type SorlConfigProperty map[string]SorlConfig

// PathSep separator
const PathSep string = "/"

const (
	CLR_BLACK = "\x1b[30;1m"

	CLR_RED = "\x1b[31;1m"

	CLR_GREEN = "\x1b[32;1m"

	CLR_YELLOW = "\x1b[33;1m"

	CLR_BLUE = "\x1b[34;1m"

	CLR_MAGENTA = "\x1b[35;1m"

	CLR_CYAN = "\x1b[36;1m"

	CLR_WHITE = "\x1b[37;1m"

	CLR_UNCOLOR = "\x1b[0m"
)
