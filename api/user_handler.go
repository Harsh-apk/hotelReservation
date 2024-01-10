package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/harsh-apk/hotelReservation/db"
	"github.com/harsh-apk/hotelReservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	res, err := h.userStore.DeleteUser(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	if res == 0 {
		return c.JSON(map[string]string{"error": "No user found"})
	}
	return c.JSON(map[string]string{"Status": fmt.Sprintf("Successfully deleted %d users", res)})
}
func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := len(params.Validate()); err != 0 {
		return c.JSON(err)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}
func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	errors := map[string]string{}
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	var values types.UpdateUserParams
	err = c.BodyParser(&values)
	if err != nil {
		return err
	}
	update := bson.M{}
	if len(values.FirstName) > types.MinFirstNameLen {
		update["firstName"] = values.FirstName
	} else {
		errors["firstNameErr"] = fmt.Sprintf("first name must have atleast %d characters :)", types.MinFirstNameLen)
	}
	if len(values.LastName) > types.MinLastNameLen {
		update["lastName"] = values.LastName
	} else {
		errors["lastNameErr"] = fmt.Sprintf("last name must have atleast %d characters :)", types.MinLastNameLen)
	}
	if update["firstName"] == nil && update["lastName"] == nil {
		return c.JSON(errors)
	}
	res, err := h.userStore.UpdateUser(c.Context(), filter, update)
	if err != nil {
		return err
	}
	errors["result"] = fmt.Sprintf("Updated the values of %d documents which didn't had erroR", res)
	return c.JSON(errors)

}
