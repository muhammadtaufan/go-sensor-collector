package delivery

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/muhammadtaufan/go-sensor-collector/config"
	"github.com/muhammadtaufan/go-sensor-collector/internal/types"
	"github.com/muhammadtaufan/go-sensor-collector/internal/usecase"
	"github.com/muhammadtaufan/go-sensor-collector/pkg"
)

type apiServer struct {
	sensorUsecase usecase.SensorSender
	userUsecase   usecase.User
	cfg           *config.Config
}

type BaseResponse struct {
	Success bool                       `json:"success"`
	Data    []types.SensorDataResponse `json:"data"`
	Error   string                     `json:"error,omitempty"`
}

func NewAPIServer(sensorUsecase usecase.SensorSender, userUsecase usecase.User, cfg *config.Config) *apiServer {
	return &apiServer{
		sensorUsecase: sensorUsecase,
		userUsecase:   userUsecase,
		cfg:           cfg,
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

	data, err := aps.sensorUsecase.GetSensorData(c.Request().Context(), id1, id2, startDate, endDate, &limit, &offset)
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

	err := aps.sensorUsecase.DeleteSensorData(c.Request().Context(), id1, id2, startDate, endDate)
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

	err := aps.sensorUsecase.UpdateSensorData(c.Request().Context(), id, &requestBody)
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

func (aps *apiServer) Login(c echo.Context) error {
	var userRequest types.UserRequest
	if err := c.Bind(&userRequest); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	user, err := aps.userUsecase.GetUser(c.Request().Context(), userRequest.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Success: false,
			Error:   "Opps, something's wrong",
		})
	}

	isValidPassword, err := pkg.ValidateHashPassword(user.Password, userRequest.Password)
	if err != nil || !isValidPassword {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Success: false,
			Error:   "Invalid Password",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte("SECRET_SENSOR"))
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})
}

func (aps *apiServer) StartServer(cfg *config.Config) error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("v1/api")
	v1.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("SECRET_SENSOR"),
	}))

	e.POST("/login", aps.Login)

	v1.GET("/sensors", aps.GetSensorData)
	v1.DELETE("/sensors", aps.DeleteSensorData)
	v1.PATCH("/sensors/:id", aps.UpdateSensorData)

	address := fmt.Sprintf("%s:%s", cfg.API_HOST, cfg.API_PORT)
	log.Printf("API server is running on %s", address)
	return e.Start(address)
}
