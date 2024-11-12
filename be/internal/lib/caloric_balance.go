package lib

import (
	"calometer/internal/db"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

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

func DeleteCaloricBalanceByLogId(logId uuid.UUID) error {
	qStr := `
		DELETE FROM user_caloric_balance
		WHERE calorie_log_id = $1
	`

	if _, err := db.GetPool().Exec(context.Background(), qStr, logId); err != nil {
		return err
	}

	return nil
}

func GetNetCaloricBalance(userId uuid.UUID) (*float64, error) {
	var netCaloricBalance sql.NullFloat64

	qStr := `
		SELECT SUM(user_caloric_balance.caloric_balance)
		FROM user_calorie_logs
		JOIN user_caloric_balance
		ON user_calorie_logs.id = user_caloric_balance.calorie_log_id
		WHERE user_calorie_logs.u_id = $1
	`

	if err := db.GetPool().QueryRow(context.Background(), qStr, userId).Scan(&netCaloricBalance); err != nil {
		return nil, err
	}

	if !netCaloricBalance.Valid {
		defaultValue := 0.0
		return &defaultValue, nil
	}

	return &netCaloricBalance.Float64, nil
}
