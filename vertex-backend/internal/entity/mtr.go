package entity

import (
	"time"

	"github.com/google/uuid"
)

// MTRSession stores the configuration and status of an MTR (traceroute) test session.
// Relation: MTRSession -> MTRResult -> MTRHop
type MTRSession struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Status    string     `gorm:"type:varchar(20);not null;index" json:"status"`
	IsReachable bool      `gorm:"not null;default:false" json:"isReachable"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	LastRunAt *time.Time `gorm:"index" json:"lastRunAt"`

	SourceIP      string `gorm:"type:varchar(45)" json:"sourceIp"`
	DestinationIP string `gorm:"type:varchar(45);not null;index" json:"destinationIp"`
	Protocol      string `gorm:"type:varchar(10);not null;default:'icmp';index" json:"protocol"`
	Port          *int   `json:"port"`
	Test          int    `gorm:"not null;default:10" json:"test"`
	Note          string `gorm:"type:varchar(500);default:''" json:"note"`

	SendNotification bool `gorm:"not null;default:false" json:"sendNotification"`

	Results []MTRResult `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"results"`
}

// MTRResult stores the summary result of a single MTR execution on a target.
type MTRResult struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SessionID     uuid.UUID `gorm:"type:uuid;not null;index" json:"sessionId"`
	SourceIP      string    `gorm:"type:varchar(45)" json:"sourceIp"`
	DestinationIP string    `gorm:"type:varchar(45)" json:"destinationIp"`
	Protocol      string    `gorm:"type:varchar(10)" json:"protocol"`
	Port          *int      `json:"port"`
	Test          int       `gorm:"not null;default:10" json:"test"`
	TotalHops     int       `json:"totalHops"`
	Reachable     bool      `json:"reachable"`
	AvgRTT        float64   `json:"avgRtt"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`

	Session  *MTRSession `gorm:"foreignKey:SessionID;references:ID" json:"session,omitempty"`
	Hops     []MTRHop    `gorm:"foreignKey:ResultID;references:ID;constraint:OnDelete:CASCADE" json:"hops"`
}

// MTRHop stores a single hop entry from an MTR traceroute result.
type MTRHop struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ResultID uuid.UUID `gorm:"type:uuid;not null;index" json:"resultId"`
	Hop      int       `gorm:"index" json:"hop"`
	Host     string    `gorm:"type:varchar(255)" json:"host"`
	DNS      string    `gorm:"type:varchar(255);index" json:"dns"`
	Loss     float64   `json:"loss"`
	Sent     int       `json:"sent"`
	Last     float64   `json:"last"`
	Avg      float64   `json:"avg"`
	Best     float64   `json:"best"`
	Worst    float64   `json:"worst"`
	StdDev   float64   `json:"stdDev"`

	Result *MTRResult `gorm:"foreignKey:ResultID;references:ID" json:"result,omitempty"`
}
