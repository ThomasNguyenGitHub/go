package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const (
	AuthorizationBasic    = "Basic"
	AuthorizationBearer   = "Bearer"
	HeaderAuthorization   = "Authorization"
	HeaderContentType     = "Content-Type"
	HeaderAccessToken     = "access_token"
	HeaderAppID           = "app_id"
	HeaderAppKey          = "app_key"
	HeaderChannel         = "channel"
	HeaderOrganization    = "organization"
	HeaderAppVersion      = "app_version"
	HeaderUserID          = "user_id"
	HeaderUsername        = "username"
	HeaderFullName        = "full_name"
	HeaderMobile          = "mobile"
	HeaderEmail           = "email"
	HeaderClientNo        = "client_no"
	HeaderDeviceID        = "device_id"
	HeaderDeviceModel     = "device_model"
	HeaderDeviceOSName    = "device_os_name"
	HeaderIPRequest       = "ip_request"
	HeaderLanguage        = "language"
	HeaderTimestamp       = "timestamp"
	HeaderOtpType         = "otp_type"
	HeaderSessionID       = "session_id"
	HeaderEncryptKey      = "encrypt_key"
	HeaderLatitude        = "lat"
	HeaderLongitude       = "lng"
	HeaderGender          = "gender"
	HeaderLoginChannel    = "login_channel"
	HeaderVendorUsers     = "vendor_users"
	HeaderClientInfo      = "client_info"
	HeaderActionBy        = "action_by"
	HeaderSignature       = "signature"
	HeaderLanguageVN      = "vi-VN"
	HeaderLanguageEN      = "en-US"
	HeaderIsDummyUser     = "is_dummy_user"
	HeaderSmoDeviceId     = "smo_device_id"
	HeaderDeviceOSVersion = "device_os_version"
)

const HeaderContentTypeValueApplicationJSON = "application/json"

type RequestContext struct {
	AccessToken     string
	AppID           string
	AppKey          string
	Channel         string
	Organization    string
	AppVersion      string
	UserID          string
	Username        string
	FullName        string
	Mobile          string
	Email           string
	ClientNo        string
	DeviceID        string
	DeviceModel     string
	DeviceOSName    string
	IPRequest       string
	Language        string
	Timestamp       string
	OtpType         string
	SessionID       string
	EncryptKey      string
	Gender          string
	LoginChannel    string
	VendorUsers     string
	ClientInfo      string
	XRequestId      string
	ActionBy        string
	IsDummyUser     string
	SmoDeviceId     string
	DeviceOSVersion string
}

type VendorUser struct {
	ID          string     `json:"vendor_id"`
	Type        string     `json:"type"`
	ExpiredDate *time.Time `json:"expired_date"`
}

func GetRequestContext(ctx context.Context) RequestContext {
	return RequestContext{
		AccessToken:     GetString(ctx, ContextKeyAccessToken),
		AppID:           GetString(ctx, ContextKeyAppID),
		AppKey:          GetString(ctx, ContextKeyAppKey),
		Channel:         GetString(ctx, ContextKeyChannel),
		Organization:    GetString(ctx, ContextKeyOrganization),
		AppVersion:      GetString(ctx, ContextKeyAppVersion),
		UserID:          GetString(ctx, ContextKeyUserID),
		Username:        GetString(ctx, ContextKeyUsername),
		FullName:        GetString(ctx, ContextKeyFullName),
		Mobile:          GetString(ctx, ContextKeyMobile),
		Email:           GetString(ctx, ContextKeyEmail),
		ClientNo:        GetString(ctx, ContextKeyClientNo),
		DeviceID:        GetString(ctx, ContextKeyDeviceID),
		DeviceModel:     GetString(ctx, ContextKeyDeviceModel),
		DeviceOSName:    GetString(ctx, ContextKeyDeviceOSName),
		IPRequest:       GetString(ctx, ContextKeyIPRequest),
		Language:        GetString(ctx, ContextKeyLanguage),
		Timestamp:       GetString(ctx, ContextKeyTimestamp),
		OtpType:         GetString(ctx, ContextKeyOtpType),
		SessionID:       GetString(ctx, ContextKeySessionID),
		EncryptKey:      GetString(ctx, ContextKeyEncryptKey),
		Gender:          GetString(ctx, ContextKeyGender),
		LoginChannel:    GetString(ctx, ContextKeyLoginChannel),
		VendorUsers:     GetString(ctx, ContextKeyVendorUsers),
		ClientInfo:      GetString(ctx, ContextKeyClientInfo),
		XRequestId:      GetString(ctx, ContextKeyRequestXRequestID),
		ActionBy:        GetString(ctx, ContextKeyActionBy),
		IsDummyUser:     GetString(ctx, ContextKeyIsDummyUser),
		SmoDeviceId:     GetString(ctx, ContextKeySmoDeviceId),
		DeviceOSVersion: GetString(ctx, ContextKeyDeviceOSVersion),
	}
}

func SetHeader(r *http.Request, key string, value ...string) {
	r.Header[key] = value
}

func GetHeader(r *http.Request, key string) string {
	if value := r.Header[key]; len(value) > 0 {
		return value[0]
	}
	return ""
}

func GetVendorUsers(value string) []VendorUser {
	var vendorUsers []VendorUser
	_ = json.Unmarshal([]byte(value), &vendorUsers)
	return vendorUsers
}

func GetString(ctx context.Context, key contextKey) string {
	if v, ok := ctx.Value(key).(string); ok {
		return v
	}
	return ""
}
