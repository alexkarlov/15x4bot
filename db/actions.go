package db

type Action struct {
	Id      int
	Command string
	Msg     string
}

func GetActionMsg(cmd string) (msg string, err error) {
	row := dbConn.QueryRow("SELECT msg FROM actions WHERE command=$1", cmd)
	err = row.Scan(&msg)
	return
}
