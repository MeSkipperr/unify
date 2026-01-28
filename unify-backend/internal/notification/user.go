package notification

import (
	"errors"
	"fmt"
	"unify-backend/internal/mailer"
	"unify-backend/internal/services"
	"unify-backend/models"
)

func UserNotificationChannel(data mailer.EmailData) (string, error) {
	users, err := SelectUserByType([]models.UserRole{
		models.RoleUser,
	})

	if err != nil {
		services.CreateAppLog(services.CreateLogParams{
			Level:       "ERROR",
			ServiceName: "user-notification-channel",
			Message:     "Failed to fetch users: " + err.Error(),
		})
		return "", err
	}

	var recipients []mailer.Recipients
	for i, user := range users {
		if user.Email == nil {
			continue
		}
		fmt.Print(i, user.Email)

		recipients = append(recipients, mailer.Recipients{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     *user.Email,
		})
	}

	if len(recipients) == 0 {
		return "", errors.New("no recipients available")
	}

	err = mailer.SendEmailSMTP(mailer.EmailStructure{
		Recipients: recipients,
		EmailData: mailer.EmailData{
			Subject:        data.Subject,
			BodyTemplate:   data.BodyTemplate,
			FileAttachment: data.FileAttachment,
		},
	})

	if err != nil {
		services.CreateAppLog(services.CreateLogParams{
			Level:       "ERROR",
			ServiceName: "user-notification-channel",
			Message:     "Failed to send user notification email: " + err.Error(),
		})
		return "", err
	}

	return "user_notification_sent", nil
}
