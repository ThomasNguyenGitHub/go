// Package auth provides generic authentication functionality.
//
// Note this is not 100% secure and should only be used for prototyping,
// not for production systems or systems that are accessed by real users.
package auth

import (
	"context"
	"fmt"
	"github.com/ThomasNguyenGitHub/go/cache"
	"github.com/ThomasNguyenGitHub/go/errors"
	ldap "github.com/ThomasNguyenGitHub/go/ldap"
	"github.com/ThomasNguyenGitHub/go/storage/local"
	httpTransport "github.com/ThomasNguyenGitHub/go/transport/http"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	authCache cache.Cacher // cache is used for storing authentication tokens.
)

const (
	ADMIN = "admin"
	STAFF = "staff"
)

// admin is a staff
type User struct {
	Id                   string    `json:"id"`
	FullName             string    `json:"full_name"`
	UserName             string    `json:"user_name"`
	TaxCode              string    `json:"tax_code"`
	BankAccount          string    `json:"bank_account"`
	BankCode             string    `json:"bank_code"`
	Phone                string    `json:"phone_num"`
	Email                string    `json:"email"`
	Gender               string    `json:"gender"`
	Avatar               string    `json:"avatar"`
	Status               string    `json:"status"`
	IdcardNum            string    `json:"idcard_num"`
	ShowroomId           string    `json:"showroom_id"`
	IdcardIssueDate      string    `json:"idcard_issue_date"`
	IdcardIssuePlaceName string    `json:"idcard_issue_place_name"`
	IdcardIssuePlaceCode string    `json:"idcard_issue_place_code"`
	DateOfBirth          string    `json:"idcard_issue_date"`
	CreateAt             time.Time `json:"-"`
	UpdateAt             time.Time `json:"-"`
}

type OMSInfo struct {
	ID                   string `json:"id"`
	Loginname            string `json:"loginName"`
	Fullname             string `json:"fullName"`
	Fullnameen           string `json:"fullNameEn"`
	Email                string `json:"email"`
	Code                 string `json:"code"`
	Address              string `json:"address"`
	Addressid            string `json:"addressId"`
	Mobile               string `json:"mobile"`
	Status               string `json:"status"`
	Ext                  string `json:"ext"`
	Imageurl             string `json:"imageUrl"`
	Joiningdate          string `json:"joiningDate"`
	StartDate            string `json:"contractStartDate"`
	Gender               string `json:"gender"`
	Birthday             string `json:"birthDay"`
	LineManagerLoginName string `json:"lineManagerLoginName"`
	LineManagerEmail     string `json:"lineManagerEmail"`
	Deptid               string `json:"deptId"`
	Deptcode             string `json:"deptCode"`
	Idno                 string `json:"idNo"`
	Idissuedate          string `json:"idIssueDate"`
	Idissueplace         string `json:"idIssuePlace"`
	Nationality          string `json:"nationality"`
	Navaddress           string `json:"navAddress"`
	Curaddress           string `json:"curAddress"`
	Addisplayname        string `json:"adDisplayName"`
	Positions            []struct {
		ID                     string `json:"id"`
		Code                   string `json:"code"`
		Nameen                 string `json:"nameEn"`
		Namevn                 string `json:"nameVn"`
		Linemanagerposcode     string `json:"lineManagerPosCode"`
		Linemanagerloginname   string `json:"lineManagerLoginName"`
		Linemanagerempcode     string `json:"lineManagerEmpCode"`
		Linemanagerempfullname string `json:"lineManagerEmpFullName"`
		JobCode                string `json:"jobCode"`
	} `json:"positions"`
}

// Authenticator defines an interface for an authenticate-able User.
type Authenticator interface {
	Identifier() string     // Identifier returns a unique reference to this user.
	HashedPassword() string // HashedPassword returns the user's password hash.
}

// SetCache sets the Cache to use for authentication tokens.
func SetCache(c cache.Cacher) {
	authCache = c
}

// HashPassword returns a hashed version of the plain-text password provided.
func HashPassword(plainText string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash[:]), nil
}

func ComparePassword(HashedPassword, plainTextPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(HashedPassword), []byte(plainTextPassword))
}

