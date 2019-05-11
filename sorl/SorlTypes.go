package main

// SorlMap for Vars
type SorlMap map[string]string

// SorlConfig for Sorl Config
type SorlConfig map[string]string

//SorlConfigProperty for Sorl Configs
type SorlConfigProperty map[string]SorlConfig

// PathSep separator
const PathSep string = "/"

// Ansi Colors const
const (
	ClrBlack = "\x1b[30;1m"

	ClrRed = "\x1b[31;1m"

	ClrGreen = "\x1b[32;1m"

	ClrYellow = "\x1b[33;1m"

	ClrBlue = "\x1b[34;1m"

	ClrMagenta = "\x1b[35;1m"

	ClrCyan = "\x1b[36;1m"

	ClrWhite = "\x1b[37;1m"

	ClrUnColor = "\x1b[0m"
)
