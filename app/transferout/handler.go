package transferout

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
// @Tags Transferout
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req query request.PageTransferout false "url query string"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/transfer-out [get]
func (h Handler) Page(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.PageTransferout)
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
// @Tags Transferout
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body request.CreateTransferout true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/transfer-out [post]
func (h Handler) Create(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	req := new(request.CreateTransferout)
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

// GetById
// @Tags Transferout
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Query preloads query string false "preloads"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/transfer-out/{id} [get]
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
// @Tags Transferout
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param req body request.UpdateTransferout true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/transfer-out [post]
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

	req := new(request.UpdateTransferout)
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
// @Tags Transferout
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/transfer-out/{id} [delete]
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

// SetInTransit
// @Tags Transferout
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/transfer-out/{id}/set-in-transit [put]
func (h Handler) SetInTransit(c echo.Context) error {
	var err error

	loginUser, err := jwt.GetUserLoginInfo(c)
	if err != nil {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetUserInfo, err, nil).SendJSON(c)
	}

	id := c.Param("id")
	if id == "" {
		return response.Error(http.StatusBadRequest, response.ErrorHandlerGetParam, err, nil).SendJSON(c)
	}

	err = h.usecase.SetInTransit(loginUser, id)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	return response.Success(http.StatusOK, response.SuccessHandler, nil).SendJSON(c)
}

// GenerateDeliveryOrder
// @Tags Transferout
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /stockmovementvehicle/transfer-out/{id}/generate-delivery-order [get]
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
		return response.Error(http.StatusBadRequest, err.Error(), err, nil).SendJSON(c)
	}

	fmt.Print(fmt.Sprintf("Delivery Order %s %s.pdf", vStockmovementvehicle.ID, utils.DisplayDate(time.Now())))

	filename := fmt.Sprintf("Delivery Order %s %s.pdf", vStockmovementvehicle.ID, utils.DisplayDate(time.Now()))
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)

	// Kirimkan PDF sebagai respons
	return c.Blob(http.StatusOK, "application/pdf", pdfBytes)
}