// Authenticate validates an Authenticator based on it's password hash and the plain-text
// password provided.
func Authenticate(a Authenticator, plainTextPassword string) (AuthenticationTokenPair, error) {
	err := ComparePassword(a.HashedPassword(), plainTextPassword)
	if err != nil {
		return AuthenticationTokenPair{}, err
	}

	// Generate and cache a new token pair for this session
	return GenerateAndStoreTokens(a)
}
func SignOut(token string) error {
	if err := authCache.Delete(getAccessTokenCacheKey(token)); err != nil {
		return err
	}
	if err := authCache.Delete(getProfileCacheKey(token)); err != nil {
		return err
	}
	return nil
}

// Refresh generates a new token pair for a given authenticator.
func Refresh(a Authenticator, refreshToken string) (AuthenticationTokenPair, error) {
	newTokens, err := GenerateAndStoreTokens(a)
	if err != nil {
		return AuthenticationTokenPair{}, err
	}

	// Clear the old tokens from the cache
	if err := clearCachedTokens(refreshToken); err != nil {
		return AuthenticationTokenPair{}, err
	}

	return newTokens, nil
}

func GetUserProfileCache(token string, userProfile interface{}) error {
	return authCache.GetMarshaled(getProfileCacheKey(token), userProfile)
}
func SetUserProfileCache(token string, userProfile interface{}) (interface{}, error) {
	expiresTime := time.Duration(3600)
	if getTime, err := strconv.Atoi(local.Getenv("TOKEN_EXPIRES_TIME")); err == nil {
		expiresTime = time.Duration(getTime)
	}
	save, err := authCache.PutMarshaled(getProfileCacheKey(token), userProfile)
	if err != nil {
		return nil, err
	} else if err := authCache.Expire(getProfileCacheKey(token), expiresTime); err != nil {
		return nil, err
	}
	return save, nil
}
func GetIdentifierForProfileKey(token string) (string, error) {
	return authCache.GetString(getProfileCacheKey(token))
}

// GetIdentifierForAccessToken returns a user's identifier, as returned by
// the Authenticator interface, if it exists in the cache.
//
// If the identifier does not exist, and empty string and error will be returned.
func GetIdentifierForAccessToken(a string) (string, error) {
	return authCache.GetString(getAccessTokenCacheKey(a))
}

// GetIdentifierForRefreshToken returns a user's identifier, as returned by
// the Authenticator interface, if it exists in the cahce.
//
// If the identifier does not exist, an empty string and error will be returned.
func GetIdentifierForRefreshToken(r string) (string, error) {
	return authCache.GetString(getRefreshTokenCacheKey(r))
}

// generateAndStoreTokens creates and caches a new AuthenticationTokenPair.
func GenerateAndStoreTokens(a Authenticator) (AuthenticationTokenPair, error) {
	t := GenerateToken()
	if err := cacheTokens(t, a); err != nil {
		return AuthenticationTokenPair{}, err
	}

	return t, nil
}

// cacheTokens stores an access token and refresh token pair for an authenticated User.
func cacheTokens(t AuthenticationTokenPair, a Authenticator) error {
	expiresTime := time.Duration(3600)
	if getTime, err := strconv.Atoi(local.Getenv("TOKEN_EXPIRES_TIME")); err == nil {
		expiresTime = time.Duration(getTime)
	}
	if _, err := authCache.PutString(getAccessTokenCacheKey(t.AccessToken), a.Identifier()); err != nil {
		return err
	}
	if err := authCache.Expire(getAccessTokenCacheKey(t.AccessToken), expiresTime); err != nil {
		return err
	}

	if _, err := authCache.PutString(getRefreshTokenCacheKey(t.RefreshToken), a.Identifier()); err != nil {
		return err
	}
	if err := authCache.Expire(getRefreshTokenCacheKey(t.RefreshToken), expiresTime); err != nil {
		return err
	}
	if _, err := authCache.PutString(getRefreshToAccessTokenCacheKey(t.RefreshToken), t.AccessToken); err != nil {
		return err
	}
	if err := authCache.Expire(getRefreshToAccessTokenCacheKey(t.RefreshToken), expiresTime*2); err != nil {
		return err
	}
	return nil
}

// getAccessTokenCacheKey returns the access token cache key.
func getAccessTokenCacheKey(accessToken string) string {
	return fmt.Sprintf("accessToken:%s", accessToken)
}

// getRefreshTokenCacheKey returns the refresh token cache key.
func getRefreshTokenCacheKey(refreshToken string) string {
	return fmt.Sprintf("refreshToken:%s", refreshToken)
}

