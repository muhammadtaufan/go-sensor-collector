package delivery

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/muhammadtaufan/go-sensor-collector/config"
	"github.com/muhammadtaufan/go-sensor-collector/internal/types"
	"github.com/muhammadtaufan/go-sensor-collector/internal/usecase"
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

func (aps *apiServer) GetDataByIDs(c echo.Context) error {
	params := c.QueryParams()
	id1 := params.Get("id1")

	if id1 == "" {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Success: false,
			Error:   "Please provide id1",
		})
	}

	id2Str := params.Get("id2")
	if id2Str == "" {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Success: false,
			Error:   "Please provide id2",
		})
	}
	id2, _ := strconv.Atoi(id2Str)

	data, err := aps.usecase.GetDataByIDs(c.Request().Context(), id1, id2)
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

	v1.GET("/sensors", aps.GetDataByIDs)

	address := fmt.Sprintf("%s:%s", cfg.API_HOST, cfg.API_PORT)
	log.Printf("API server is running on %s", address)
	return e.Start(address)
}