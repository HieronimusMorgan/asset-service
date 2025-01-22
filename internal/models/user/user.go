package user

type User struct {
	UserID         uint   `json:"user_id,omitempty"`
	ClientID       string `json:"client_id,omitempty"`
	Username       string `json:"username,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	FullName       string `json:"full_name,omitempty"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	ProfilePicture string `json:"profile_picture,omitempty"`
	RoleID         uint   `json:"role_id,omitempty"`
}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUUID   string `json:"access_uuid"`
	RefreshUUID  string `json:"refresh_uuid"`
	AtExpires    int64  `json:"at_expires"`
	RtExpires    int64  `json:"rt_expires"`
}
