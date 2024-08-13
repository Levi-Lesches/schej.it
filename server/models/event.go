package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventType string

const (
	SPECIFIC_DATES EventType = "specific_dates"
	DOW            EventType = "dow"
	GROUP          EventType = "group"
)

// Representation of an Event in the mongoDB database
type Event struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	ShortId     *string            `json:"shortId" bson:"shortId,omitempty"`
	OwnerId     primitive.ObjectID `json:"ownerId" bson:"ownerId,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Description *string            `json:"description" bson:"description,omitempty"`

	Duration                 *float32             `json:"duration" bson:"duration,omitempty"`
	Dates                    []primitive.DateTime `json:"dates" bson:"dates,omitempty"`
	NotificationsEnabled     *bool                `json:"notificationsEnabled" bson:"notificationsEnabled,omitempty"`
	SendEmailAfterXResponses *int                 `json:"sendEmailAfterXResponses" bson:"sendEmailAfterXResponses,omitempty"`
	When2meetHref            *string              `json:"when2meetHref" bson:"when2meetHref,omitempty"`
	CollectEmails            *bool                `json:"collectEmails" bson:"collectEmails,omitempty"`

	Type EventType `json:"type" bson:"type,omitempty"`

	// Whether to enable blind availability
	BlindAvailabilityEnabled *bool `json:"blindAvailabilityEnabled" bson:"blindAvailabilityEnabled,omitempty"`

	// Whether to only poll for days, not times
	DaysOnly *bool `json:"daysOnly" bson:"daysOnly,omitempty"`

	// Availability responses
	Responses map[string]*Response `json:"responses" bson:"responses"`

	// Scheduled event
	ScheduledEvent  *CalendarEvent `json:"scheduledEvent" bson:"scheduledEvent,omitempty"`
	CalendarEventId string         `json:"calendarEventId" bson:"calendarEventId,omitempty"`

	// Remindees
	Remindees *[]Remindee `json:"remindees" bson:"remindees,omitempty"`

	// Attendees for an availability group
	Attendees *[]Attendee `json:"attendees" bson:"attendees,omitempty"`
}

// A response object containing an array of times that the given user is available
type Response struct {
	// Guest information
	Name  string `json:"name" bson:"name,omitempty"`
	Email string `json:"email" bson:"email,omitempty"`

	// User information
	UserId primitive.ObjectID `json:"userId" bson:"userId,omitempty"`
	User   *User              `json:"user" bson:",omitempty"`

	// Availability
	Availability []primitive.DateTime `json:"availability" bson:"availability"`
	IfNeeded     []primitive.DateTime `json:"ifNeeded" bson:"ifNeeded"`

	// Mapping from the start date of a day to the available times for that day
	ManualAvailability *map[primitive.DateTime][]primitive.DateTime `json:"manualAvailability" bson:"manualAvailability,omitempty"`

	// Calendar availability variables for Availability Groups feature
	UseCalendarAvailability *bool                `json:"useCalendarAvailability" bson:"useCalendarAvailability,omitempty"`
	EnabledCalendars        *map[string][]string `json:"enabledCalendars" bson:"enabledCalendars,omitempty"` // Maps email to an array of sub calendar ids
	CalendarOptions         *CalendarOptions     `json:"calendarOptions" bson:"calendarOptions,omitempty"`
}

// Object containing information associated with the remindee
type Remindee struct {
	Email     string   `json:"email" bson:"email,omitempty"`
	TaskIds   []string `json:"-" bson:"taskIds,omitempty"` // Task IDs of the scheduled emails
	Responded *bool    `json:"responded" bson:"responded,omitempty"`
}

type Attendee struct {
	Email    string `json:"email" bson:"email,omitempty"`
	Declined *bool  `json:"declined" bson:"declined,omitempty"`
}

func (e *Event) GetId() string {
	if e.ShortId != nil {
		return *e.ShortId
	}

	return e.Id.Hex()
}
