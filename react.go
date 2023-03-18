package react

type Tool struct {
	// Name is the name of the tool, will be used for lookup
	Name string
	// Description is a short description of the tool with usage info
	Description string
	// Run is the function that will be called when the tool is invoked
	Run func(arg string) string
}

