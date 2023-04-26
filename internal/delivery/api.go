package delivery

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/muhammadtaufan/go-sensor-collector/config"
	"github.com/muhammadtaufan/go-sensor-collector/internal/types"
	"github.com/muhammadtaufan/go-sensor-collector/internal/usecase"
	"github.com/muhammadtaufan/go-sensor-collector/pkg"
)

type apiServer struct {
	usecase usecase.SensorSender
	cfg     *config.Config
}

type BaseResponse struct {
	Success bool                       `json:"success"`
	Data    []types.SensorDataResponse `json:"data"`
	Error   string                     `json:"error,omitempty"`
}

func NewAPIServer(usecase usecase.SensorSender, cfg *config.Config) *apiServer {
	return &apiServer{
		usecase: usecase,
		cfg:     cfg,
	}
}

func (aps *apiServer) GetSensorData(c echo.Context) error {
	params := c.QueryParams()

	limitStr := params.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offsetStr := params.Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 1 {
		offset = 0
	}

	id1Str := params.Get("id1")
	id2Str := params.Get("id2")

	var id1 *string
	if id1Str != "" {
		id1 = &id1Str
	}

	var id2 *int
	if id2Str != "" {
		id2Int, err := strconv.Atoi(id2Str)
		if err != nil {
			return c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Error:   "Please provide valid id2",
			})
		}
		id2 = &id2Int
	}

	var startDate, endDate *time.Time

	startDateStr := params.Get("start_date")
	if startDateStr != "" {
		startDateParsed, err := pkg.ParseDateWithFallback(startDateStr, "2006-01-02", " 00:00:00", aps.cfg.TIMEZONE)
		if err != nil {
			return c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Error:   "Please provide valid start_date",
			})
		}
		startDateParsedUTC := startDateParsed.UTC()
		startDate = &startDateParsedUTC
	}

	endDateStr := params.Get("end_date")
	if endDateStr != "" {
		endDateParsed, err := pkg.ParseDateWithFallback(endDateStr, "2006-01-02", " 23:59:59", aps.cfg.TIMEZONE)
		if err != nil {
			return c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Error:   "Please provide valid end_date",
			})
		}
		startDateParsedUTC := endDateParsed.UTC()
		endDate = &startDateParsedUTC
	}

	data, err := aps.usecase.GetSensorData(c.Request().Context(), id1, id2, startDate, endDate, &limit, &offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Success: false,
			Error:   "Opps, something's wrong",
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Success: true,
		Data:    data,
	})
}

func (aps *apiServer) DeleteSensorData(c echo.Context) error {
	var requestBody types.SensorDataRequest
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	var id1 *string
	if requestBody.ID1 != "" {
		id1 = &requestBody.ID1
	}

	var id2 *int
	if requestBody.ID2 != "" {
		id2Int, err := strconv.Atoi(requestBody.ID2)
		if err != nil {
			return c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Error:   "Please provide valid id2",
			})
		}
		id2 = &id2Int
	}

	var startDate, endDate *time.Time

	if requestBody.StartDate != "" {
		startDateParsed, err := pkg.ParseDateWithFallback(requestBody.StartDate, "2006-01-02", " 00:00:00", aps.cfg.TIMEZONE)
		if err != nil {
			return c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Error:   "Please provide valid start_date",
			})
		}
		startDate = startDateParsed
	}

	if requestBody.EndDate != "" {
		endDateParsed, err := pkg.ParseDateWithFallback(requestBody.EndDate, "2006-01-02", " 23:59:59", aps.cfg.TIMEZONE)
		if err != nil {
			return c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Error:   "Please provide valid end_date",
			})
		}
		endDate = endDateParsed
	}

	if id1 == nil && id2 == nil && startDate == nil && endDate == nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Success: false,
			Error:   "Please provide valid request",
		})
	}

	err := aps.usecase.DeleteSensorData(c.Request().Context(), id1, id2, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, BaseResponse{
			Success: false,
			Error:   "Opps, something's wrong",
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Success: true,
	})
}

func (aps *apiServer) UpdateSensorData(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Success: false,
			Error:   "Please provide id",
		})
	}

	var requestBody types.UpdateSensorDataRequest
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	err := aps.usecase.UpdateSensorData(c.Request().Context(), id, &requestBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Success: false,
			Error:   "Opps, something's wrong",
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Success: true,
	})
}

func (aps *apiServer) StartServer(cfg *config.Config) error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("v1/api")

	v1.GET("/sensors", aps.GetSensorData)
	v1.DELETE("/sensors", aps.DeleteSensorData)
	v1.PATCH("/sensors/:id", aps.UpdateSensorData)

	address := fmt.Sprintf("%s:%s", cfg.API_HOST, cfg.API_PORT)
	log.Printf("API server is running on %s", address)
	return e.Start(address)
}
