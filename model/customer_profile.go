package model

type CustomerProfile struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	UserStatus  string   `json:"user_status"`
	ClientNo    string   `json:"client_no"`
	DeviceID    string   `json:"device_id"`
	SessionID   string   `json:"session_id"`
	TokenAccess string   `json:"token_access"`
	OtpType     string   `json:"otp_type"`
	FullName    string   `json:"full_name"`
	Email       string   `json:"email"`
	Mobile      string   `json:"mobile"`
	Dob         string   `json:"dob"`
	DobStr      string   `json:"dob_str"`
	HasPasscode string   `json:"has_passcode"`
	EncryptKey  string   `json:"encrypt_key"`
	Avatar      string   `json:"avatar"`
	LoginAt     string   `json:"login_at"`
	LoginAtStr  string   `json:"login_at_str"`
	Gender      string   `json:"gender"`
	Language    string   `json:"language"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
}
