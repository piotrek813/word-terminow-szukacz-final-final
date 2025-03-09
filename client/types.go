package client

type Exam struct {
	ID             string `json:"id"`
	Places         int    `json:"places"`
	Date           string `json:"date"`
	Amount         int    `json:"amount"`
	AdditionalInfo string `json:"additionalInfo"`
}

type ScheduledHour struct {
	Time          string `json:"time"`
	TheoryExams   []Exam `json:"theoryExams"`
	PracticeExams []Exam `json:"practiceExams"`
	LinkedExams   []Exam `json:"linkedExamsDto"`
}

type ScheduledDay struct {
	Day            string          `json:"day"`
	ScheduledHours []ScheduledHour `json:"scheduledHours"`
}

type Schedule struct {
	ScheduledDays []ScheduledDay `json:"scheduledDays"`
}

type Reservation struct {
	OrganizationID                 string   `json:"organizationId"`
	IsOskVehicleReservationEnabled bool     `json:"isOskVehicleReservationEnabled"`
	IsRescheduleReservation        bool     `json:"isRescheduleReservation"`
	Category                       string   `json:"category"`
	Schedule                       Schedule `json:"schedule"`
}
