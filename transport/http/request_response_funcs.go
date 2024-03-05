package http

import (
	"context"
	"github.com/ThomasNguyenGitHub/go/util"
	"net/http"
)

// RequestFunc may take information from an HTTP request and put it into a
// request context. In Servers, RequestFuncs are executed prior to invoking the
// endpoint. In Clients, RequestFuncs are executed after creating the request
// but prior to invoking the HTTP client.
type RequestFunc func(context.Context, *http.Request) context.Context

// ServerResponseFunc may take information from a request context and use it to
// manipulate a ResponseWriter. ServerResponseFuncs are only executed in
// servers, after invoking the endpoint but prior to writing a response.
type ServerResponseFunc func(context.Context, http.ResponseWriter) context.Context

// ClientResponseFunc may take information from an HTTP request and make the
// response available for consumption. ClientResponseFuncs are only executed in
// clients, after a request has been made, but prior to it being decoded.
type ClientResponseFunc func(context.Context, *http.Response) context.Context

// SetContentType returns a ServerResponseFunc that sets the Content-Type header
// to the provided value.
func SetContentType(contentType string) ServerResponseFunc {
	return SetResponseHeader("Content-Type", contentType)
}

// SetResponseHeader returns a ServerResponseFunc that sets the given header.
func SetResponseHeader(key, val string) ServerResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter) context.Context {
		w.Header().Set(key, val)
		return ctx
	}
}

// SetRequestHeader returns a RequestFunc that sets the given header.
func SetRequestHeader(key, val string) RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		r.Header.Set(key, val)
		return ctx
	}
}

