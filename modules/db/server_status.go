package db

type Status int

const (
	OPERATIONAL Status = iota
	SHUNNED
	OFF
)

func (s Status) String() string {
	strStatuses := [...]string{"OPERATIONAL", "SHUNNED", "OFF"}
	if s < OPERATIONAL || s > OFF {
		return "Unknown"
	}
	return strStatuses[s]
}
