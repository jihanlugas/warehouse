package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jihanlugas/warehouse/app/auditlog"
	"github.com/jihanlugas/warehouse/config"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/response"
	"github.com/jihanlugas/warehouse/utils"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type Handler struct {
	usecase         Usecase
	auditlogUsecase auditlog.Usecase
}

func NewHandler(usecase Usecase, auditlogUsecase auditlog.Usecase) Handler {
	return Handler{
		usecase:         usecase,
		auditlogUsecase: auditlogUsecase,
	}
}

// SignIn
// @Tags Authentication
// @Accept json
// @Produce json
// @Param req body request.Signin true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /auth/sign-in [post]
func (h Handler) SignIn(c echo.Context) error {
	var err error

	req := new(request.Signin)
	if err = c.Bind(req); err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	token, userLogin, err := h.usecase.SignIn(*req)
	if err != nil {
		//go h.auditlogUsecase.CreateAuditlogFailed(userLogin, "Login", err.Error(), req, nil)
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	go h.auditlogUsecase.CreateAuditlog(userLogin, model.AuditlogTypeSuccess, request.CreateAuditlog{
		StockmovementvehicleID: "",
		Title:                  fmt.Sprintf("Login"),
		Description:            fmt.Sprintf("Success login"),
		Request:                nil,
		Response: response.Payload{
			"token":     token,
			"userLogin": userLogin,
		},
	})

	return response.Success(http.StatusOK, response.SuccessHandler, response.Payload{
		"token":     token,
		"userLogin": userLogin,
	}).SendJSON(c)
}

// SignOut Sign out user
// @Tags Authentication
// @Accept json
// @Produce json
// // @Param req body request.Signout true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /auth/sign-out [get]
func (h Handler) SignOut(c echo.Context) error {
	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// RefreshToken
// @Tags Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /auth/refresh-token [get]
func (h Handler) RefreshToken(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	token, err := h.usecase.RefreshToken(loginUser)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, response.Payload{
		"token": token,
	}).SendJSON(c)
}

// Init
// @Tags Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /auth/init [get]
func (h Handler) Init(c echo.Context) error {
	var err error
	var res response.Init

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	vUser, vWarehouse, err := h.usecase.Init(loginUser)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	if vUser.UserRole == model.UserRoleOperator {
		res = response.Init{
			User:      vUser,
			Warehouse: &vWarehouse,
		}
	} else {
		res = response.Init{
			User: vUser,
		}
	}

	return response.Success(http.StatusOK, response.SuccessHandler, res).SendJSON(c)
}

func (h Handler) GoogleSignIn(c echo.Context) (err error) {
	url := googleOAuthConfig.AuthCodeURL("", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return c.Redirect(http.StatusFound, url)
}

func (h Handler) GoogleCallback(c echo.Context) (err error) {
	rawState := c.QueryParam("state")
	if rawState == "" {
		return h.googleSigninCallback(c)
	} else {
		return h.googleLinkCallback(c)
	}
}

func (h Handler) GoogleUnlink(c echo.Context) (err error) {
	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	err = h.usecase.GoogleUnlink(loginUser)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

func (h Handler) GoogleLink(c echo.Context) (err error) {
	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	stateJSON, _ := json.Marshal(loginUser)
	state := base64.URLEncoding.EncodeToString(stateJSON)

	url := googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return c.Redirect(http.StatusFound, url)
}

func (h Handler) googleSigninCallback(c echo.Context) (err error) {
	var redirectURL string
	code := c.QueryParam("code")
	if code == "" {
		redirectURL = fmt.Sprintf("%s/callback?state=%s&status=failed&message=%s", config.OauthFeCallback, "sign-in", err.Error())
		return c.Redirect(http.StatusFound, redirectURL)
	}

	providerToken, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		redirectURL = fmt.Sprintf("%s/callback?state=%s&status=failed&message=%s", config.OauthFeCallback, "sign-in", err.Error())
		return c.Redirect(http.StatusFound, redirectURL)
	}

	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+providerToken.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		redirectURL = fmt.Sprintf("%s/callback?state=%s&status=failed&message=%s", config.OauthFeCallback, "sign-in", err.Error())
		return c.Redirect(http.StatusFound, redirectURL)
	}
	defer resp.Body.Close()

	var callback response.GoogleCallback
	err = json.NewDecoder(resp.Body).Decode(&callback)
	if err != nil {
		redirectURL = fmt.Sprintf("%s/callback?state=%s&status=failed&message=%s", config.OauthFeCallback, "sign-in", err.Error())
		return c.Redirect(http.StatusFound, redirectURL)
	}

	token, loginUser, err := h.usecase.GoogleCallback(callback)
	if err != nil {
		redirectURL = fmt.Sprintf("%s/callback?state=%s&status=failed&message=%s", config.OauthFeCallback, "sign-in", err.Error())
		return c.Redirect(http.StatusFound, redirectURL)
	}

	go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeSuccess, request.CreateAuditlog{
		StockmovementvehicleID: "",
		Title:                  fmt.Sprintf("Login"),
		Description:            fmt.Sprintf("Success login with SSO"),
		Request:                nil,
		Response: response.Payload{
			"token":     token,
			"userLogin": loginUser,
		},
	})

	redirectURL = fmt.Sprintf("%s/callback?state=%s&status=success&token=%s&role=%s", config.OauthFeCallback, "sign-in", token, loginUser.UserRole)
	return c.Redirect(http.StatusFound, redirectURL)
}