// PopulateRequestContext is a RequestFunc that populates several values into
// the context from the HTTP request. Those values may be extracted using the
// corresponding ContextKey type in this package.
func PopulateRequestContext(ctx context.Context, r *http.Request) context.Context {
	ipRequest, _ := util.GetIP(r)
	for k, v := range map[contextKey]string{
		ContextKeyRequestMethod:          r.Method,
		ContextKeyRequestURI:             r.RequestURI,
		ContextKeyRequestPath:            r.URL.Path,
		ContextKeyRequestProto:           r.Proto,
		ContextKeyRequestHost:            r.Host,
		ContextKeyRequestRemoteAddr:      r.RemoteAddr,
		ContextKeyRequestXForwardedFor:   r.Header.Get("X-Forwarded-For"),
		ContextKeyRequestXForwardedProto: r.Header.Get("X-Forwarded-Proto"),
		ContextKeyRequestAuthorization:   r.Header.Get("Authorization"),
		ContextKeyRequestReferer:         r.Header.Get("Referer"),
		ContextKeyRequestUserAgent:       r.Header.Get("User-Agent"),
		ContextKeyRequestXRequestID:      r.Header.Get("X-Request-Id"),
		ContextKeyRequestAccept:          r.Header.Get("Accept"),
		ContextKeyAccessToken:            r.Header.Get(HeaderAccessToken),
		ContextKeyAppID:                  r.Header.Get(HeaderAppID),
		ContextKeyAppKey:                 r.Header.Get(HeaderAppKey),
		ContextKeyAppVersion:             r.Header.Get(HeaderAppVersion),
		ContextKeyDeviceID:               r.Header.Get(HeaderDeviceID),
		ContextKeyDeviceModel:            r.Header.Get(HeaderDeviceModel),
		ContextKeyDeviceOSName:           r.Header.Get(HeaderDeviceOSName),
		ContextKeyLanguage:               r.Header.Get(HeaderLanguage),
		ContextKeyTimestamp:              r.Header.Get(HeaderTimestamp),
		ContextKeyUserID:                 r.Header.Get(HeaderUserID),
		ContextKeyUsername:               r.Header.Get(HeaderUsername),
		ContextKeyFullName:               r.Header.Get(HeaderFullName),
		ContextKeyMobile:                 r.Header.Get(HeaderMobile),
		ContextKeyEmail:                  r.Header.Get(HeaderEmail),
		ContextKeyClientNo:               r.Header.Get(HeaderClientNo),
		ContextKeyOtpType:                r.Header.Get(HeaderOtpType),
		ContextKeySessionID:              r.Header.Get(HeaderSessionID),
		ContextKeyEncryptKey:             r.Header.Get(HeaderEncryptKey),
		ContextKeyGender:                 r.Header.Get(HeaderGender),
		ContextKeyVendorUsers:            r.Header.Get(HeaderVendorUsers),
		ContextKeyClientInfo:             r.Header.Get(HeaderClientInfo),
		ContextKeyChannel:                r.Header.Get(HeaderChannel),
		ContextKeyOrganization:           r.Header.Get(HeaderOrganization),
		ContextKeyLoginChannel:           r.Header.Get(HeaderLoginChannel),
		ContextKeyActionBy:               r.Header.Get(HeaderActionBy),
		ContextKeyIsDummyUser:            r.Header.Get(HeaderIsDummyUser),
		ContextKeySmoDeviceId:            r.Header.Get(HeaderSmoDeviceId),
		ContextKeyDeviceOSVersion:        r.Header.Get(HeaderDeviceOSVersion),
		ContextKeyIPRequest:              ipRequest,
	} {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

type contextKey int

const (
	// ContextKeyRequestMethod is populated in the context by
	// PopulateRequestContext. Its value is r.Method.
	ContextKeyRequestMethod contextKey = iota

	// ContextKeyRequestURI is populated in the context by
	// PopulateRequestContext. Its value is r.RequestURI.
	ContextKeyRequestURI

	// ContextKeyRequestPath is populated in the context by
	// PopulateRequestContext. Its value is r.URL.Path.
	ContextKeyRequestPath

	// ContextKeyRequestProto is populated in the context by
	// PopulateRequestContext. Its value is r.Proto.
	ContextKeyRequestProto

	// ContextKeyRequestHost is populated in the context by
	// PopulateRequestContext. Its value is r.Host.
	ContextKeyRequestHost

	// ContextKeyRequestRemoteAddr is populated in the context by
	// PopulateRequestContext. Its value is r.RemoteAddr.
	ContextKeyRequestRemoteAddr

	// ContextKeyRequestXForwardedFor is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("X-Forwarded-For").
	ContextKeyRequestXForwardedFor

	// ContextKeyRequestXForwardedProto is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("X-Forwarded-Proto").
	ContextKeyRequestXForwardedProto

	// ContextKeyRequestAuthorization is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("Authorization").
	ContextKeyRequestAuthorization

	// ContextKeyRequestReferer is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("Referer").
	ContextKeyRequestReferer

	// ContextKeyRequestUserAgent is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("User-Agent").
	ContextKeyRequestUserAgent

	// ContextKeyRequestXRequestID is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("X-Request-Id").
	ContextKeyRequestXRequestID

	// ContextKeyRequestAccept is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("Accept").
	ContextKeyRequestAccept

	// ContextKeyResponseHeaders is populated in the context whenever a
	// ServerFinalizerFunc is specified. Its value is of type http.Header, and
	// is captured only once the entire response has been written.
	ContextKeyResponseHeaders

	// ContextKeyResponseSize is populated in the context whenever a
	// ServerFinalizerFunc is specified. Its value is of type int64.
	ContextKeyResponseSize

	// ContextKeyAppID is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("app_id").
	ContextKeyAppID

	// ContextKeyAppKey is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("app_key").
	ContextKeyAppKey

	// ContextKeyAppVersion is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("app_version").
	ContextKeyAppVersion

	// ContextKeyDeviceID is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("device_id").
	ContextKeyDeviceID

	// ContextKeyDeviceModel is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("device_model").
	ContextKeyDeviceModel

	// ContextKeyDeviceOSName is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("device_os_name").
	ContextKeyDeviceOSName

	// ContextKeyLanguage is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("language").
	ContextKeyLanguage

	// ContextKeyTimestamp is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("timestamp").
	ContextKeyTimestamp

	// ContextKeyUserID is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("user_id").
	ContextKeyUserID

	// ContextKeyUsername is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("username").
	ContextKeyUsername

	// ContextKeyFullName is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("full_name").
	ContextKeyFullName

	// ContextKeyMobile is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("mobile").
	ContextKeyMobile

	// ContextKeyEmail is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("email").
	ContextKeyEmail

	// ContextKeyClientNo is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("client_no").
	ContextKeyClientNo

	// ContextKeyIPRequest is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("ip_request").
	ContextKeyIPRequest

	// ContextKeyAccessToken is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("access_token").
	ContextKeyAccessToken

	// ContextKeyOtpType is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("otp_type").
	ContextKeyOtpType

	// ContextKeySessionID is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("session_id").
	ContextKeySessionID

	// ContextKeyEncryptKey is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("encrypt_key").
	ContextKeyEncryptKey

	// ContextKeyGender is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("gender").
	ContextKeyGender

	// ContextKeyVendorUsers is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("vendor_users").
	ContextKeyVendorUsers

	// ContextKeyClientInfo is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("client_info").
	ContextKeyClientInfo

	// ContextKeyChannel is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("channel").
	ContextKeyChannel

	// ContextKeyOrganization is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("organization").
	ContextKeyOrganization

	// ContextKeyLoginChannel is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("login_channel").
	ContextKeyLoginChannel

	// ContextKeyActionBy is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("action_by").
	ContextKeyActionBy

	// ContextKeyIsDummyUser is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("is_dummy_user").
	ContextKeyIsDummyUser

	// ContextKeySmoDeviceId is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("smo_device_id").
	ContextKeySmoDeviceId

	// ContextKeyDeviceOSVersion is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("device_os_version").
	ContextKeyDeviceOSVersion
)

const (
	ResponseCodeSuccess = "000000"
)
