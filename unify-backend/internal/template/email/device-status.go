package email

import (
	"fmt"
	"strings"
	"unify-backend/models"
	"unify-backend/utils"
)

func DeviceStatusEmail(dev models.Devices, isConnect bool) string {
	currentTime := utils.GetCurrentTimeFormatted()

	// optional description
	desc := ""
	if trimmed := strings.TrimSpace(dev.Description); trimmed != "" {
		desc = fmt.Sprintf("- Description : %s\n", trimmed)
	}

	// default (DOWN)
	intro := "We would like to inform you that an error has occurred in the network system. Below are the details:"
	action := "Kindly review the information provided and take necessary actions to resolve the issue at your earliest convenience."
	
	// Recovery state 
	if isConnect {
		intro = "This is to notify you of a network system recovery update. Below are the recovery details:"
		action = "Inspect the details and confirm the system is back to normal."
	}

	return fmt.Sprintf(`
Dear {{firstName}} {{lastName}},

%s

- Time : %s
- Host Name: %s
- IP Address: %s
- Device: %s
- Mac Address: %s
%s
%s

Best regards,
{{PROPERTY}}
`,
		intro,
		currentTime,
		dev.Name,
		dev.IPAddress,
		dev.Name,
		dev.MacAddress,
		desc,
		action,
	)
}
