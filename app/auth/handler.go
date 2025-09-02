package auth

import (
	"context"
	"encoding/json"
	"errors"
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
		Description:            fmt.Sprintf("%s Success login", userLogin.Fullname),
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
	// bisa tambahkan state random untuk CSRF protection
	state := config.OauthKey

	url := googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return c.Redirect(http.StatusFound, url)
}

func (h Handler) GoogleCallback(c echo.Context) (err error) {
	code := c.QueryParam("code")
	if code == "" {
		return response.Error(http.StatusBadRequest, "missing code", errors.New("missing code"), nil).SendJSON(c)
	}

	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return response.Error(http.StatusBadRequest, "token exchange failed: ", err, nil).SendJSON(c)
	}

	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, "fetch userinfo failed", err, nil).SendJSON(c)
	}
	defer resp.Body.Close()

	var user map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return response.Error(http.StatusBadRequest, "decode userinfo failed", err, nil).SendJSON(c)
	}

	// STEP 4: Redirect balik ke FE dengan token
	redirectURL := fmt.Sprintf("http://localhost:3000/callback?token=%s&role=ADMIN", "asdjn23asd")
	return c.Redirect(http.StatusFound, redirectURL)
}
