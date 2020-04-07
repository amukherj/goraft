package state

type ServerState struct {
	state byte
}

func (ss *ServerState) String() string {
	switch ss.state {
	case 1:
		return "Follower"
	case 2:
		return "Candidate"
	case 3:
		return "Leader"
	default:
		return "(invalid)"
	}
}
