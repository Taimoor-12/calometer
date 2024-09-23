package lib

import (
	"calometer/internal/db"
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

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

func DoesLogExistForTheDay(userId uuid.UUID, logDate string) (*bool, error) {
	var logExists bool

	qStr := `
		SELECT EXISTS (
			SELECT 1
			FROM user_calorie_logs
			WHERE u_id = $1 AND log_date = $2
		)
	`

	if err := db.GetPool().QueryRow(context.Background(), qStr, userId, logDate).Scan(&logExists); err != nil {
		return nil, err
	}

	return &logExists, nil
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

func CreateUserLog(userId uuid.UUID, bmr float64, logDate string) error {
	qStr := `
		INSERT INTO user_calorie_logs (
			u_id,
			tdee,
			log_date
		) VALUES (
			$1,
			$2,
			$3
		)
	`
	if _, err := db.GetPool().Exec(context.Background(), qStr, userId, bmr, logDate); err != nil {
		return err
	}

	return nil
}

func LogCaloriesConsumed(userId uuid.UUID, caloriesConsumed float64, logDate string) error {
	qStr := `
		UPDATE user_calorie_logs
		SET
			calories_consumed = COALESCE(user_calorie_logs.calories_consumed, 0) + $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE u_id = $1 AND log_date = $3
		`

	if _, err := db.GetPool().Exec(
		context.Background(),
		qStr,
		userId,
		caloriesConsumed,
		logDate,
	); err != nil {
		return err
	}

	return nil
}

func LogCaloriesBurnt(userId uuid.UUID, caloriesBurnt float64, logDate string) error {
	qStr := `
		UPDATE user_calorie_logs
		SET
			calories_burnt = COALESCE(user_calorie_logs.calories_burnt, 0) + $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE u_id = $1 AND log_date = $3
		`

	if _, err := db.GetPool().Exec(
		context.Background(),
		qStr,
		userId,
		caloriesBurnt,
		logDate,
	); err != nil {
		return err
	}

	return nil
}

func FetchCaloriesConsumedForTheDay(userId uuid.UUID, logDate string) (*float64, error) {
	var caloriesConsumed float64

	qStr := `
		SELECT calories_consumed
		FROM user_calorie_logs
		WHERE u_id = $1 AND log_date = $2
	`

	if err := db.GetPool().QueryRow(context.Background(), qStr, userId, logDate).Scan(&caloriesConsumed); err != nil {
		return nil, err
	}

	return &caloriesConsumed, nil
}

func FetchCaloriesBurntForTheDay(userId uuid.UUID, logDate string) (*float64, error) {
	var caloriesBurnt float64

	qStr := `
		SELECT calories_burnt
		FROM user_calorie_logs
		WHERE u_id = $1 AND log_date = $2
	`

	if err := db.GetPool().QueryRow(context.Background(), qStr, userId, logDate).Scan(&caloriesBurnt); err != nil {
		return nil, err
	}

	return &caloriesBurnt, nil
}

func AddCaloriesBurntInTDEE(userId uuid.UUID, logDate string, caloriesBurnt float64) error {
	qStr := `
		UPDATE user_calorie_logs
		SET tdee = user_calorie_logs.tdee + $3
		WHERE u_id = $1 AND log_date = $2
	`

	if _, err := db.GetPool().Exec(
		context.Background(),
		qStr,
		userId,
		logDate,
		caloriesBurnt,
	); err != nil {
		return err
	}

	return nil
}

func MarkLoggingStatus(userId uuid.UUID, logDate string, logStatus string) error {
	qStr := `
		UPDATE user_calorie_logs
		SET log_status = $3
		WHERE u_id = $1 AND log_date = $2
	`

	if _, err := db.GetPool().Exec(
		context.Background(),
		qStr,
		userId,
		logDate,
		logStatus,
	); err != nil {
		return err
	}

	return nil
}

func CalculateCaloricBalanceForTheDay(userId uuid.UUID, logDate string) (*float64, error) {
	var caloricBalance float64

	qStr := `
		SELECT (tdee - calories_consumed) AS caloric_balance
		FROM user_calorie_logs
		WHERE u_id = $1 AND log_date = $2 AND log_status = 'D'
	`

	if err := db.GetPool().QueryRow(context.Background(), qStr, userId, logDate).Scan(&caloricBalance); err != nil {
		return nil, err
	}

	return &caloricBalance, nil
}

func GetCalorieLogId(userId uuid.UUID, logDate string) (*uuid.UUID, error) {
	var calorieLogId uuid.UUID

	qStr := `
		SELECT id
		FROM user_calorie_logs
		WHERE u_id = $1 AND log_date = $2
	`

	if err := db.GetPool().QueryRow(context.Background(), qStr, userId, logDate).Scan(&calorieLogId); err != nil {
		return nil, err
	}

	return &calorieLogId, nil
}

func AddCaloricBalanceForTheDay(logId uuid.UUID, caloricBalance float64) error {
	qStr := `
		INSERT INTO user_caloric_balance (
			calorie_log_id,
			caloric_balance
		) VALUES (
			$1,
			$2
		)
	`

	if _, err := db.GetPool().Exec(context.Background(), qStr, logId, caloricBalance); err != nil {
		return err
	}

	return nil
}

func ResetCaloricBalanceForTheDay(logId uuid.UUID) error {
	qStr := `
		UPDATE user_caloric_balance
		SET caloric_balance = 0.00
		WHERE calorie_log_id = $1
	`

	if _, err := db.GetPool().Exec(context.Background(), qStr, logId); err != nil {
		return err
	}

	return nil
}

type UserCalorieLogs struct {
	LogDate          string
	CaloriesBurnt    float64
	CaloriesConsumed float64
	Tdee             float64
	Updated_at       time.Time
	LogStatus        string
}

func GetCalorieLogs(userId uuid.UUID) (*[]UserCalorieLogs, error) {
	var userCalorieLogs []UserCalorieLogs

	qStr := `
		SELECT
			log_date,
			calories_burnt,
			calories_consumed,
			tdee,
			updated_at,
			log_status
		FROM user_calorie_logs
		WHERE u_id = $1
	`

	rows, err := db.GetPool().Query(context.Background(), qStr, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var log UserCalorieLogs
		var logDate time.Time

		err := rows.Scan(
			&logDate,
			&log.CaloriesBurnt,
			&log.CaloriesConsumed,
			&log.Tdee,
			&log.Updated_at,
			&log.LogStatus,
		)
		if err != nil {
			return nil, err
		}

		log.LogDate = logDate.Format("2006-01-02")

		userCalorieLogs = append(userCalorieLogs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &userCalorieLogs, nil
}
