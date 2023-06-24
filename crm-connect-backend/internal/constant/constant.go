package constant

const (
	/*
	 . ServiceName is the default service name
	 . with multiple deployables in picture, currently the service name is picked from config vars
	 . all usage of this const is removed from client setups
	 . only usage exist in default service test cases
	*/
	EmptyString  = ""
	NoneString   = "none"
	DefaultQueue = "unknown"
)

// Date formats
const (
	YYYYMMDD = "2006-01-02"
)

const (
	ReopenActionType     = "reopened"
	ReinitActionType     = "reinitiated"
	ActionTypePrefix     = "conversations-"
	ReasonPrefix         = "conversations."
	SetActiveKeyConstant = ".conversation.active."
)

const (
	DefaultResolutionReason = "manual"
)
const (
	RealTime = "real-time"
)

const (
	True  = "true"
	False = "false"
)

const (
	TCP = "tcp"
	All = "all"
)
