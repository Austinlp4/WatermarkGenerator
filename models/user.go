package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email                 string             `bson:"email" json:"email"`
	Password              string             `bson:"password" json:"password"` // Make sure this line exists
	StripeCustomerID      string             `bson:"stripeCustomerId" json:"stripeCustomerId,omitempty"`
	SubscriptionStatus    string             `bson:"subscriptionStatus" json:"subscriptionStatus"`
	SubscriptionId        string             `bson:"subscriptionId" json:"subscriptionId,omitempty"`
	SubscriptionExpiresAt time.Time          `bson:"subscriptionExpiresAt" json:"subscriptionExpiresAt"`
	DailyDownloads        int                `bson:"dailyDownloads" json:"dailyDownloads"`
	LastDownloadDate      time.Time          `bson:"lastDownloadDate" json:"lastDownloadDate"`
}