func (h Handler) googleLinkCallback(c echo.Context) (err error) {
	var redirectURL string
	code := c.QueryParam("code")
	if code == "" {
		redirectURL = fmt.Sprintf("%s/callback?state=%s&status=failed&message=%s", config.OauthFeCallback, "link", err.Error())
		return c.Redirect(http.StatusFound, redirectURL)
	}

	providerToken, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		redirectURL = fmt.Sprintf("%s/callback?state=%s&status=failed&message=%s", config.OauthFeCallback, "link", err.Error())
		return c.Redirect(http.StatusFound, redirectURL)
	}

	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+providerToken.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		redirectURL = fmt.Sprintf("%s/callback?state=%s&status=failed&message=%s", config.OauthFeCallback, "link", err.Error())
		return c.Redirect(http.StatusFound, redirectURL)
	}
	defer resp.Body.Close()

	var callback response.GoogleCallback
	err = json.NewDecoder(resp.Body).Decode(&callback)
	if err != nil {
		redirectURL = fmt.Sprintf("%s/callback?state=%s&status=failed&message=%s", config.OauthFeCallback, "link", err.Error())
		return c.Redirect(http.StatusFound, redirectURL)
	}

	rawState := c.QueryParam("state")

	decoded, err := base64.URLEncoding.DecodeString(rawState)
	if err != nil {
		redirectURL = fmt.Sprintf("%s/callback?state=%s&status=failed&message=%s", config.OauthFeCallback, "link", err.Error())
		return c.Redirect(http.StatusFound, redirectURL)
	}
	var loginUser jwt.UserLogin
	err = json.Unmarshal(decoded, &loginUser)
	if err != nil {
		redirectURL = fmt.Sprintf("%s/callback?state=%s&status=failed&message=%s", config.OauthFeCallback, "link", err.Error())
		return c.Redirect(http.StatusFound, redirectURL)
	}

	err = h.usecase.GoogleLinkCallback(loginUser, callback)
	if err != nil {
		redirectURL = fmt.Sprintf("%s/callback?state=%s&status=failed&message=%s", config.OauthFeCallback, "link", err.Error())
		return c.Redirect(http.StatusFound, redirectURL)
	}

	redirectURL = fmt.Sprintf("%s/callback?state=%s&status=success&message=%s", config.OauthFeCallback, "link", "success link akun")
	return c.Redirect(http.StatusFound, redirectURL)
}
