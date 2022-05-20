package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/middleware"
	"schej.it/server/models"
	"schej.it/server/utils"
)

func InitUser(router *gin.Engine) {
	userRouter := router.Group("/user")
	userRouter.Use(middleware.AuthRequired())

	userRouter.GET("/profile", getProfile)
	userRouter.GET("/events", getEvents)
	userRouter.GET("/calendar", getCalendar)
}

// @Summary Gets the user's profile
// @Tags user
// @Produce json
// @Success 200 {object} models.UserProfile "A user profile object"
// @Router /user/profile [get]
func getProfile(c *gin.Context) {
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	c.JSON(http.StatusOK, user.GetProfile())
}

// @Summary Gets all the user's events
// @Description Returns an array containing all the user's events
// @Tags user
// @Produce json
// @Success 200 {object} []models.Event
// @Router /user/events [get]
func getEvents(c *gin.Context) {
	session := sessions.Default(c)

	var events []models.Event
	cursor, err := db.EventsCollection.Find(context.Background(), bson.M{
		"ownerId": utils.GetUserId(session),
	})
	if err != nil {
		panic(err)
	}
	if err := cursor.All(context.Background(), &events); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, events)
}

// @Summary Gets the user's calendar events
// @Description Gets the user's calendar events between "timeMin" and "timeMax"
// @Tags user
// @Produce json
// @Param timeMin query string true "Lower bound for event's start time to filter by"
// @Param timeMax query string true "Upper bound for event's end time to filter by"
// @Success 200 {object} []models.CalendarEvent
// @Router /user/calendar [get]
func getCalendar(c *gin.Context) {
	// Bind query parameters
	payload := struct {
		TimeMin time.Time `form:"timeMin" binding:"required"`
		TimeMax time.Time `form:"timeMax" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}
	session := sessions.Default(c)

	// Refresh token if necessary
	user, _ := c.Get("authUser")
	db.RefreshUserTokenIfNecessary(user.(*models.User))

	// Call the google calendar API to get a list of calendar events from the user's gcal
	// TODO: get events for all user's calendars, not just primary
	min, _ := payload.TimeMin.MarshalText()
	max, _ := payload.TimeMax.MarshalText()
	fmt.Printf("https://www.googleapis.com/calendar/v3/calendars/primary/events?timeMin=%s&timeMax=%s\n", min, max)
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/primary/events?timeMin=%s&timeMax=%s&singleEvents=true", min, max),
		nil,
	)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Get("accessToken")))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	// Define some structs to parse the json response
	type TimeInfo struct {
		DateTime time.Time `json:"dateTime" binding:"required"`
	}
	type Item struct {
		Summary string   `json:"summary"`
		Start   TimeInfo `json:"start"`
		End     TimeInfo `json:"end"`
	}
	type Response struct {
		Items []Item      `json:"items"`
		Error interface{} `json:"error"`
	}

	// Parse the response
	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		panic(err)
	}

	// Check if the response returned an error
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Format response to return
	calendarEvents := make([]models.CalendarEvent, 0)
	for _, item := range res.Items {
		// Filter out invalid dates and restructure
		if payload.TimeMin.Before(item.Start.DateTime) && payload.TimeMax.After(item.End.DateTime) {
			calendarEvents = append(calendarEvents, models.CalendarEvent{
				Summary:   item.Summary,
				StartDate: primitive.NewDateTimeFromTime(item.Start.DateTime),
				EndDate:   primitive.NewDateTimeFromTime(item.End.DateTime),
			})
		}
	}

	c.JSON(http.StatusOK, calendarEvents)
}
