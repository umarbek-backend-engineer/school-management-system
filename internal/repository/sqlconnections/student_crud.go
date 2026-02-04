package sqlconnections

import (
	"RESTAPI/internal/models"
	"RESTAPI/pkg/utils"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func GetStudentDBHandler(w http.ResponseWriter, id int, student models.Student) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return student, err
	}

	row := db.QueryRow("select id, first_name, last_name, email, class from students where id = ?", id)
	err = row.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
			return student, err
		}
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return student, err
	}
	return student, nil
}

func GetStudentsDBHandler(w http.ResponseWriter, r *http.Request, page, limit int) ([]models.Student, int, error) {

	query := "select id, first_name, last_name, email, class from students where 1=1"
	var args []interface{}

	query, args = utils.Getfilters(r, query, args)

	//ading pagination
	offset := (page - 1) * limit
	query += " limit ? offset ?"
	args = append(args, limit, offset)

	sortParams := r.URL.Query()["sortby"]

	if len(sortParams) > 0 {
		query += " ORDER BY "
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !utils.IsValidSortField(field) || !utils.IsValidSortOrder(order) {
				continue
			}
			if i > 0 {
				query += ","
			}
			query += " " + field + " " + order
		}
	}

	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return nil, 0, err
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student(s) not found", http.StatusNotFound)
			return nil, 0, err
		}
		fmt.Println("Error in geting muliple students", err)
		return nil, 0, err
	}
	defer rows.Close()

	studentlist := make([]models.Student, 0)

	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Student not found", http.StatusNotFound)
				return nil, 0, err
			}
			log.Println("Error occured: ", err)
			http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
			return nil, 0, err
		}
		studentlist = append(studentlist, student)
	}

	var totalNstudents int
	err = db.QueryRow("select count(*) from students").Scan(&totalNstudents)
	if err != nil {
		log.Println(err)
	}
	return studentlist, totalNstudents, nil
}

func PostStudentsDBHandler(w http.ResponseWriter, newStudent []models.Student) ([]models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into students (first_name, last_name, email, class) values (?,?,?,?)")
	if err != nil {
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		log.Println("Error in preparing stmt in add students handler: ", err)
		return nil, err
	}

	addedStudents := make([]models.Student, len(newStudent))

	for i, v := range newStudent {
		res, err := stmt.Exec(v.FirstName, v.LastName, v.Email, v.Class)
		if err != nil {
			http.Error(w, "Ooops something went wrong ", http.StatusInternalServerError)
			log.Println("Error in inserting data into database: ", err)
			return nil, err
		}
		lastid, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting last id", http.StatusInternalServerError)
			log.Println("Error in getting lastid in post students handler ", err)
			return nil, err
		}
		v.ID = int(lastid)

		addedStudents[i] = v
	}
	return addedStudents, nil
}

func UpdataStudentDBHandler(w http.ResponseWriter, id int, existingStudent models.Student, updatedStudent models.Student) error {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return err
	}
	defer db.Close()
	err = db.QueryRow("select id, first_name, last_name, email, class from students where id=?", id).Scan(
		&existingStudent.ID,
		&existingStudent.FirstName,
		&existingStudent.LastName,
		&existingStudent.Email,
		&existingStudent.Class,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
			return err
		}
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong in retriving data", http.StatusInternalServerError)
		return err
	}

	updatedStudent.ID = existingStudent.ID
	_, err = db.Exec("update students set first_name = ?, last_name = ?, email = ?, class = ? where id = ?",
		updatedStudent.FirstName,
		updatedStudent.LastName,
		updatedStudent.Email,
		updatedStudent.Class,
		updatedStudent.ID,
	)

	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Updating student", http.StatusInternalServerError)
		return err
	}
	return nil
}

