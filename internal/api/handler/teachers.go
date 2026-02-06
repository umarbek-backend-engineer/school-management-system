package handler

import (
	"RESTAPI/internal/models"
	"RESTAPI/internal/repository/sqlconnections"
	"RESTAPI/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

var mutex = &sync.Mutex{}

// GetTeacherHandler retrieves a teacher by ID
// @Summary Get teacher by ID
// @Description Get a specific teacher by their ID
// @Tags Teachers
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} models.Teacher
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /teachers/{id} [get]
// @Security Bearer
func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error in converting the string id into int")
		return
	}

	var teacher models.Teacher

	foundteacher, err := sqlconnections.GetTeacherDBHandler(w, id, teacher)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(foundteacher)
	if err != nil {
		log.Println("Error in encoding the responce into json: ", err)
		http.Error(w, "Error in enconding the responce", http.StatusInternalServerError)
		return
	}

}

// GetTeachersHandler retrieves all teachers
// @Summary Get all teachers
// @Description Get a list of all teachers
// @Tags Teachers
// @Produce json
// @Success 200 {object} map[string]interface{} "List of teachers"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /teachers/ [get]
// @Security Bearer
func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {

	teacherlist, err := sqlconnections.GetTeachersDBHandler(w, r)

	if err != nil {
		log.Println(err)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(teacherlist),
		Data:   teacherlist,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error in encoding the responce into json: ", err)
		http.Error(w, "Error in enconding the responce", http.StatusInternalServerError)
		return
	}
}

// AddTeacherHandler adds new teachers
// @Summary Add new teachers
// @Description Create one or multiple new teachers
// @Tags Teachers
// @Accept json
// @Produce json
// @Param teachers body []models.Teacher true "Teacher data"
// @Success 201 {object} map[string]interface{} "Teachers created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /teachers/ [post]
// @Security Bearer
func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []models.Teacher

	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		log.Println("Error in decoding recieved data: ", err)
		return
	}

	for _, teacher := range newTeachers {
		val := reflect.ValueOf(teacher)
		fmt.Println(val)
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.Kind() == reflect.String && field.String() == "" {
				http.Error(w, "All field must be filled", http.StatusBadRequest)
				return
			}
		}
	}

	addedTeachers, err := sqlconnections.PostTeachersDBHandler(w, newTeachers)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error in sending the response", http.StatusInternalServerError)
		log.Println("Error in sending the response: ", err)
		return
	}
}

// UpdateTeacherHandler updates a teacher by ID
// @Summary Update teacher
// @Description Update a specific teacher by their ID
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher ID"
// @Param teacher body models.Teacher true "Teacher data to update"
// @Success 200 {object} map[string]interface{} "Teacher updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /teachers/{id} [put]
// @Security Bearer
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {

	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		log.Println("Error in decoding recieved data: ", err)
		return
	}

	var existingTeacher models.Teacher

	err = sqlconnections.UpdataTeacherDBHandler(w, id, existingTeacher, updatedTeacher)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string         `json:"status"`
		Data   models.Teacher `json:"data"`
	}{
		Status: "Success",
		Data:   updatedTeacher,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error in sending the response", http.StatusInternalServerError)
		log.Println("Error in sending the response: ", err)
		return
	}

}

func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println("Error in decoding json val: ", err)
		http.Error(w, "Invalid request payload ", http.StatusBadRequest)
		return
	}

	updatedTeacher, err := sqlconnections.PatchTeachersDBhandler(w, updates)
	if err != nil {
		log.Println(err)
		return
	}

	responce := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(updatedTeacher),
		Data:   updatedTeacher,
	}

	err = json.NewEncoder(w).Encode(responce)
	if err != nil {
		log.Println("Error in sending response", err)
		http.Error(w, "Error in sending response", http.StatusInternalServerError)
		return
	}
}

