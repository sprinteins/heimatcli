package heimat

// Balances _
type Balances struct {
	Holidays            float32 `json:"holidays"`
	HolidayEntitlement  float32 `json:"holidayEntitlement"`
	UnpaidHolidays      float32 `json:"unpaidHolidays"`
	Flexidays           float32 `json:"flexidays"`
	DaysOfIllness       float32 `json:"daysOfIllness"`
	WorkingHours        float32 `json:"workingHours"`
	PlannedWorkingHours float32 `json:"plannedWorkingHours"`
	BalanceWorkingHours float32 `json:"balanceWorkingHours"`
}
