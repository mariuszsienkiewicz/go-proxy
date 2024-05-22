package db

type Status int

const (
	OPERATIONAL Status = iota // server is working in current moment
	SHUNNED                   // there is a problem with this server, should be ignored
	OFF                       // server is turned off and should be ignored too
)

func (s Status) String() string {
	strStatuses := [...]string{"OPERATIONAL", "SHUNNED", "OFF"}
	if s < OPERATIONAL || s > OFF {
		return "Unknown"
	}
	return strStatuses[s]
}