// PatchTeachersHandler patches multiple teachers
// @Summary Patch multiple teachers
// @Description Partially update multiple teachers
// @Tags Teachers
// @Accept json
// @Produce json
// @Param updates body []map[string]interface{} true "Teacher updates"
// @Success 200 {object} map[string]interface{} "Teachers patched successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /teachers/ [patch]
// @Security Bearer// PatchTeacherHandler patches a specific teacher
// @Summary Patch teacher
// @Description Partially update a specific teacher by ID
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher ID"
// @Param updates body map[string]interface{} true "Teacher updates"
// @Success 200 {object} map[string]interface{} "Teacher patched successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /teachers/{id} [patch]
// @Security Bearer
func PatchTeacherHandler(w http.ResponseWriter, r *http.Request) {

	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)

	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		log.Println("Error in decoding recieved data: ", err)
		return
	}

	var existingTeacher models.Teacher

	err = sqlconnections.PatchTeacherDBHandler(w, id, existingTeacher, updates)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string         `json:"status"`
		Data   models.Teacher `json:"data"`
	}{
		Status: "Success",
		Data:   existingTeacher,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error in sending the response", http.StatusInternalServerError)
		log.Println("Error in sending the response: ", err)
		return
	}

}

// TeacherDeleteHandler deletes a specific teacher
// @Summary Delete teacher
// @Description Delete a specific teacher by their ID
// @Tags Teachers
// @Param id path int true "Teacher ID"
// @Success 204 "Teacher deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /teachers/{id} [delete]
// @Security Bearer
func TeacherDeleteHandler(w http.ResponseWriter, r *http.Request) {

	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	err = sqlconnections.DeleteTeacherDBHandler(w, id)
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

// DeleteTeachersHandler deletes multiple teachers
// @Summary Delete multiple teachers
// @Description Delete multiple teachers by their IDs
// @Tags Teachers
// @Accept json
// @Produce json
// @Param ids body []int true "Teacher IDs to delete"
// @Success 200 {object} map[string]interface{} "Teachers deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /teachers/ [delete]
// @Security Bearer
func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Println("Error in decoding ids: ", err)
		http.Error(w, "Error in decoding ids", http.StatusInternalServerError)
		return
	}

	deletedIDS, err := sqlconnections.DeleteTeachersDbHandler(w, ids)
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

// GetStudentsByTeacherID retrieves all students for a specific teacher
// @Summary Get students by teacher
// @Description Get a list of all students taught by a specific teacher
// @Tags Teachers,Students
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} map[string]interface{} "List of students"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /teachers/{id}/students/ [get]
// @Security Bearer
func GetStudentsByTeacherID(w http.ResponseWriter, r *http.Request) {

	_, err := utils.AuthorizeUsers(r.Context().Value("user_role").(string), "exec")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	students, err := sqlconnections.GetStudentsByTeacherIDDB(w, r, id)

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "Success",
		Count:  len(students),
		Data:   students,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error in sending the response", http.StatusInternalServerError)
		log.Println("Error in sending the response: ", err)
		return
	}
}

// DeleteAllTeachersHandler deletes all teachers
// @Summary Delete all teachers
// @Description Delete all teachers from the system
// @Tags Teachers
// @Success 200 {object} map[string]interface{} "All teachers deleted successfully"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /allteachers/ [delete]
// @Security Bearer
func DeleteAllTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var DeletedTeachers []models.Teacher

	teacherlist, err := sqlconnections.DeleteAllTeacherDBHandler(w, DeletedTeachers)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(DeletedTeachers)

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Successfully Deleted",
		Count:  len(teacherlist),
		Data:   teacherlist,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
		return
	}
}

// GetstudnetCountForATeacher retrieves the count of students for a teacher
// @Summary Get student count for teacher
// @Description Get the number of students taught by a specific teacher
// @Tags Teachers,Students
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} map[string]interface{} "Student count"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /teachers/{id}/studentcount/ [get]
// @Security Bearer
func GetstudnetCountForATeacher(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	studentCount, err := sqlconnections.GetStudentCountTeacherDbHandler(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error in retriving data from db", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}{
		Status: "Success",
		Count:  studentCount,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
		return
	}
}
