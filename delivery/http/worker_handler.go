package http

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"net/http"
	"prakarsa-app/domain"
	"prakarsa-app/transport/request"
	"prakarsa-app/utils"
)

type ReferenceHandler struct {
	ReferenceUC domain.WorkerUsecase
}

// NewReferenceHandler will initialize the todo resources endpoint
func NewReferenceHandler(e *echo.Echo, referenceUC domain.WorkerUsecase) {
	handler := &ReferenceHandler{
		ReferenceUC: referenceUC,
	}

	apiV1 := e.Group("/worker/v1")
	apiV1.GET("/references/report-reasons", handler.GetReportReason)
}

func (h *ReferenceHandler) GetReportReason(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.GetReportReasonReferenceReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if res, meta, err := h.ReferenceUC.GetReportReason(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Reference successfully retrieved",
			"data":    res,
			"meta":    meta,
		})
	}
}
