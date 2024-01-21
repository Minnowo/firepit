package data

const (
	// byte to indicate the thread should be running
	CHAN__RUNNING byte = iota

	// byte to indicate the thread should be paused
	CHAN__PAUSED byte = iota

	// byte to indicate the thread should die
	CHAN__DEAD byte = iota
)

// Gets the state as a human-readable string
func ChannelStateToString(state byte) string {
	switch state {
	case CHAN__DEAD:
		return "DEAD"
	case CHAN__PAUSED:
		return "PAUSED"
	case CHAN__RUNNING:
		return "RUNNING"
	}
	return "UNKNOWN"
}
