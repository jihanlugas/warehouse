package stockmovementvehicleretail

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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
// @Tags StockmovementvehicleRetail
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req query request.PageStockmovementvehicleRetail false "url query string"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/retail [get]
func (h Handler) Page(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.PageStockmovementvehicleRetail)
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
// @Tags StockmovementvehicleRetail
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.CreateStockmovementvehicleRetail true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/retail [post]
func (h Handler) Create(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.CreateStockmovementvehicleRetail)
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
			Title:                  fmt.Sprintf("Buat Pengiriman Keluar Retail"),
			Description:            strings.TrimSpace(fmt.Sprintf("%s %s", vStockmovementvehicle.Number, err.Error())),
			Request:                req,
			Response:               nil,
		})
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeSuccess, request.CreateAuditlog{
		StockmovementvehicleID: vStockmovementvehicle.ID,
		Title:                  fmt.Sprintf("Buat Pengiriman Keluar Retail"),
		Description:            fmt.Sprintf("Buat Pengiriman Keluar Retail %s", vStockmovementvehicle.Number),
		Request:                req,
		Response:               nil,
	})

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// GetById
// @Tags StockmovementvehicleRetail
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Query preloads query string false "preloads"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/retail/{id} [get]
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

// Update
// @Tags StockmovementvehicleRetail
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param req body request.UpdateStockmovementvehicleRetail true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/retail [post]
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

	req := new(request.UpdateStockmovementvehicleRetail)
	if err = c.Bind(req); err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	vStockmovementvehicle, err := h.usecase.Update(loginUser, id, *req)
	if err != nil {
		go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeFailed, request.CreateAuditlog{
			StockmovementvehicleID: vStockmovementvehicle.ID,
			Title:                  fmt.Sprintf("Edit Loading Pengiriman Keluar Retail"),
			Description:            strings.TrimSpace(fmt.Sprintf("%s %s", vStockmovementvehicle.Number, err.Error())),
			Request:                req,
			Response:               nil,
		})
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeSuccess, request.CreateAuditlog{
		StockmovementvehicleID: vStockmovementvehicle.ID,
		Title:                  fmt.Sprintf("Edit Loading Pengiriman Keluar Retail"),
		Description:            fmt.Sprintf("Edit Loading Pengiriman Keluar Retail %s", vStockmovementvehicle.Number),
		Request:                req,
		Response:               nil,
	})

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)

}

// Delete
// @Tags StockmovementvehicleRetail
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/retail/{id} [delete]
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
			Title:                  fmt.Sprintf("Hapus Pengiriman Keluar Retail"),
			Description:            strings.TrimSpace(fmt.Sprintf("%s %s", vStockmovementvehicle.Number, err.Error())),
			Request:                nil,
			Response:               nil,
		})
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeSuccess, request.CreateAuditlog{
		StockmovementvehicleID: vStockmovementvehicle.ID,
		Title:                  fmt.Sprintf("Hapus Pengiriman Keluar Retail"),
		Description:            fmt.Sprintf("Hapus Pengiriman Keluar Retail %s", vStockmovementvehicle.Number),
		Request:                nil,
		Response:               nil,
	})

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// SetComplete
// @Tags StockmovementvehicleRetail
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/retail/{id}/set-complete [put]
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
			Title:                  fmt.Sprintf("Set In Transit Pengiriman Keluar Retail"),
			Description:            strings.TrimSpace(fmt.Sprintf("%s %s", vStockmovementvehicle.Number, err.Error())),
			Request:                nil,
			Response:               nil,
		})
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeSuccess, request.CreateAuditlog{
		StockmovementvehicleID: vStockmovementvehicle.ID,
		Title:                  fmt.Sprintf("Set In Transit Pengiriman Keluar Retail"),
		Description:            fmt.Sprintf("Set In Transit Pengiriman Keluar Retail %s", vStockmovementvehicle.Number),
		Request:                nil,
		Response:               nil,
	})

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// GenerateDeliveryOrder
// @Tags StockmovementvehicleRetail
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/retail/{id}/generate-delivery-order [get]
func (h Handler) GenerateDeliveryOrder(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	id := c.Param("id")
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	pdfBytes, vStockmovementvehicle, err := h.usecase.GenerateDeliveryOrder(loginUser, id)
	if err != nil {
		go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeFailed, request.CreateAuditlog{
			StockmovementvehicleID: vStockmovementvehicle.ID,
			Title:                  fmt.Sprintf("Generate Surat Jalan Pengiriman Keluar Retail"),
			Description:            strings.TrimSpace(fmt.Sprintf("%s %s", vStockmovementvehicle.Number, err.Error())),
			Request:                nil,
			Response:               nil,
		})
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	go h.auditlogUsecase.CreateAuditlog(loginUser, model.AuditlogTypeSuccess, request.CreateAuditlog{
		StockmovementvehicleID: vStockmovementvehicle.ID,
		Title:                  fmt.Sprintf("Generate Surat Jalan Pengiriman Keluar Retail"),
		Description:            fmt.Sprintf("Generate Surat Jalan Pengiriman Keluar Retail %s", vStockmovementvehicle.Number),
		Request:                nil,
		Response:               nil,
	})

	fmt.Print(fmt.Sprintf("Delivery Order %s %s.pdf", vStockmovementvehicle.ID, utils.DisplayDate(time.Now())))

	filename := fmt.Sprintf("Delivery Order %s %s.pdf", vStockmovementvehicle.ID, utils.DisplayDate(time.Now()))
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)

	// Kirimkan PDF sebagai respons
	return c.Blob(http.StatusOK, "application/pdf", pdfBytes)
}
