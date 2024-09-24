package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"calometer/internal/db"
	"calometer/internal/lib"
	"calometer/internal/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

var log = logger.GetLogger()

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

type SignupReq struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user SignupReq
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Info("failed to decode incoming json")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" || user.Name == "" {
		log.Info("invalid input data")
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	doesExist, err := lib.DoesUserExists(user.Username)
	if err != nil {
		log.Info(
			"failed to check user's existence by username",
			zap.String("username", user.Username),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if *doesExist {
		log.Info("username already exists")
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	password, err := lib.HashPassword(user.Password)
	if err != nil {
		log.Info(
			"failed to hash user's password",
			zap.String("username", user.Username),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// Save user to the database
	qStr := `
	INSERT INTO users (name, username, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id
  `

	var userId uuid.UUID
	if err := db.GetPool().QueryRow(context.Background(), qStr, user.Name, user.Username, password).Scan(&userId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			log.Info(
				"username already exists",
				zap.String("username", user.Username),
			)
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}

		log.Info(
			"failed to create the user",
			zap.String("username", user.Username),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	response := Response{
		Code: http.StatusOK,
		Data: map[string]uuid.UUID{"user_id": userId},
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: Check for JWT in the request
	tokenStr := lib.ExtractTokenFromHeader(r)
	if tokenStr != "" {
		// Step 2: Validate the JWT
		if err := lib.ValidateToken(tokenStr); err == nil {
			// Token is valid, return a success response
			w.WriteHeader(http.StatusOK)
			resp := Response{
				Code: http.StatusOK,
				Data: map[string]string{"token": tokenStr},
			}
			json.NewEncoder(w).Encode(&resp)
			return
		}
	}

	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info("failed to decode incoming json")
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if req.Username == "" && req.Password == "" {
		log.Info("invalid input data")
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	userId, err := lib.GetUserIdByUsername(req.Username)
	if err != nil {
		log.Fatal(
			"failed to get user id by username",
			zap.String("username", req.Username),
		)
	}

	exists, err := lib.DoesUserExists(req.Username)
	if err != nil {
		log.Info(
			"failed to check user's existence by username",
			zap.String("username", req.Username),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !*exists {
		http.Error(w, "Username or Password is incorrect", http.StatusUnauthorized)
		return
	}

	passwordHash, err := lib.GetHashedPass(req.Username)
	if err != nil {
		log.Info(
			"failed to fetch user's hashed password",
			zap.String("username", req.Username),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := lib.CheckPasswordValidity(req.Password, passwordHash); err != nil {
		log.Info(
			"failed to check password validity",
			zap.String("username", req.Username),
		)
		http.Error(w, "Username or Password is incorrect", http.StatusUnauthorized)
		return
	}

	token, err := lib.GenerateJWT(*userId, req.Username)
	if err != nil {
		log.Info(
			"failed to generate JWT for user id",
			zap.String("userId", userId.String()),
			zap.String("username", req.Username),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := Response{
		Code: http.StatusOK,
		Data: map[string]string{"token": token},
	}
	json.NewEncoder(w).Encode(&resp)
}

type AddBodyDetailsReq struct {
	Age       int     `json:"age"`
	Weight_kg float64 `json:"weight"`
	Height_cm int     `json:"height"`
	Gender    string  `json:"gender"`
	Goal      string  `json:"goal,omitempty"`
}

func AddBodyDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the token from the context
	tokenStr, ok := r.Context().Value(TokenContextKey).(string)
	if !ok {
		log.Info(
			"token not found in context",
		)

		// Token is not present in context
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var req AddBodyDetailsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info(
			"failed to decode incoming json",
			zap.Error(err),
		)

		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userId, err := lib.ExtractUserIdFromToken(tokenStr)
	if err != nil {
		log.Info(
			"failed to extract username from token",
			zap.String("tokenStr", tokenStr),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := lib.AddUserBodyDetails(
		*userId,
		req.Age,
		req.Height_cm,
		req.Weight_kg,
		req.Gender,
	); err != nil {
		log.Info(
			"failed to add user body details by id",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := Response{
		Code: http.StatusOK,
	}
	json.NewEncoder(w).Encode(&resp)
}

type SetUserWeightGoalReq struct {
	Goal string `json:"goal"`
}

func SetUserWeightGoalHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the token from the context
	tokenStr, ok := r.Context().Value(TokenContextKey).(string)
	if !ok {
		log.Info(
			"token not found in context",
		)

		// Token is not present in context
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var req SetUserWeightGoalReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info(
			"failed to decode incoming json",
			zap.Error(err),
		)

		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userId, err := lib.ExtractUserIdFromToken(tokenStr)
	if err != nil {
		log.Info(
			"failed to extract username from token",
			zap.String("tokenStr", tokenStr),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := lib.SetUserGoal(*userId, req.Goal); err != nil {
		log.Info(
			"failed to set user's goal by id",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := Response{
		Code: http.StatusOK,
	}
	json.NewEncoder(w).Encode(&resp)
}

type MarkLoggingStatusReq struct {
	Status  string    `json:"status"`
	LogDate time.Time `json:"log_date"`
}

func MarkLoggingStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the token from the context
	tokenStr, ok := r.Context().Value(TokenContextKey).(string)
	if !ok {
		log.Info(
			"token not found in context",
		)

		// Token is not present in context
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userId, err := lib.ExtractUserIdFromToken(tokenStr)
	if err != nil {
		log.Info(
			"failed to get user id by token",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var req MarkLoggingStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info(
			"failed to decode incoming json",
			zap.Error(err),
		)

		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	logDate := req.LogDate.Format("2006-01-02")

	if err := lib.MarkLoggingStatus(*userId, logDate, req.Status); err != nil {
		log.Info(
			"failed to mark logging status by id and date",
			zap.String("userId", userId.String()),
			zap.String("logDate", logDate),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	calorieLogId, err := lib.GetCalorieLogId(*userId, logDate)
	if err != nil {
		log.Info(
			"failed to get calorie log by id and date",
			zap.String("userId", userId.String()),
			zap.String("logDate", logDate),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if req.Status == "D" {
		caloricBalance, err := lib.CalculateCaloricBalanceForTheDay(*userId, logDate)
		if err != nil {
			log.Info(
				"failed to calculate caloric balance by id and date",
				zap.String("userId", userId.String()),
				zap.String("logDate", logDate),
				zap.Error(err),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := lib.AddCaloricBalanceForTheDay(*calorieLogId, *caloricBalance); err != nil {
			log.Info(
				"failed to add caloric balance by log id",
				zap.String("logId", calorieLogId.String()),
				zap.Error(err),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		if err := lib.ResetCaloricBalanceForTheDay(*calorieLogId); err != nil {
			log.Info(
				"failed to reset caloric balance by log id",
				zap.String("logId", calorieLogId.String()),
				zap.Error(err),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	resp := Response{
		Code: http.StatusOK,
	}
	json.NewEncoder(w).Encode(&resp)
}

type CreateCalorieLogReq struct {
	LogDate time.Time `json:"log_date,omitempty"`
}

func CreateCalorieLogHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the token from the context
	tokenStr, ok := r.Context().Value(TokenContextKey).(string)
	if !ok {
		log.Info(
			"token not found in context",
		)

		// Token is not present in context
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userId, err := lib.ExtractUserIdFromToken(tokenStr)
	if err != nil {
		log.Info(
			"failed to get user id by token",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var req CreateCalorieLogReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info(
			"failed to decode incoming json",
			zap.Error(err),
		)

		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var logDate string
	if req.LogDate.IsZero() {
		logDate = time.Now().Format("2006-01-02")
	} else {
		// Check if the logDate is in the future
		if req.LogDate.After(time.Now()) {
			http.Error(w, "Log date cannot be a future date", http.StatusBadRequest)
			return
		}

		logDate = req.LogDate.Format("2006-01-02")
	}

	exists, err := lib.DoesLogExistForTheDay(*userId, logDate)
	if err != nil {
		log.Info(
			"failed to determine user log's existence",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)
	}

	if *exists {
		http.Error(w, "Log already exists for this day", http.StatusConflict)
		return
	}

	bmr, err := lib.GetUserBmr(*userId)
	if err != nil {
		log.Info(
			"failed to get user bmr by id",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := lib.CreateUserLog(*userId, *bmr, logDate); err != nil {
		log.Info(
			"failed to create user log by id",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := Response{
		Code: http.StatusOK,
	}
	json.NewEncoder(w).Encode(&resp)
}

func GetCalorieLogs(w http.ResponseWriter, r *http.Request) {
	// Retrieve the token from the context
	tokenStr, ok := r.Context().Value(TokenContextKey).(string)
	if !ok {
		log.Info(
			"token not found in context",
		)

		// Token is not present in context
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userId, err := lib.ExtractUserIdFromToken(tokenStr)
	if err != nil {
		log.Info(
			"failed to get user id by token",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var userCalorieLogs *[]lib.UserCalorieLogs
	userCalorieLogs, err = lib.GetCalorieLogs(*userId)
	if err != nil {
		log.Info(
			"failed to get logs by user id",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := Response{
		Code: http.StatusOK,
		Data: userCalorieLogs,
	}
	json.NewEncoder(w).Encode(&resp)
}

type UpdateCalorieLogReq struct {
	CaloriesConsumed float64   `json:"calories_consumed"`
	CaloriesBurnt    float64   `json:"calories_burnt"`
	LogDate          time.Time `json:"log_date"`
}

func UpdateCalorieLogHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the token from the context
	tokenStr, ok := r.Context().Value(TokenContextKey).(string)
	if !ok {
		log.Info(
			"token not found in context",
		)

		// Token is not present in context
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userId, err := lib.ExtractUserIdFromToken(tokenStr)
	if err != nil {
		log.Info(
			"failed to get user id by token",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var req UpdateCalorieLogReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info(
			"failed to decode incoming json",
			zap.Error(err),
		)

		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	logDate := req.LogDate.Format("2006-01-02")

	logStatus, err := lib.CheckLogStatusByIdAndDate(*userId, logDate)
	if err != nil {
		log.Info(
			"failed to check log status by id and date",
			zap.String("userId", userId.String()),
			zap.String("logDate", logDate),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if *logStatus == "D" {
		http.Error(w, "Log is already completed", http.StatusConflict)
		return
	}

	if req.CaloriesBurnt != 0.00 {
		currValue, err := lib.FetchCaloriesBurntForTheDay(*userId, logDate)
		if err != nil {
			log.Info(
				"failed to fetch calories burnt by id and date",
				zap.String("userId", userId.String()),
				zap.String("logDate", logDate),
				zap.Error(err),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if *currValue+req.CaloriesBurnt < 0 {
			http.Error(w, "Resulting calories burnt can't be negative", http.StatusBadRequest)
			return
		}

		if err := lib.AddCaloriesBurntInTDEE(*userId, logDate, req.CaloriesBurnt); err != nil {
			log.Info(
				"failed to add burnt calories in tdee by id and date",
				zap.String("userId", userId.String()),
				zap.String("logDate", logDate),
				zap.Error(err),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if req.CaloriesConsumed != 0.00 {
		currValue, err := lib.FetchCaloriesConsumedForTheDay(*userId, logDate)
		if err != nil {
			log.Info(
				"failed to fetch calories consumed by id and date",
				zap.String("userId", userId.String()),
				zap.String("logDate", logDate),
				zap.Error(err),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if *currValue+req.CaloriesConsumed < 0 {
			http.Error(w, "Resulting calories consumed can't be negative", http.StatusBadRequest)
			return
		}
	}

	if err := lib.UpdateCalorieLog(*userId, logDate, req.CaloriesConsumed, req.CaloriesBurnt); err != nil {
		log.Info(
			"failed to update calorie log by id and date",
			zap.String("userId", userId.String()),
			zap.String("logDate", logDate),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := Response{
		Code: http.StatusOK,
	}
	json.NewEncoder(w).Encode(&resp)
}
