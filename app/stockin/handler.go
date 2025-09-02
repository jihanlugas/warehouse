package stockin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jihanlugas/warehouse/app/auditlog"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/response"
	"github.com/jihanlugas/warehouse/utils"
	"github.com/labstack/echo/v4"
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

// Page
// @Tags Stockin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req query request.PageStockin false "url query string"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/stock-in [get]
func (h Handler) Page(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.PageStockin)
	if err = c.Bind(req); err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	data, count, err := h.usecase.Page(loginUser, *req)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, response.PayloadPagination(req, data, count)).SendJSON(c)
}

// Create
// @Tags Stockin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.CreateStockin true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/stock-in [post]
func (h Handler) Create(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.CreateStockin)
	if err = c.Bind(req); err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	vStockmovementvehicle, err := h.usecase.Create(loginUser, *req)
	if err != nil {
		go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeFailed, request.CreateAuditlog{
			StockmovementvehicleID: vStockmovementvehicle.ID,
			Title:                  fmt.Sprintf("Buat Stock Masuk"),
			Description:            strings.TrimSpace(fmt.Sprintf("%s %s", vStockmovementvehicle.Number, err.Error())),
			Request:                req,
			Response:               nil,
		})
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeSuccess, request.CreateAuditlog{
		StockmovementvehicleID: vStockmovementvehicle.ID,
		Title:                  fmt.Sprintf("Buat Stock Masuk"),
		Description:            fmt.Sprintf("Buat Stock Masuk %s", vStockmovementvehicle.Number),
		Request:                req,
		Response:               nil,
	})

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// GetById
// @Tags Stockin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Query preloads query string false "preloads"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/stock-in/{id} [get]
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

	vStockmovementvehicle, err := h.usecase.GetById(loginUser, id, preloadSlice...)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, vStockmovementvehicle).SendJSON(c)
}

// Delete
// @Tags Stockin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/stock-in/{id} [delete]
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

	vStockmovementvehicle, err := h.usecase.Delete(loginUser, id)
	if err != nil {
		go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeFailed, request.CreateAuditlog{
			StockmovementvehicleID: vStockmovementvehicle.ID,
			Title:                  fmt.Sprintf("Hapus Stock Masuk"),
			Description:            strings.TrimSpace(fmt.Sprintf("%s %s", vStockmovementvehicle.Number, err.Error())),
			Request:                nil,
			Response:               nil,
		})
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}
	go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeSuccess, request.CreateAuditlog{
		StockmovementvehicleID: vStockmovementvehicle.ID,
		Title:                  fmt.Sprintf("Hapus Stock Masuk"),
		Description:            fmt.Sprintf("Hapus Stock Masuk %s", vStockmovementvehicle.Number),
		Request:                nil,
		Response:               nil,
	})
	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// SetComplete
// @Tags Stockin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/stock-in/{id}/set-complete [put]
func (h Handler) SetComplete(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	id := c.Param("id")
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	vStockmovementvehicle, err := h.usecase.SetComplete(loginUser, id)
	if err != nil {
		go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeFailed, request.CreateAuditlog{
			StockmovementvehicleID: vStockmovementvehicle.ID,
			Title:                  fmt.Sprintf("Set Complete Stock Masuk"),
			Description:            strings.TrimSpace(fmt.Sprintf("%s %s", vStockmovementvehicle.Number, err.Error())),
			Request:                nil,
			Response:               nil,
		})
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}
	go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeSuccess, request.CreateAuditlog{
		StockmovementvehicleID: vStockmovementvehicle.ID,
		Title:                  fmt.Sprintf("Set Complete Stock Masuk"),
		Description:            fmt.Sprintf("Set Complete Stock Masuk %s", vStockmovementvehicle.Number),
		Request:                nil,
		Response:               nil,
	})
	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}