func PatchStudentsDBhandler(w http.ResponseWriter, updates []map[string]interface{}) ([]models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin() // transaction is one or more sql commends to execute that all will ether will succeed nither fail
	if err != nil {
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return nil, err
	}
	var updateStudent []models.Student
	for _, v := range updates {
		idstr, ok := v["id"].(float64)
		if !ok {
			tx.Rollback()
			http.Error(w, "Invalid Student id", http.StatusBadRequest)
			return nil, err
		}

		id := int(idstr)

		var StudentFromDb models.Student
		err = db.QueryRow("select id, first_name, last_name, email, class from students where id=?", id).Scan(
			&StudentFromDb.ID,
			&StudentFromDb.FirstName,
			&StudentFromDb.LastName,
			&StudentFromDb.Email,
			&StudentFromDb.Class,
		)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				http.Error(w, "Student not found", http.StatusNotFound)
				return nil, err
			}
			log.Println("Error occured: ", err)
			http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
			return nil, err
		}
		// apply update using reflection
		studentVal := reflect.ValueOf(&StudentFromDb).Elem()
		studentType := studentVal.Type()

		for k, v := range v {
			if k == "id" {
				continue
			}
			for i := 0; i < studentVal.NumField(); i++ {
				field := studentType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := studentVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							log.Printf("Cannot convert %v to %v\n", val.Type(), fieldVal.Type())
							return nil, err
						}
					}
					break
				}
			}
		}
		_, err = tx.Exec("update students set first_name = ?, last_name = ?, email = ?, class = ? where id = ?",
			StudentFromDb.FirstName,
			StudentFromDb.LastName,
			StudentFromDb.Email,
			StudentFromDb.Class,
			StudentFromDb.ID,
		)
		updateStudent = append(updateStudent, StudentFromDb)

		if err != nil {
			log.Println("Error in executing transaction: ", err)
			http.Error(w, "Error in updating Student", http.StatusInternalServerError)
			return nil, err
		}

	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Error in commiting transaction", http.StatusInternalServerError)
		log.Println("Error in commiting transaction: ", err)
		return nil, err
	}
	return updateStudent, nil
}

func PatchStudentDBHandler(w http.ResponseWriter, id int, existingStudent models.Student, updates map[string]interface{}) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return existingStudent, err
	}
	defer db.Close()
	err = db.QueryRow("select id, first_name, last_name, email, class from students where id=?", id).Scan(
		&existingStudent.ID,
		&existingStudent.FirstName,
		&existingStudent.LastName,
		&existingStudent.Email,
		&existingStudent.Class,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
			return existingStudent, err
		}
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return existingStudent, err
	}

	//apply update using reflect
	StudentVal := reflect.ValueOf(&existingStudent).Elem()
	studentType := StudentVal.Type()

	for k, v := range updates {
		for i := 0; i < StudentVal.NumField(); i++ {
			field := studentType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if StudentVal.Field(i).CanSet() {
					StudentVal.Field(i).Set(reflect.ValueOf(v).Convert(StudentVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("update students set first_name = ?, last_name = ?, email = ?, class = ? where id = ?",
		existingStudent.FirstName,
		existingStudent.LastName,
		existingStudent.Email,
		existingStudent.Class,
		existingStudent.ID,
	)

	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Updating student", http.StatusInternalServerError)
		return existingStudent, err
	}
	return existingStudent, nil
}

func DeleteStudentDBHandler(w http.ResponseWriter, id int) error {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	res, err := db.Exec("delete from students where id=?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "student not found", http.StatusBadRequest)
			return err
		}
		http.Error(w, "Error occured", http.StatusInternalServerError)
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Error retriving deleted student", http.StatusInternalServerError)
		log.Println(err)
		return err
	}

	if rowAffected == 0 {
		http.Error(w, "Student not found", http.StatusNotFound)
		return err
	}
	return nil
}

func DeleteStudentsDbHandler(w http.ResponseWriter, ids []int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Invalid Students ID", http.StatusBadRequest)
		return nil, err
	}

	stmt, err := tx.Prepare("delete from students where id = ?")
	if err != nil {
		log.Println("Error in preparing stmt: ", err)
		http.Error(w, "Error in preparing stmt", http.StatusInternalServerError)
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	deletedIDS := []int{}

	for _, v := range ids {
		res, err := stmt.Exec(v)

		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				http.Error(w, "Students not found", http.StatusBadRequest)
				return nil, err
			}
			http.Error(w, "Error occured", http.StatusInternalServerError)
			return nil, err
		}
		rowAffected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error retriving deleted Student", http.StatusInternalServerError)
			log.Println(err)
			return nil, err
		}

		if rowAffected == 0 {
			tx.Rollback()
			http.Error(w, "Students not found", http.StatusNotFound)
			return nil, err
		}

		deletedIDS = append(deletedIDS, v)

	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error in commiting transaction")
		http.Error(w, "Error in commiting transaction", http.StatusInternalServerError)
		return nil, err
	}
	return deletedIDS, nil
}

func DeleteAllStudentsDBHandler(w http.ResponseWriter, studentList []models.Student) ([]models.Student, error) {

	db, err := ConnectDb()
	if err != nil {
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	row, err := tx.Query("select id, first_name, last_name, email, class from students where 1=1")
	if err != nil {
		http.Error(w, "Error in retriving data", http.StatusBadRequest)
		tx.Rollback()
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var students models.Student
		err = row.Scan(&students.ID, &students.FirstName, &students.LastName, &students.Email, &students.Class)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := row.Err(); err != nil {
			tx.Rollback()
			return nil, err
		}
		studentList = append(studentList, students)
	}

	_, err = tx.Exec("delete from students where 1 = 1 ")
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error in deleteing students", http.StatusInternalServerError)
		return nil, err
	}

	tx.Commit()

	return studentList, nil
}
