package store

// Action represents a table action, which needs for simple responses on a specific messages
type Action struct {
	ID      int
	Command string
	Msg     string
}

// ActionMsg returns message string from actions table
func ActionMsg(cmd string) (msg string, err error) {
	err = dbConn.QueryRow("SELECT msg FROM actions WHERE command=$1", cmd).Scan(&msg)
	return
}
