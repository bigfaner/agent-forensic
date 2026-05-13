package parser

// IsReadTool returns true for tool names that read files.
// Accepts known aliases across Claude Code JSONL format versions.
func IsReadTool(name string) bool {
	return name == "Read" || name == "Read1"
}

// IsEditTool returns true for tool names that modify files.
func IsEditTool(name string) bool {
	return name == "Write" || name == "Edit"
}

// IsFileTool returns true for tool names that operate on files (read or edit).
func IsFileTool(name string) bool {
	return IsReadTool(name) || IsEditTool(name)
}

// IsBashTool returns true for tool names that run shell commands.
func IsBashTool(name string) bool {
	return name == "Bash"
}

// IsAgentTool returns true for tool names that spawn sub-agents.
func IsAgentTool(name string) bool {
	return name == "Agent" || name == "SubAgent"
}