// getRefreshToAccessTokenCacheKey returns the refresh -> access token cache key.
func getRefreshToAccessTokenCacheKey(refreshToken string) string {
	return fmt.Sprintf("refreshToAccessToken:%s", refreshToken)
}

func getProfileCacheKey(token string) string {
	return fmt.Sprintf("profile:%s", token)
}

// clearCachedTokens clears all tokens associated to a refresh token.
func clearCachedTokens(r string) error {
	if a, err := authCache.GetString(getRefreshToAccessTokenCacheKey(r)); err != nil {
		return err
	} else if err = authCache.Delete(getAccessTokenCacheKey(a)); err != nil {
		return err
	} else if err = authCache.Delete(getRefreshTokenCacheKey(r)); err != nil {
		return err
	} else if err = authCache.Delete(getRefreshToAccessTokenCacheKey(r)); err != nil {
		return err
	}

	return nil
}
func GetUserIdByRequestContext(ctx context.Context) (string, error) {
	token, ok := ctx.Value(local.Getenv("JWT_TOKEN_CONTENT_KEY")).(string)
	if ok == false {
		return "", errors.BadRequest("miss token")
	}
	userId, err := GetIdentifierForAccessToken(token)
	if err != nil || len(userId) < 12 {
		return "", errors.Unauthorized("Token not correct")
	}
	return userId, nil
}

///---------------------------LDAP---AUTH------------------------------

func LAPConfig() *Config {
	return &Config{
		Server:   local.Getenv("LDAP_URL_ADDRESS"),
		Port:     389,
		BaseDN:   "OU=Users,DC=example,DC=com",
		Security: SecurityNone,
	}
}

// Authenticate checks if the given credentials are valid, or returns an error if one occurred.
// username may be either the sAMAccountName or the userPrincipalName.
func LDAPAutht(config *Config, username, password string) (bool, error) {
	upn, err := config.UPN(username)
	if err != nil {
		return false, err
	}

	conn, err := config.Connect()
	if err != nil {
		return false, err
	}
	defer conn.Conn.Close()

	return conn.Bind(upn, password)
}

// AuthenticateExtended checks if the given credentials are valid, or returns an error if one occurred.
// username may be either the sAMAccountName or the userPrincipalName.
// entry is the *ldap.Entry that holds the DN and any request attributes of the user.
// If groups is non-empty, userGroups will hold which of those groups the user is a member of.
// groups can be a list of groups referenced by DN or cn and the format provided will be the format returned.
func AuthenticateExtended(config *Config, username, password string, attrs, groups []string) (status bool, entry *ldap.Entry, userGroups []string, err error) {
	upn, err := config.UPN(username)
	if err != nil {
		return false, nil, nil, err
	}

	conn, err := config.Connect()
	if err != nil {
		return false, nil, nil, err
	}
	defer conn.Conn.Close()

	//bind
	status, err = conn.Bind(upn, password)
	if err != nil {
		return false, nil, nil, err
	}
	if !status {
		return false, nil, nil, nil
	}

	//get entry
	entry, err = conn.GetAttributes("userPrincipalName", upn, attrs)
	if err != nil {
		return false, nil, nil, err
	}

	if len(groups) > 0 {
		//get all groups
		foundGroups, err := conn.Search(fmt.Sprintf("(member:%s:=%s)", LDAPMatchingRuleInChain, entry.DN), []string{""}, 1000)
		if err != nil {
			return false, nil, nil, err
		}

		for _, group := range groups {
			groupDN, err := conn.GroupDN(group)
			if err != nil {
				return false, nil, nil, err
			}

			for _, userGroup := range foundGroups {
				if userGroup.DN == groupDN {
					userGroups = append(userGroups, group)
					break
				}
			}
		}
	}

	return status, entry, userGroups, nil
}

func GetLDABProfile(ctx context.Context, rmId string) (*OMSInfo, error) {
	urls := fmt.Sprintf("%s/v2/employee/find?input=%s", local.Getenv("OMS_HOST_ADDR"), strings.ToLower(rmId))
	var data []OMSInfo
	req, err := http.NewRequest(http.MethodGet, urls, nil)
	if err != nil {
		return nil, errors.InternalServerError(err)
	}
	req.Header.Add("Authorization", local.Getenv("OMS_API_BASIC_TOKEN"))
	err = httpTransport.ExecuteHTTP(ctx, req, &data)
	if err != nil {
		return nil, errors.InternalServerError(err)
	}
	return &data[0], nil
}
