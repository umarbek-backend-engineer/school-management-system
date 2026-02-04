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

func GetTeacherDBHandler(w http.ResponseWriter, id int, teacher models.Teacher) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return teacher, err
	}
	defer db.Close()
	row := db.QueryRow("select id, first_name, last_name, email, class, subject from teachers where id = ?", id)
	err = row.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return teacher, err
		}
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return teacher, err
	}
	return teacher, nil
}

func GetTeachersDBHandler(w http.ResponseWriter, r *http.Request) ([]models.Teacher, error) {

	query := "select id, first_name, last_name, email, class, subject from teachers where 1=1"
	var args []interface{}

	query, args = utils.Getfilters(r, query, args)
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
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher(s) not found", http.StatusNotFound)
			return nil, err
		}
		fmt.Println("Error in geting muliple teachers", err)
		return nil, err
	}
	defer rows.Close()

	teacherlist := make([]models.Teacher, 0)

	for rows.Next() {
		var teacher models.Teacher
		err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Teacher not found", http.StatusNotFound)
				return nil, err
			}
			log.Println("Error occured: ", err)
			http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
			return nil, err
		}
		teacherlist = append(teacherlist, teacher)
	}
	return teacherlist, nil
}

func PostTeachersDBHandler(w http.ResponseWriter, newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into teachers(first_name, last_name, email, class, subject) values (?,?,?,?,?)")
	if err != nil {
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		log.Println("Error in preparing stmt in add teachers handler: ", err)
		return nil, err
	}

	addedTeachers := make([]models.Teacher, len(newTeachers))

	for i, v := range newTeachers {
		res, err := stmt.Exec(v.FirstName, v.LastName, v.Email, v.Class, v.Subject)
		if err != nil {
			http.Error(w, "Ooops something went wrong ", http.StatusInternalServerError)
			log.Println("Error in inserting data into database: ", err)
			return nil, err
		}
		lastid, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting last id", http.StatusInternalServerError)
			log.Println("Error in getting lastid in post teachers handler ", err)
			return nil, err
		}
		v.ID = int(lastid)

		addedTeachers[i] = v
	}
	return addedTeachers, nil
}

func UpdataTeacherDBHandler(w http.ResponseWriter, id int, existingTeacher models.Teacher, updatedTeacher models.Teacher) error {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return err
	}
	defer db.Close()
	err = db.QueryRow("select id, first_name, last_name, email, class, subject from teachers where id=?", id).Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Class,
		&existingTeacher.Subject,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return err
		}
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return err
	}

	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("update teachers set first_name = ?, last_name = ?, email = ?, class = ?, subject = ? where id = ?",
		updatedTeacher.FirstName,
		updatedTeacher.LastName,
		updatedTeacher.Email,
		updatedTeacher.Class,
		updatedTeacher.Subject,
		updatedTeacher.ID,
	)

	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Updating teacher", http.StatusInternalServerError)
		return err
	}
	return nil
}

func PatchTeachersDBhandler(w http.ResponseWriter, updates []map[string]interface{}) ([]models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin() // transaction is one or more sql commends to execute that all will ether will succeed nither fail
	if err != nil {
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return nil, err
	}
	var updatedTeacher []models.Teacher
	for _, v := range updates {
		idstr, ok := v["id"].(float64)
		if !ok {
			tx.Rollback()
			http.Error(w, "Invalid Teacher id", http.StatusBadRequest)
			return nil, err
		}

		id := int(idstr)

		var teacherFromDB models.Teacher
		err = db.QueryRow("select id, first_name, last_name, email, class, subject from teachers where id=?", id).Scan(
			&teacherFromDB.ID,
			&teacherFromDB.FirstName,
			&teacherFromDB.LastName,
			&teacherFromDB.Email,
			&teacherFromDB.Class,
			&teacherFromDB.Subject,
		)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				http.Error(w, "Teacher not found", http.StatusNotFound)
				return nil, err
			}
			log.Println("Error occured: ", err)
			http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
			return nil, err
		}
		// apply update using reflection
		teacherVal := reflect.ValueOf(&teacherFromDB).Elem()
		teacherType := teacherVal.Type()

		for k, v := range v {
			if k == "id" {
				continue
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := teacherVal.Field(i)
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
		_, err = tx.Exec("update teachers set first_name = ?, last_name = ?, email = ?, class = ?, subject = ? where id = ?",
			teacherFromDB.FirstName,
			teacherFromDB.LastName,
			teacherFromDB.Email,
			teacherFromDB.Class,
			teacherFromDB.Subject,
			teacherFromDB.ID,
		)
		updatedTeacher = append(updatedTeacher, teacherFromDB)

		if err != nil {
			log.Println("Error in executing transaction: ", err)
			http.Error(w, "Error in updating teacher", http.StatusInternalServerError)
			return nil, err
		}

	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Error in commiting transaction", http.StatusInternalServerError)
		log.Println("Error in commiting transaction: ", err)
		return nil, err
	}
	return updatedTeacher, nil
}

