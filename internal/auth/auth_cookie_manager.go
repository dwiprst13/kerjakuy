package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CookieManager interface {
	SetTokens(c *gin.Context, tokens AuthTokens)
	ClearTokens(c *gin.Context)
}

type CookieOptions struct {
	AccessCookieName  string
	RefreshCookieName string
	Domain            string
	Path              string
	Secure            bool
	HTTPOnly          bool
	SameSite          http.SameSite
	AccessTTL         time.Duration
	RefreshTTL        time.Duration
}

type cookieManager struct {
	opts CookieOptions
}

func NewCookieManager(opts CookieOptions) CookieManager {
	defaults := CookieOptions{
		AccessCookieName:  "kerjakuy_access",
		RefreshCookieName: "kerjakuy_refresh",
		Path:              "/",
		HTTPOnly:          true,
		SameSite:          http.SameSiteLaxMode,
		AccessTTL:         15 * time.Minute,
		RefreshTTL:        7 * 24 * time.Hour,
	}

	if opts.AccessCookieName == "" {
		opts.AccessCookieName = defaults.AccessCookieName
	}
	if opts.RefreshCookieName == "" {
		opts.RefreshCookieName = defaults.RefreshCookieName
	}
	if opts.Path == "" {
		opts.Path = defaults.Path
	}
	if opts.SameSite == 0 {
		opts.SameSite = defaults.SameSite
	}
	if opts.AccessTTL == 0 {
		opts.AccessTTL = defaults.AccessTTL
	}
	if opts.RefreshTTL == 0 {
		opts.RefreshTTL = defaults.RefreshTTL
	}
	if opts.HTTPOnly == false && !opts.Secure {
		opts.HTTPOnly = defaults.HTTPOnly
	}

	return &cookieManager{opts: opts}
}

func (m *cookieManager) SetTokens(c *gin.Context, tokens AuthTokens) {
	accessMaxAge := durationOrSeconds(m.opts.AccessTTL, tokens.ExpiresIn)
	refreshMaxAge := int(m.opts.RefreshTTL.Seconds())

	m.setCookie(c, m.opts.AccessCookieName, tokens.AccessToken, accessMaxAge)
	m.setCookie(c, m.opts.RefreshCookieName, tokens.RefreshToken, refreshMaxAge)
}

func (m *cookieManager) ClearTokens(c *gin.Context) {
	m.setCookie(c, m.opts.AccessCookieName, "", -1)
	m.setCookie(c, m.opts.RefreshCookieName, "", -1)
}

func (m *cookieManager) setCookie(c *gin.Context, name, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     m.opts.Path,
		Domain:   m.opts.Domain,
		Secure:   m.opts.Secure,
		HttpOnly: m.opts.HTTPOnly,
		SameSite: m.opts.SameSite,
	}

	if maxAge >= 0 {
		cookie.MaxAge = maxAge
		cookie.Expires = time.Now().Add(time.Duration(maxAge) * time.Second)
	} else {
		cookie.MaxAge = -1
		cookie.Expires = time.Unix(0, 0)
	}

	http.SetCookie(c.Writer, cookie)
}

func durationOrSeconds(d time.Duration, seconds int64) int {
	if d > 0 {
		return int(d.Seconds())
	}
	if seconds > 0 {
		return int(seconds)
	}
	return 0
}
