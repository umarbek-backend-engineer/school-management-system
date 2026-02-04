package handler

import (
	"RESTAPI/internal/models"
	"RESTAPI/internal/repository/sqlconnections"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

func GetStudentHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error in converting the string id into int")
		return
	}

	var student models.Student

	foundStudent, err := sqlconnections.GetStudentDBHandler(w, id, student)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(foundStudent)
	if err != nil {
		log.Println("Error in encoding the responce into json: ", err)
		http.Error(w, "Error in enconding the responce", http.StatusInternalServerError)
		return
	}

}

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {

	page, limit := getPaginationparams(r)

	studentlist, TotalNumberStudents, err := sqlconnections.GetStudentsDBHandler(w, r, page, limit)

	if err != nil {
		log.Println(err)
		return
	}

	response := struct {
		Status   string           `json:"status"`
		Count    int              `json:"count"`
		Page     int              `json:"page"`
		PageSize int              `json:"page_size"`
		Data     []models.Student `json:"data"`
	}{
		Status:   "Success",
		Count:    TotalNumberStudents,
		Page:     page,
		PageSize: limit,
		Data:     studentlist,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error in encoding the responce into json: ", err)
		http.Error(w, "Error in enconding the responce", http.StatusInternalServerError)
		return
	}
}

func getPaginationparams(r *http.Request) (int, int) {
	pagestr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pagestr)
	if err != nil {
		log.Println(err)
		page = 1
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Println(err)
		limit = 10
	}

	return page, limit
}

func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newStudents []models.Student

	err := json.NewDecoder(r.Body).Decode(&newStudents)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		log.Println("Error in decoding recieved data: ", err)
		return
	}

	for _, student := range newStudents {
		val := reflect.ValueOf(student)
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.Kind() == reflect.String && field.String() == "" {
				http.Error(w, "All field must be filled", http.StatusBadRequest)
				return
			}
		}
	}

	addedStudents, err := sqlconnections.PostStudentsDBHandler(w, newStudents)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "Success",
		Count:  len(addedStudents),
		Data:   addedStudents,
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error in sending the response", http.StatusInternalServerError)
		log.Println("Error in sending the response: ", err)
		return
	}
}

func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {

	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	var updatedStudent models.Student
	err = json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		log.Println("Error in decoding recieved data: ", err)
		return
	}

	var existingStudent models.Student

	err = sqlconnections.UpdataStudentDBHandler(w, id, existingStudent, updatedStudent)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string         `json:"status"`
		Data   models.Student `json:"data"`
	}{
		Status: "Success",
		Data:   updatedStudent,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error in sending the response", http.StatusInternalServerError)
		log.Println("Error in sending the response: ", err)
		return
	}

}

func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println("Error in decoding json val: ", err)
		http.Error(w, "Invalid request payload ", http.StatusBadRequest)
		return
	}

	updatedStudent, err := sqlconnections.PatchStudentsDBhandler(w, updates)
	if err != nil {
		log.Println(err)
		return
	}

	responce := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(updatedStudent),
		Data:   updatedStudent,
	}

	err = json.NewEncoder(w).Encode(responce)
	if err != nil {
		log.Println("Error in sending response", err)
		http.Error(w, "Error in sending response", http.StatusInternalServerError)
		return
	}
}

func PatchStudentHandler(w http.ResponseWriter, r *http.Request) {

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

	var existingStudent models.Student

	existingStudent, err = sqlconnections.PatchStudentDBHandler(w, id, existingStudent, updates)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string         `json:"status"`
		Data   models.Student `json:"data"`
	}{
		Status: "Success",
		Data:   existingStudent,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error in sending the response", http.StatusInternalServerError)
		log.Println("Error in sending the response: ", err)
		return
	}

}

func StudentDeleteHandler(w http.ResponseWriter, r *http.Request) {

	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	err = sqlconnections.DeleteStudentDBHandler(w, id)
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

func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Println("Error in decoding ids: ", err)
		http.Error(w, "Error in decoding ids", http.StatusInternalServerError)
		return
	}

	deletedIDS, err := sqlconnections.DeleteStudentsDbHandler(w, ids)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
		Data   []int  `json:"data"`
	}{
		Status: "Successfully deleted",
		Count:  len(deletedIDS),
		Data:   deletedIDS,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error in sending the response", http.StatusInternalServerError)
		log.Println("Error in sending the response: ", err)
		return
	}
}

func DeleteAllStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var DeletedStudents []models.Student

	studentlist, err := sqlconnections.DeleteAllStudentsDBHandler(w, DeletedStudents)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(DeletedStudents)

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "Successfully Deleted",
		Count:  len(studentlist),
		Data:   studentlist,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
		return
	}
}
