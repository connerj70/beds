package actions

import (
	"beds/models"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
)

// FriendsCreate default implementation.
func FriendsCreate(c buffalo.Context) error {
	var friend models.Friend

	if err := c.Bind(&friend); err != nil {
		return fmt.Errorf("failed to bind: %w", err)
	}

	// Check to ensure this user is not trying to create a friendship on behalf of another user.
	userID, err := c.Cookies().Get("user_id")
	if err != nil {
		return fmt.Errorf("failed to get userID from cookies: %w", err)
	}
	if friend.Requester.String() != userID {
		return fmt.Errorf("You can only create a friendship between yourself and another user")
	}

	// Check if this friendship already exists.
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	sqlStatement := `SELECT created_at FROM friends where (requester_id = ? AND receiver_id = ?) OR (receiver_id = ? AND requester_id = ?)`

	var friends []models.Friend
	qry := tx.RawQuery(sqlStatement, friend.Requester, friend.Receiver, friend.Requester, friend.Receiver)
	if err := qry.All(&friends); err != nil {
		return fmt.Errorf("failed to check if this friendship already exists: %w", err)
	}

	if len(friends) != 0 {
		return fmt.Errorf("you already became friends with this user on %v", friends[0].CreatedAt.Local())
	}

	if err := tx.Create(&friend); err != nil {
		return fmt.Errorf("failed to create friendship: %w", err)
	}

	// Go get information needed to send back a new friend struct
	var user models.User
	if err := tx.Where("id = ? ", friend.Receiver).First(&user); err != nil {
		return fmt.Errorf("failed to get information about a user: %w", err)
	}

	newFriend := struct {
		Email    string    `json:"email"`
		FriendID uuid.UUID `json:"friend_id"`
		Approved bool      `json:"approved"`
	}{
		Email:    user.Email,
		FriendID: user.ID,
		Approved: false,
	}

	return c.Render(http.StatusOK, r.JSON(newFriend))
}

func FriendsList(c buffalo.Context) error {

	userID := c.Param("id")

	txn, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	var friends []models.Friend

	if err := txn.Where("requester_id = ? OR receiver_id = ?", userID, userID).All(&friends); err != nil {
		return fmt.Errorf("failed to list friends: %w", err)
	}

	userFriends := []struct {
		Email    string    `json:"email"`
		FriendID uuid.UUID `json:"friend_id"`
		Approved bool      `json:"approved"`
	}{}

	for _, friend := range friends {
		var user models.User

		var friendLookupID uuid.UUID
		userIDUUID, err := uuid.FromString(userID)
		if err != nil {
			continue
		}
		if friend.Requester == userIDUUID {
			friendLookupID = friend.Receiver
		} else {
			friendLookupID = friend.Requester
		}

		if err := txn.Find(&user, friendLookupID); err != nil {
			c.Logger().Warn("user %s is friends with a deleted user %s", userID, friendLookupID)
			continue
		}
		// Remove password from user
		user.Password = ""

		userFriend := struct {
			Email    string    `json:"email"`
			FriendID uuid.UUID `json:"friend_id"`
			Approved bool      `json:"approved"`
		}{
			Email:    user.Email,
			FriendID: friendLookupID,
			Approved: friend.Accepted,
		}

		userFriends = append(userFriends, userFriend)
	}

	return c.Render(http.StatusOK, r.JSON(userFriends))
}

func FriendsListPage(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("/friends/index.plush.html"))
}
