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
}

type BaseResponse struct {
	Success bool                       `json:"success"`
	Data    []types.SensorDataResponse `json:"data"`
	Error   string                     `json:"error,omitempty"`
}

func NewAPIServer(usecase usecase.SensorSender) *apiServer {
	return &apiServer{
		usecase: usecase,
	}
}

func (aps *apiServer) GetSensorData(c echo.Context) error {
	params := c.QueryParams()

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
		startDateParsed, err := pkg.ParseDateWithFallback(startDateStr, "2006-01-02", " 00:00:00")
		if err != nil {
			return c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Error:   "Please provide valid start_date",
			})
		}
		startDate = startDateParsed
	}

	endDateStr := params.Get("end_date")
	if endDateStr != "" {
		endDateParsed, err := pkg.ParseDateWithFallback(endDateStr, "2006-01-02", " 23:59:59")
		if err != nil {
			return c.JSON(http.StatusBadRequest, BaseResponse{
				Success: false,
				Error:   "Please provide valid end_date",
			})
		}
		endDate = endDateParsed
	}

	data, err := aps.usecase.GetSensorData(c.Request().Context(), id1, id2, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, BaseResponse{
			Success: false,
			Error:   "Opps, something's wrong",
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Success: true,
		Data:    data,
	})
}

func (aps *apiServer) StartServer(cfg *config.Config) error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("v1/api")

	v1.GET("/sensors", aps.GetSensorData)

	address := fmt.Sprintf("%s:%s", cfg.API_HOST, cfg.API_PORT)
	log.Printf("API server is running on %s", address)
	return e.Start(address)
}
