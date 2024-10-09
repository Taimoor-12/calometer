package lib

func CalculateBMR(
	gender string,
	age int,
	weight_kg float64,
	height_cm int,
) float64 {
	var bmr float64
	// Using Harris-Benedict Equation
	if gender == "M" {
		bmr = 66.5 + (13.75 * float64(weight_kg)) + (5.003 * float64(height_cm)) - (6.75 * float64(age))
	} else {
		bmr = 655.1 + (9.563 * float64(weight_kg)) + (1.850 * float64(height_cm)) - (4.676 * float64(age))
	}

	return bmr
}
