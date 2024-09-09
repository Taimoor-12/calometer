package lib

func CalculateBMR(
	gender byte,
	age int,
	weight_kg float64,
	height_cm int,
) float64 {
	var bmr float64
	// Using Mifflin-St Jeor Equation
	if gender == 'M' {
		bmr = (10 * weight_kg) + (6.25 * float64(height_cm)) - (5 * float64(age)) + 5
	} else {
		bmr = (10 * weight_kg) + (6.25 * float64(height_cm)) - (5 * float64(age)) - 161
	}

	return bmr
}
