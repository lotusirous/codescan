package core

type (
	// User represents a user in the system.
	User struct {
		Name    string `json:"username"`
		Created int64  `json:"created"`
		Updated int64  `json:"updated"`
	}
)
