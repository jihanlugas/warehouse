package user

import (
	"net/http"
	"strings"

	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/response"
	"github.com/jihanlugas/warehouse/utils"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) Handler {
	return Handler{
		usecase: usecase,
	}
}

// Page
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req query request.PageUser false "url query string"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /user [get]
func (h Handler) Page(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.PageUser)
	if err = c.Bind(req); err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	if req.WarehouseID == "" {
		req.WarehouseID = loginUser.WarehouseID
	} else {
		if jwt.IsSaveWarehouseIDOR(loginUser, req.WarehouseID) {
			return response.Error(http.StatusBadRequest, response.ErrorHandlerIDOR, err, nil).SendJSON(c)
		}
	}

	data, count, err := h.usecase.Page(loginUser, *req)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, response.PayloadPagination(req, data, count)).SendJSON(c)
}

// GetById
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Query preloads query string false "preloads"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /user/{id} [get]
func (h Handler) GetById(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	id := c.Param("id")
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	preloads := c.QueryParam("preloads")
	var preloadSlice []string
	if preloads != "" {
		preloadSlice = strings.Split(preloads, ",")
	}

	vUser, err := h.usecase.GetById(loginUser, id, preloadSlice...)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, vUser).SendJSON(c)
}

// Create
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.CreateUser true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /user [post]
func (h Handler) Create(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.CreateUser)
	if err = c.Bind(req); err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, req.WarehouseID) {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerIDOR, err, nil).SendJSON(c)
	}

	err = h.usecase.Create(loginUser, *req)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// Update
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param req body request.UpdateUser true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /user/{id} [put]
func (h Handler) Update(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	id := c.Param("id")
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	req := new(request.UpdateUser)
	if err = c.Bind(req); err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	err = h.usecase.Update(loginUser, id, *req)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// UpdateUserPrivilege
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param req body request.UpdateUserprivilege true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /user/{id}/privilege [post]
func (h Handler) UpdateUserPrivilege(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	id := c.Param("id")
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	req := new(request.UpdateUserprivilege)
	if err = c.Bind(req); err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	err = h.usecase.UpdateUserprivilege(loginUser, id, *req)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// ChangePassword
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.ChangePassword true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /user/change-password [post]
func (h Handler) ChangePassword(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.ChangePassword)
	if err = c.Bind(req); err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	err = h.usecase.ChangePassword(loginUser, *req)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// Delete
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /user/{id} [delete]
func (h Handler) Delete(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	id := c.Param("id")
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	err = h.usecase.Delete(loginUser, id)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}
