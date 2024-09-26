package lib

import (
	"calometer/internal/db"
	"context"

	"github.com/google/uuid"
)

func DoesUserExists(username string) (*bool, error) {
	qStr := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE username = $1
		)`

	var doesExist bool
	if err := db.GetPool().QueryRow(context.Background(), qStr, username).Scan(&doesExist); err != nil {
		return nil, err
	}

	return &doesExist, nil
}

func GetUserIdByUsername(username string) (*uuid.UUID, error) {
	var userId uuid.UUID

	qStr := `
		SELECT id
		FROM users
		WHERE username = $1`

	if err := db.GetPool().QueryRow(context.Background(), qStr, username).Scan(&userId); err != nil {
		return nil, err
	}

	return &userId, nil
}

func AddUserBodyDetails(userId uuid.UUID, age int, height_cm int, weight_kg float64, gender string) error {
	bmr := CalculateBMR(gender, age, weight_kg, height_cm)

	qStr := `
		INSERT INTO user_body_details (
			u_id,
			age,
			height_cm,
			weight_kg,
			gender,
			bmr
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		) ON CONFLICT (u_id) DO UPDATE
		SET
			age = COALESCE(NULLIF($2, 0), user_body_details.age),
			height_cm = COALESCE(NULLIF($3, 0), user_body_details.height_cm),
			weight_kg = COALESCE(NULLIF($4, 0), user_body_details.weight_kg),
			gender = COALESCE(NULLIF($5, ''), user_body_details.gender),
			bmr = CASE
				WHEN user_body_details.age <> COALESCE(NULLIF($2, 0), user_body_details.age)
					OR user_body_details.height_cm <> COALESCE(NULLIF($3, 0), user_body_details.height_cm)
					OR user_body_details.weight_kg <> COALESCE(NULLIF($4, 0), user_body_details.weight_kg)
					OR user_body_details.gender <> COALESCE(NULLIF($5, ''), user_body_details.gender)
				THEN $6
				ELSE user_body_details.bmr
			END
	`

	if _, err := db.GetPool().Exec(
		context.Background(),
		qStr,
		userId,
		age,
		height_cm,
		weight_kg,
		gender,
		bmr,
	); err != nil {
		return err
	}

	return nil
}

func SetUserGoal(userId uuid.UUID, goal string) error {
	qStr := `
		INSERT INTO user_weight_goal (
			u_id,
			goal
		) VALUES (
			$1,
			$2
		) ON CONFLICT (u_id) DO UPDATE
		SET goal = COALESCE(NULLIF($2, ''), user_weight_goal.goal)
		`

	if _, err := db.GetPool().Exec(
		context.Background(),
		qStr,
		userId,
		goal,
	); err != nil {
		return err
	}

	return nil
}

func GetUserBmr(userId uuid.UUID) (*float64, error) {
	var bmr float64

	qStr := `
		SELECT bmr
		FROM user_body_details
		WHERE u_id = $1
	`

	if err := db.GetPool().QueryRow(context.Background(), qStr, userId).Scan(&bmr); err != nil {
		return nil, err
	}

	return &bmr, nil
}

func GetUserWeightGoalById(userId uuid.UUID) (*string, error) {
	var goal string

	qStr := `
		SELECT goal
		FROM user_weight_goal
		WHERE u_id = $1
	`

	if err := db.GetPool().QueryRow(context.Background(), qStr, userId).Scan(&goal); err != nil {
		return nil, err
	}

	return &goal, nil
}
