package handlers

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"go-webapp/config"
	"go-webapp/dbapi/user"
	"go-webapp/instagram"
	"go-webapp/models"
	"go-webapp/utils"
)

// UserViewModel - User view model
type UserViewModel struct {
	ID        bson.ObjectId `json:"id"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	Token     string        `json:"token"`
}

// Save user, handle response and return token if necessary
func saveUser(user *models.User, c *gin.Context) {
	session := c.MustGet("UserSession").(*models.Session)

	// if oauth is already registered with existing user id then raise error
	if user.ID.Valid() && session.UserID.Valid() && user.ID != session.UserID {
		c.JSON(422, "Already registered with some other id")
		return
	}

	if session.UserID.Valid() {
		user.ID = session.UserID
	}

	err := userapi.Upsert(user)

	if err != nil {
		panic(err)
	}

	userVM := &UserViewModel{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		ID:        user.BaseModel.ID,
	}

	if session.UserID == "" {
		userVM.Token = userapi.GenerateToken(*user)
	}

	c.JSON(200, userVM)
}


func Instagram(c *gin.Context) {
	var ig models.Instagram

	// get access token
	q, _, err := utils.MakeOauthQparams(c.Request.Body)
	q.Add("client_secret", config.Settings.IgSecret)
	q.Add("scope", "basic+public_content+follower_list")
	err = instagram.GetAccessToken(&ig, q)
	if err != nil {
		panic(err.Error())
	}

	user, err := userapi.GetByIgID(ig.ProfileID)
	if user == nil {
		user = &models.User{
			Ig:        ig,
			FirstName: ig.FirstName,
			LastName:  ig.LastName,
			Email:     ig.Email,
		}
	} else {
		user.Ig = ig
	}

	err = instagram.SaveInitFeed(&user.Ig)
	if err != nil {
		panic(err)
	}

	saveUser(user, c)
}
