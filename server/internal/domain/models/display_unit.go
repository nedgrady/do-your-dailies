package models

type DisplayUnit string

const (
	DisplayUnitDay   DisplayUnit = "DAY"
	DisplayUnitWeek  DisplayUnit = "WEEK"
	DisplayUnitMonth DisplayUnit = "MONTH"
	DisplayUnitYear  DisplayUnit = "YEAR"
)

func (u DisplayUnit) Valid() bool {
	switch u {
	case DisplayUnitDay, DisplayUnitWeek, DisplayUnitMonth, DisplayUnitYear:
		return true
	default:
		return false
	}
}
