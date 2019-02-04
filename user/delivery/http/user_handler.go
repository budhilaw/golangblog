package http

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"

	_mw "github.com/budhilaw/blog/middleware"
	"github.com/budhilaw/blog/models"
	"github.com/budhilaw/blog/user"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError : Represent of error response message
type ResponseError struct {
	Message string `json:"message"`
}

// UserHandler : Represent of HTTPHandler of user
type UserHandler struct {
	UserUcase user.Usecase
}

type Claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// New : Will create object that represent of the user.Usecase handler interface
func New(e *echo.Echo, u user.Usecase) {
	handler := &UserHandler{
		UserUcase: u,
	}

	mw := _mw.InitMiddleware()
	e.Use(mw.CORS)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	config := middleware.JWTConfig{
		Claims:     &Claims{},
		SigningKey: []byte(viper.GetString("jwt.key")),
	}
	user := e.Group("/user", middleware.JWTWithConfig(config))
	user.GET("", handler.Fetch)
	user.POST("/new", handler.Store)

	e.POST("/login", handler.Login)
}

// Fetch : Fetch list of users
func (u *UserHandler) Fetch(c echo.Context) error {
	param := c.QueryParam("num")
	num, err := strconv.Atoi(param)
	if err != nil {
		num = 10
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	listUser, err := u.UserUcase.Fetch(ctx, int64(num))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, listUser)
}

// Store : Add a new user
func (u *UserHandler) Store(c echo.Context) error {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	if ok, err := isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = u.UserUcase.Store(ctx, &user)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

// Login : Login handler
func (u *UserHandler) Login(c echo.Context) error {
	var user models.User
	email := c.FormValue("email")
	pass := c.FormValue("password")

	if len(email) == 0 {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Username cannot be empty"})
	} else if len(pass) == 0 {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Password cannot be empty"})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := u.UserUcase.Login(ctx, email, pass, &user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	expiredTime := time.Now().Add(10 * time.Minute)

	claims := &Claims{
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
		},
	}

	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := sign.SignedString([]byte(viper.GetString("jwt.key")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"name":  user.Name,
		"email": user.Email,
		"token": token,
	})
}

func isRequestValid(m *models.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getStatusCode(err error) int {

	if err == nil {
		return http.StatusOK
	}
	log.Fatal(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
