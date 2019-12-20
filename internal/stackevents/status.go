package stackevents

const (
	Unknown = iota
	Ok
	Progress
	Fail
)

var statusMap map[string]int = map[string]int{
	"CREATE_COMPLETE":                              Ok,
	"CREATE_IN_PROGRESS":                           Progress,
	"CREATE_FAILED":                                Fail,
	"DELETE_COMPLETE":                              Ok,
	"DELETE_FAILED":                                Fail,
	"DELETE_IN_PROGRESS":                           Progress,
	"ROLLBACK_COMPLETE":                            Ok,
	"ROLLBACK_FAILED":                              Fail,
	"ROLLBACK_IN_PROGRESS":                         Progress,
	"UPDATE_COMPLETE":                              Ok,
	"UPDATE_COMPLETE_CLEANUP_IN_PROGRESS":          Progress,
	"UPDATE_IN_PROGRESS":                           Progress,
	"UPDATE_ROLLBACK_COMPLETE":                     Ok,
	"UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS": Progress,
	"UPDATE_ROLLBACK_FAILED":                       Fail,
	"UPDATE_ROLLBACK_IN_PROGRESS":                  Progress,
	"UPDATE_FAILED":                                Fail,
	"IMPORT_IN_PROGRESS":                           Progress,
	"IMPORT_COMPLETE":                              Ok,
	"IMPORT_ROLLBACK_IN_PROGRESS":                  Progress,
	"IMPORT_ROLLBACK_FAILED":                       Fail,
	"IMPORT_ROLLBACK_COMPLETE":                     Ok,
}

// StatusType returns the broad status type
func StatusType(s string) int {
	return statusMap[s]
}
