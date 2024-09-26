package lib

import (
	"calometer/internal/db"
	"context"
	"time"

	"github.com/google/uuid"
)

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

func UpdateCalorieLog(
	userId uuid.UUID,
	logDate string,
	caloriesConsumed float64,
	caloriesBurnt float64,
) error {
	qStr := `
		UPDATE user_calorie_logs
		SET
			calories_consumed = user_calorie_logs.calories_consumed + $3,
			calories_burnt = user_calorie_logs.calories_burnt + $4
		WHERE u_id = $1 AND log_date = $2
	`

	if _, err := db.GetPool().Exec(
		context.Background(),
		qStr,
		userId,
		logDate,
		caloriesConsumed,
		caloriesBurnt,
	); err != nil {
		return err
	}

	return nil
}

func CheckLogStatusByIdAndDate(userId uuid.UUID, logDate string) (*string, error) {
	var logStatus string

	qStr := `
		SELECT log_status
		FROM user_calorie_logs
		WHERE u_id = $1 AND log_date = $2
	`

	if err := db.GetPool().QueryRow(context.Background(), qStr, userId, logDate).Scan(&logStatus); err != nil {
		return nil, err
	}

	return &logDate, nil
}

func DeleteCalorieLog(userId uuid.UUID, logDate string) error {
	qStr := `
		DELETE FROM user_calorie_logs
		WHERE u_id = $1 AND log_date = $2
	`

	if _, err := db.GetPool().Exec(context.Background(), qStr, userId, logDate); err != nil {
		return err
	}

	return nil
}
