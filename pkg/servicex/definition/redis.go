package def

const (
	RedisKeySessionToken      = "SessionToken:%s:%s"       // SessionToken:BJMSessionToken:xxx
	RedisKeyMutexSessionToken = "SessionToken:Mutex:%s:%s" // SessionToken:Mutex:BJMSessionToken:xxx
)
