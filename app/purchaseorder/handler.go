package purchaseorder

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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
// @Tags Purchaseorder
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req query request.PagePurchaseorder false "url query string"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /purchase-order [get]
func (h Handler) Page(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.PagePurchaseorder)
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

// GetById
// @Tags Purchaseorder
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Query preloads query string false "preloads"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /purchase-order/{id} [get]
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

	vPurchaseorder, err := h.usecase.GetById(loginUser, id, preloadSlice...)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, vPurchaseorder).SendJSON(c)
}

// Create
// @Tags Purchaseorder
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.CreatePurchaseorder true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /purchase-order [post]
func (h Handler) Create(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.CreatePurchaseorder)
	if err = c.Bind(req); err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerBind, err, nil).SendJSON(c)
	}

	utils.TrimWhitespace(req)

	err = c.Validate(req)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerFailedValidation, err, response.ValidationError(err)).SendJSON(c)
	}

	err = h.usecase.Create(loginUser, *req)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// Update
// @Tags Purchaseorder
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param req body request.UpdatePurchaseorder true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /purchase-order/{id} [put]
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

	req := new(request.UpdatePurchaseorder)
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

// Delete
// @Tags Purchaseorder
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /purchase-order/{id} [delete]
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

// SetStatusOpen
// @Tags Purchaseorder
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /purchase-order/{id}/set-status-open [put]
func (h Handler) SetStatusOpen(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	id := c.Param("id")
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	err = h.usecase.SetStatusOpen(loginUser, id)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// SetStatusClose
// @Tags Purchaseorder
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /purchase-order/{id}/set-status-close [put]
func (h Handler) SetStatusClose(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	id := c.Param("id")
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	err = h.usecase.SetStatusClose(loginUser, id)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// GenerateInvoice
// @Tags Purchaseorder
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Query preloads query string false "preloads"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /purchase-order/{id}/generate-invoice [get]
func (h Handler) GenerateInvoice(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	id := c.Param("id")
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	pdfBytes, vPurchaseorder, err := h.usecase.GenerateInvoice(loginUser, id)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	fmt.Print(fmt.Sprintf("Purchaseorder Invoice %s %s.pdf", vPurchaseorder.ID, utils.DisplayDate(time.Now())))

	filename := fmt.Sprintf("Purchaseorder Invoice %s %s.pdf", vPurchaseorder.ID, utils.DisplayDate(time.Now()))
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)

	// Kirimkan PDF sebagai respons
	return c.Blob(http.StatusOK, "application/pdf", pdfBytes)
}
