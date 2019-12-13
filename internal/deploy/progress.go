package deploy

import "strings"

const (
	ProgressUnmapped                    = iota
	ProgressNoChangeset                 = iota
	ProgressAwaitingChangesetCreation   = iota
	ProgressAwaitingChangesetCompletion = iota
	ProgressStackCompletion             = iota
	ProgressEOF                         = iota
)

// Progress parses a deploy stdout line to it's progress value
func Progress(line string) (progress int) {
	l := strings.Trim(strings.ToLower(line), " .")

	switch l {
	case "waiting for changeset to be created":
		progress = ProgressAwaitingChangesetCreation
	case "waiting for stack create/update to complete":
		progress = ProgressAwaitingChangesetCompletion
	case "no changes to deploy":
		progress = ProgressNoChangeset
	case "successfully created/updated stack":
		progress = ProgressStackCompletion
	default:
		progress = ProgressUnmapped
	}
	return
}
