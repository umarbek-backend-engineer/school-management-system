package handler

import (
	"RESTAPI/internal/models"
	"RESTAPI/internal/repository/sqlconnections"
	"RESTAPI/pkg/utils"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/go-mail/mail/v2"
)

func GetExecHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error in converting the string id into int")
		return
	}

	var exec models.Exec

	foundExec, err := sqlconnections.GetExecDBHandler(w, id, exec)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(foundExec)
	if err != nil {
		log.Println("Error in encoding the responce into json: ", err)
		http.Error(w, "Error in enconding the responce", http.StatusInternalServerError)
		return
	}

}

func GetExecsHandler(w http.ResponseWriter, r *http.Request) {

	execslist, err := sqlconnections.GetExecsDBHandler(w, r)

	if err != nil {
		log.Println(err)
		return
	}

	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "Success",
		Count:  len(execslist),
		Data:   execslist,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error in encoding the responce into json: ", err)
		http.Error(w, "Error in enconding the responce", http.StatusInternalServerError)
		return
	}
}

func AddExecsHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newExecs []models.Exec

	err := json.NewDecoder(r.Body).Decode(&newExecs)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		log.Println("Error in decoding recieved data: ", err)
		return
	}

	for _, exec := range newExecs {
		val := reflect.ValueOf(exec)
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.Kind() == reflect.String && field.String() == "" {
				http.Error(w, "All field must be filled", http.StatusBadRequest)
				return
			}
		}
	}

	addedExecs, err := sqlconnections.PostExecsDBHandler(w, newExecs)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "Success",
		Count:  len(addedExecs),
		Data:   addedExecs,
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error in sending the response", http.StatusInternalServerError)
		log.Println("Error in sending the response: ", err)
		return
	}
}

func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println("Error in decoding json val: ", err)
		http.Error(w, "Invalid request payload ", http.StatusBadRequest)
		return
	}

	updatedExecs, err := sqlconnections.PatchExecsDBhandler(w, updates)
	if err != nil {
		log.Println(err)
		return
	}

	responce := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(updatedExecs),
		Data:   updatedExecs,
	}

	err = json.NewEncoder(w).Encode(responce)
	if err != nil {
		log.Println("Error in sending response", err)
		http.Error(w, "Error in sending response", http.StatusInternalServerError)
		return
	}
}

func PatchExecHandler(w http.ResponseWriter, r *http.Request) {

	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)

	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		log.Println("Error in decoding recieved data: ", err)
		return
	}

	var existingExec models.Exec

	existingExec, err = sqlconnections.PatchExecDBHandler(w, id, existingExec, updates)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string      `json:"status"`
		Data   models.Exec `json:"data"`
	}{
		Status: "Success",
		Data:   existingExec,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error in sending the response", http.StatusInternalServerError)
		log.Println("Error in sending the response: ", err)
		return
	}

}

func ExecDeleteHandler(w http.ResponseWriter, r *http.Request) {

	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	err = sqlconnections.DeleteExecDBHandler(w, id)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Success",
		ID:     id,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error in sending the response", http.StatusInternalServerError)
		log.Println("Error in sending the response: ", err)
		return
	}
}

func LogInHandler(w http.ResponseWriter, r *http.Request) {
	// Data validation
	var req models.Exec
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Error in decoding creds: ", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Username == "" || req.Password == "" {
		http.Error(w, "UserName and Password are required", http.StatusBadRequest)
		return
	}

	// Search for a user if user actually exists

	db, err := sqlconnections.ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	user := &models.Exec{}
	err = db.QueryRow("select id, first_name, last_name, email, username, password, inactive_status, role from execs where username = ?", req.Username).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.InacvtiveStatus,
		&user.Role,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		log.Println("Error: ", err)
		http.Error(w, "Database query error", http.StatusBadRequest)
		return
	}

	// Check if user is active

	if user.InacvtiveStatus {
		http.Error(w, "Account is inactive", http.StatusForbidden)
	}

	// Verify password

	err = utils.VerifyPassword(user.Password, w, req.Password)
	if err != nil {
		log.Println(err)
		return
	}

	// Generate JWT token
	tokenString, err := utils.Sign_Token(user.ID, user.Username, user.Role)
	if err != nil {
		log.Println("Error from JWT utils .go file: ", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	// Send token as a response or as a cookie

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 4380),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}
	json.NewEncoder(w).Encode(response)
}

func LogOuthandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Loged out successfylly"}`))
}

func UpdatePasswordhandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	db, err := sqlconnections.ConnectDb()
	if err != nil {
		log.Println("Error in retriving data from database: ", err)
		http.Error(w, "Internal Error in retriving data", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var req models.Exec_Update_password_request

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid payload", http.StatusInternalServerError)
		return
	}

	reqValue := reflect.ValueOf(req)

	for i := 0; i < reqValue.NumField(); i++ {
		value := reqValue.Field(i)
		if value.Kind() == reflect.String && value.String() == "" {
			http.Error(w, "Please enter password", http.StatusBadRequest)
			return
		}
	}

	var username string
	var password string
	var role string

	err = db.QueryRow("select username, password, role from execs where id = ?", id).Scan(&username, &password, &role)
	if err != nil {
		log.Println(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = utils.VerifyPassword(password, w, req.Current_Password)
	if err != nil {
		log.Println(err)
		return
	}

	if req.New_Password != req.Confurmation_Password {
		http.Error(w, "Password confirmation does not match the new password.", http.StatusBadRequest)
		return
	}

	encodedHash, err := utils.HashPassword(req.New_Password)
	if err != nil {
		fmt.Println(err)
		return
	}

	current_time := time.Now().Format(time.RFC3339)
	_, err = db.Exec("update execs set password = ?, password_changed_at = ? where id = ?", encodedHash, current_time, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update the password", http.StatusInternalServerError)
		return
	}

	idint, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return
	}
	tokenString, err := utils.Sign_Token(idint, username, role)
	if err != nil {
		log.Println(err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 4380),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}
	json.NewEncoder(w).Encode(response)

}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Enter your email", http.StatusBadRequest)
		return
	}

	db, err := sqlconnections.ConnectDb()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to connect to DB", http.StatusNotFound)
		return
	}
	defer db.Close()

	var exec models.Exec
	err = db.QueryRow("select id from execs where email = ?", req.Email).Scan(&exec.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	duration, err := strconv.Atoi(os.Getenv("RESET_TOKEN_EXP_DURATION"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Faild to send reset email", http.StatusInternalServerError)
		return
	}

	mins := time.Duration(duration)

	expiry := time.Now().Add(mins * time.Minute).Format(time.RFC3339)

	tokenbytes := make([]byte, 32)

	_, err = rand.Read(tokenbytes)
	if err != nil {
		log.Println(err)
		http.Error(w, "Faild to send reset email", http.StatusInternalServerError)
		return
	}

	log.Println("Token bytes: ", tokenbytes)
	token := hex.EncodeToString(tokenbytes)
	log.Println("Token String: ", token)

	hashedToken := sha256.Sum256(tokenbytes)
	log.Println("Hashed Token bytes: ", tokenbytes)

	TokenString := hex.EncodeToString(hashedToken[:])
	log.Println("Hashed Token bytes: ", TokenString)

	_, err = db.Exec("update execs set password_reset_token = ?, password_token_epires = ? where id = ?", TokenString, expiry, exec.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Faild to send reset email", http.StatusInternalServerError)
		return
	}

	//emailing

	resetURL := fmt.Sprintf("https://localhost:8080/execs/resetpassword/reset/%s", token)
	message := fmt.Sprintf("Forgot you password? Reset your password using the following link: \n%s\nIf you did not request a password reset, please ignore this email.\nThis link is only valid for %d minutes.", resetURL, duration)
	m := mail.NewMessage()
	m.SetHeader("From", "schooladmin@shool.edu") // replace with your own sender
	m.SetHeader("To", req.Email)
	m.SetHeader("Subject", "You password reset link")
	m.SetBody("text/plain", message)

	d := mail.NewDialer("localhost", 1025, "", "")
	err = d.DialAndSend(m)
	if err != nil {
		log.Println(err)
		http.Error(w, "Faild to send reset email", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Password reset link send to %s", req.Email)
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("resetcode")
	type request struct {
		New_Password    string `json:"new_password"`
		Confirm_Pasword string `json:"confirm_password"`
	}
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Payload", http.StatusBadRequest)
		return
	}

	if req.Confirm_Pasword == "" || req.New_Password == "" {
		http.Error(w, "Password is missing", http.StatusBadRequest)
		return
	}

	if req.Confirm_Pasword != req.New_Password {
		http.Error(w, "Passwords should match", http.StatusBadRequest)
		return
	}

	db, err := sqlconnections.ConnectDb()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to connect to DB", http.StatusNotFound)
		return
	}
	defer db.Close()

	var user models.Exec

	bytes, err := hex.DecodeString(token)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusNotFound)
		return
	}

	hashedToken := sha256.Sum256(bytes)
	heshedTokenString := hex.EncodeToString(hashedToken[:])

	query := "select id, email from execs where password_reset_token = ? and password_token_epires > ?"
	err = db.QueryRow(query, heshedTokenString, time.Now().Format(time.RFC3339)).Scan(&user.ID, &user.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid or Expired reset code", http.StatusNotFound)
		return
	}

	hashPassword, err := utils.HashPassword(req.New_Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update password", http.StatusNotFound)
		return
	}

	updateQuery := "update execs set password = ?, password_reset_token = null, password_token_epires = null, password_changed_at = ? where id = ?"
	_, err = db.Exec(updateQuery, hashPassword, time.Now().Format(time.RFC3339), user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update password", http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "Password reset succfully")
}
