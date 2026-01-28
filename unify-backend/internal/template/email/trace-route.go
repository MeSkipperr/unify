package email

import (
	"fmt"
	"unify-backend/internal/core/mtr"
	"unify-backend/models"
	"unify-backend/utils"
)

type TraceRouteEmailParms struct {
	Session models.MTRSession
	Result  mtr.MtrResultJson
}

func TraceRouteEmail(payload TraceRouteEmailParms, isReachable bool) string {
	session := payload.Session
	result := payload.Result.Report
	currentTime := utils.GetCurrentTimeFormatted()

	// default (DOWN)
	intro := "We would like to inform you that a network issue has been detected. Please find the details below:"
	action := "Kindly review the information provided and take the necessary actions to resolve the issue at your earliest convenience."
	status := "Unreachable"
	// Recovery state
	if isReachable {
		intro = "This is to notify you of a network system recovery update. Below are the recovery details:"
		action = "No further action is required. The network system is operating normally."
		status = "Reachable"
	}

	portText := "N/A"
	if session.Port != nil {
		portText = fmt.Sprintf("%d", *session.Port)
	}

	

	return fmt.Sprintf(`
Dear {{firstName}} {{lastName}},

%s

- Timestamp : %s
- Session ID : %s
- Source IP : %s
- Destination IP : %s
- Protocol : %s
- Port : %s
- Status : %s
- Note : %s

%s

Best regards,
{{PROPERTY}}
`,
		intro,
		currentTime,
		session.ID,
		result.HopResult[0].Host,
		session.DestinationIP,
		session.Protocol,
		portText,
		status,
		session.Note,
		action,
	)
}