func PatchTeacherDBHandler(w http.ResponseWriter, id int, existingTeacher models.Teacher, updates map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return err
	}
	defer db.Close()
	err = db.QueryRow("select id, first_name, last_name, email, class, subject from teachers where id=?", id).Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Class,
		&existingTeacher.Subject,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return err
		}
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return err
	}

	// for k, v := range updates {
	// 	switch k {
	// 	case "first_name":
	// 		existingTeacher.FirstName = v.(string)
	// 	case "last_name":
	// 		existingTeacher.LastName = v.(string)
	// 	case "email":
	// 		existingTeacher.Email = v.(string)
	// 	case "class":
	// 		existingTeacher.Class = v.(string)
	// 	case "subject":
	// 		existingTeacher.Subject = v.(string)
	// 	}
	// }

	//apply update using reflect
	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {
					teacherVal.Field(i).Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("update teachers set first_name = ?, last_name = ?, email = ?, class = ?, subject = ? where id = ?",
		existingTeacher.FirstName,
		existingTeacher.LastName,
		existingTeacher.Email,
		existingTeacher.Class,
		existingTeacher.Subject,
		existingTeacher.ID,
	)

	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Updating teacher", http.StatusInternalServerError)
		return err
	}
	return nil
}

func DeleteTeacherDBHandler(w http.ResponseWriter, id int) error {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	res, err := db.Exec("delete from teachers where id=?", id)
	if err != nil {
		http.Error(w, "Error occured", http.StatusInternalServerError)
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Error retriving deleted teacher", http.StatusInternalServerError)
		log.Println(err)
		return err
	}

	if rowAffected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return sql.ErrNoRows
	}
	return nil
}

func DeleteTeachersDbHandler(w http.ResponseWriter, ids []int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return nil, err
	}
	defer db.Close()

	stmt, err := tx.Prepare("delete from teachers where id = ?")
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
				http.Error(w, "Teacher not found", http.StatusBadRequest)
				return nil, err
			}
			http.Error(w, "Error occured", http.StatusInternalServerError)
			return nil, err
		}
		rowAffected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error retriving deleted teacher", http.StatusInternalServerError)
			log.Println(err)
			return nil, err
		}

		if rowAffected == 0 {
			tx.Rollback()
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return nil, sql.ErrNoRows
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

func GetStudentsByTeacherIDDB(w http.ResponseWriter, r *http.Request, id int) ([]models.Student, error) {

	var studentList []models.Student

	query := "select id, first_name, last_name, email, class from students where class = ( select class from teachers where id = ?)"

	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return nil, err
		}
		fmt.Println("Error in geting students for a teacher", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Student not found", http.StatusNotFound)
				return nil, err
			}
			log.Println("Error occured: ", err)
			http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
			return nil, err
		}
		studentList = append(studentList, student)
	}
	return studentList, nil
}

func DeleteAllTeacherDBHandler(w http.ResponseWriter, teacherList []models.Teacher) ([]models.Teacher, error) {

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

	row, err := tx.Query("select id, first_name, last_name, email, class, subject from teachers where 1=1")
	if err != nil {
		http.Error(w, "Error in retriving data", http.StatusBadRequest)
		tx.Rollback()
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var teachers models.Teacher
		err = row.Scan(&teachers.ID, &teachers.FirstName, &teachers.LastName, &teachers.Email, &teachers.Class, &teachers.Subject)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := row.Err(); err != nil {
			tx.Rollback()
			return nil, err
		}
		teacherList = append(teacherList, teachers)
	}

	_, err = tx.Exec("delete from teachers where 1 = 1 ")
	if err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "Cannot delete or update a parent row: a foreign key constraint fails (`school`.`students`, CONSTRAINT `students_ibfk_1` FOREIGN KEY (`class`) REFERENCES `teachers` (`class`))") {
			http.Error(w, "First, related student must be deleted!", http.StatusBadRequest)
			return nil, err
		}
		http.Error(w, "Error in deleteing teachers", http.StatusInternalServerError)
		return nil, err
	}

	tx.Commit()

	return teacherList, nil
}

func GetStudentCountTeacherDbHandler(id string) (int, error) {

	var studentcount int

	db, err := ConnectDb()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	query := "select count(*) from students where class = (select class from teachers where id = ?)"

	err = db.QueryRow(query, id).Scan(&studentcount)
	if err != nil {
		return 0, err
	}

	return studentcount, nil

}
