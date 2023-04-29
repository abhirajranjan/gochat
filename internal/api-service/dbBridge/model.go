package dbBridge

type userroles []string

func (ur userroles) Has(key string) bool {
	for _, k := range ur {
		if k == key {
			return true
		}
	}
	return false
}

type user struct {
	UserID    string
	UserRoles userroles
}

type status struct {
	Err     string
	ErrCode int64
}

type userData struct {
	User   user
	Status status
}

func (lr *userData) GetMap() (map[string]interface{}, error) {
	user := make(map[string]interface{})
	user["userID"] = lr.User.UserID
	user["userRoles"] = lr.User.UserRoles

	return user, nil
}

func (lr *userData) GetErr() string {
	return lr.Status.Err
}

func (lr *userData) GetErrCode() int64 {
	return lr.Status.ErrCode
}
