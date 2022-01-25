package books

import (
	"net/http"
	"project-airbnb/delivery/common"
	"project-airbnb/entities"
	"project-airbnb/repository/books"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BooksController struct {
	Repo books.BooksInterface
}

func NewBooksControllers(bkrep books.BooksInterface) *BooksController {
	return &BooksController{Repo: bkrep}
}

func (bkrep BooksController) Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Get("user").(*jwt.Token)
		claims := uid.Claims.(jwt.MapClaims)
		userID := int(claims["userid"].(float64))

		findBookReq := GetBookingRequestFormat{}
		if err := c.Bind(&findBookReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}
		if res, err := bkrep.Repo.Get(uint(userID), findBookReq.RoomID); err != nil {
			return c.JSON(http.StatusNotFound, common.NewStatusNotAcceptable())
		} else {
			responses := BookingsResponseFormat{
				Code:    http.StatusOK,
				Message: "Successful Operation",
				Data:    res,
			}
			return c.JSON(http.StatusOK, responses)
		}
	}
}

func (bkrep BooksController) Gets() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Get("user").(*jwt.Token)
		claims := uid.Claims.(jwt.MapClaims)
		userID := int(claims["userid"].(float64))

		if res, err := bkrep.Repo.Gets(uint(userID)); err != nil || len(res) == 0 {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		} else {
			responses := BookingsResponseFormat{
				Code:    http.StatusOK,
				Message: "Successful Operation",
				Data:    res,
			}
			return c.JSON(http.StatusOK, responses)
		}
	}
}

func (bkrep BooksController) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Get("user").(*jwt.Token)
		claims := uid.Claims.(jwt.MapClaims)
		userID := int(claims["userid"].(float64))

		newBookReq := NewBookingRequestFormat{}
		if err := c.Bind(&newBookReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}
		newUUID := uuid.New().String()
		newInvoice := "INV-" + strconv.Itoa(userID) + "/book/" + newUUID

		if res, err := bkrep.Repo.CreateTransactions(uint(userID), newBookReq.RoomID, newInvoice); err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		} else {
			responses := TransactionsResponseFormat{
				Code:    http.StatusOK,
				Message: "Successful Operation",
				Data:    res.Url,
			}

			newBook := entities.Book{
				User_id:        uint(userID),
				Room_id:        newBookReq.RoomID,
				Transaction_id: res.ID,
			}

			if res, err := bkrep.Repo.Create(newBook); err != nil || res.ID == 0 {
				return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
			}
			return c.JSON(http.StatusOK, responses)

		}
	}
}