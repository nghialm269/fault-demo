package users

type ServiceGetUserParams struct {
	ID int64
}

type ServiceGetUserResult struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type RepositoryGetUserByIDParams struct {
	ID int64
}

type RepositoryGetUserByIDResult struct {
	ID       int64
	Username string
}
