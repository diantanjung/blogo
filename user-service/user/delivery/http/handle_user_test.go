package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"strconv"
	"time"
	"encoding/json"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/diantanjung/blogo/user-service/domain"
	"github.com/diantanjung/blogo/user-service/domain/mocks"
	userHttp "github.com/diantanjung/blogo/user-service/user/delivery/http"
)

func TestFetch(t *testing.T) {
	var mockUser domain.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)
	mockUCase := new(mocks.UserUsecase)
	mockListUser := make([]domain.User, 0)
	mockListUser = append(mockListUser, mockUser)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListUser, "10", nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/users?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := userHttp.UserHandler{
		UserUsecase: mockUCase,
	}
	err = handler.Fetch(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	var mockUser domain.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)

	mockUCase := new(mocks.UserUsecase)

	num := int(mockUser.ID)

	mockUCase.On("GetByID", mock.Anything, int64(num)).Return(mockUser, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/user/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("user/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := userHttp.UserHandler{
		UserUsecase: mockUCase,
	}
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockUser := domain.User{
		Username:   "username",
		Name:   	"Name",
		Email:   	"Email@gmail.com",
		Password:   "ASDF1234",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tempMockUser := mockUser
	tempMockUser.ID = 0
	mockUCase := new(mocks.UserUsecase)

	j, err := json.Marshal(tempMockUser)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/user", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/user")

	handler := userHttp.UserHandler{
		UserUsecase: mockUCase,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	mockUser := domain.User{
		ID:   		1,
		Username:   "username",
		Name:   	"Name",
		Email:   	"Email@gmail.com",
		Password:   "ASDF1234",
		UpdatedAt: time.Now(),
	}

	tempMockUser := mockUser
	tempMockUser.ID = 0
	mockUCase := new(mocks.UserUsecase)

	j, err := json.Marshal(tempMockUser)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.PATCH, "/user", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/user")

	handler := userHttp.UserHandler{
		UserUsecase: mockUCase,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	var mockUser domain.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)

	mockUCase := new(mocks.UserUsecase)

	num := int(mockUser.ID)

	mockUCase.On("Delete", mock.Anything, int64(num)).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/user/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("user/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := userHttp.UserHandler{
		UserUsecase: mockUCase,
	}
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUCase.AssertExpectations(t)

}
