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

func GetExecDBHandler(w http.ResponseWriter, id int, exec models.Exec) (models.Exec, error) {

	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return exec, err
	}

	row := db.QueryRow("select id, first_name, last_name, email, username, user_created_at, inactive_status, role from execs where id = ?", id)
	err = row.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.UserCreatedAT, &exec.InacvtiveStatus, &exec.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Exec not found", http.StatusNotFound)
			return exec, err
		}
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return exec, err
	}
	return exec, nil
}

func GetExecsDBHandler(w http.ResponseWriter, r *http.Request) ([]models.Exec, error) {

	query := "select id, first_name, last_name, email, username, user_created_at, inactive_status, role from execs where 1=1"
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

	rows, err := db.Query(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Exec(s) not found", http.StatusNotFound)
			return nil, err
		}
		fmt.Println("Error in geting muliple execs", err)
		return nil, err
	}
	defer rows.Close()

	execslist := make([]models.Exec, 0)

	for rows.Next() {
		var exec models.Exec
		err = rows.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.UserCreatedAT, &exec.InacvtiveStatus, &exec.Role)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "exec not found", http.StatusNotFound)
				return nil, err
			}
			log.Println("Error occured: ", err)
			http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
			return nil, err
		}
		execslist = append(execslist, exec)
	}
	return execslist, nil
}

func PostExecsDBHandler(w http.ResponseWriter, newExec []models.Exec) ([]models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into execs (first_name, last_name, email, username, password, inactive_status, role) values (?,?,?,?,?,?,?)")
	if err != nil {
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		log.Println("Error in preparing stmt in add execs handler: ", err)
		return nil, err
	}

	addedExecs := make([]models.Exec, len(newExec))

	for i, v := range newExec {
		encodedhash, err := utils.HashPassword(v.Password)
		if err != nil {
			return nil, err
		}
		v.Password = encodedhash

		res, err := stmt.Exec(v.FirstName, v.LastName, v.Email, v.Username, v.Password, v.InacvtiveStatus, v.Role)
		if err != nil {
			http.Error(w, "Ooops something went wrong ", http.StatusInternalServerError)
			log.Println("Error in inserting data into database: ", err)
			return nil, err
		}
		lastid, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting last id", http.StatusInternalServerError)
			log.Println("Error in getting lastid in post execs handler ", err)
			return nil, err
		}
		v.ID = int(lastid)

		addedExecs[i] = v
	}
	return addedExecs, nil
}

func PatchExecsDBhandler(w http.ResponseWriter, updates []map[string]interface{}) ([]models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin() // transaction is one or more sql commends to execute that all will ether will succeed nither fail
	if err != nil {
		http.Error(w, "Invalid Exec ID", http.StatusBadRequest)
		return nil, err
	}
	var updateExec []models.Exec
	for _, v := range updates {
		idstr, ok := v["id"].(float64)
		if !ok {
			tx.Rollback()
			http.Error(w, "Invalid Exec id", http.StatusBadRequest)
			return nil, err
		}

		id := int(idstr)

		var ExecFromDb models.Exec
		err = db.QueryRow("select id, first_name, last_name, email, username, role from execs where id=?", id).Scan(
			&ExecFromDb.ID,
			&ExecFromDb.FirstName,
			&ExecFromDb.LastName,
			&ExecFromDb.Email,
			&ExecFromDb.Username,
			&ExecFromDb.Role,
		)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				http.Error(w, "exec not found", http.StatusNotFound)
				return nil, err
			}
			log.Println("Error occured: ", err)
			http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
			return nil, err
		}
		// apply update using reflection
		execVal := reflect.ValueOf(&ExecFromDb).Elem()
		execType := execVal.Type()

		for k, v := range v {
			if k == "id" {
				continue
			}
			for i := 0; i < execVal.NumField(); i++ {
				field := execType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := execVal.Field(i)
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
		_, err = tx.Exec("update execs set first_name = ?, last_name = ?, email = ?, username = ?,  role = ? where id = ?",
			ExecFromDb.FirstName,
			ExecFromDb.LastName,
			ExecFromDb.Email,
			ExecFromDb.Username,
			ExecFromDb.Role,
			ExecFromDb.ID,
		)
		updateExec = append(updateExec, ExecFromDb)

		if err != nil {
			log.Println("Error in executing transaction: ", err)
			http.Error(w, "Error in updating exec", http.StatusInternalServerError)
			return nil, err
		}

	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Error in commiting transaction", http.StatusInternalServerError)
		log.Println("Error in commiting transaction: ", err)
		return nil, err
	}
	return updateExec, nil
}

func PatchExecDBHandler(w http.ResponseWriter, id int, existingExec models.Exec, updates map[string]interface{}) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return existingExec, err
	}
	defer db.Close()
	err = db.QueryRow("select id, first_name, last_name, email, username, user_created_at, inactive_status, role from execs where id=?", id).Scan(
		&existingExec.ID,
		&existingExec.FirstName,
		&existingExec.LastName,
		&existingExec.Email,
		&existingExec.Username,
		&existingExec.UserCreatedAT,
		&existingExec.InacvtiveStatus,
		&existingExec.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Exec not found", http.StatusNotFound)
			return existingExec, err
		}
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong", http.StatusInternalServerError)
		return existingExec, err
	}

	//apply update using reflect
	ExecVal := reflect.ValueOf(&existingExec).Elem()
	execType := ExecVal.Type()

	for k, v := range updates {
		for i := 0; i < ExecVal.NumField(); i++ {
			field := execType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if ExecVal.Field(i).CanSet() {
					ExecVal.Field(i).Set(reflect.ValueOf(v).Convert(ExecVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("update Execs set first_name = ?, last_name = ?, email = ?, username = ?, role = ? where id = ?",
		existingExec.FirstName,
		existingExec.LastName,
		existingExec.Email,
		existingExec.Username,
		existingExec.Role,
		existingExec.ID,
	)

	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Updating execs", http.StatusInternalServerError)
		return existingExec, err
	}
	return existingExec, nil
}

func DeleteExecDBHandler(w http.ResponseWriter, id int) error {
	db, err := ConnectDb()
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, "Ooops something went wrong: Unable to connect to DB", http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	res, err := db.Exec("delete from execs where id=?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "exec not found", http.StatusBadRequest)
			return err
		}
		http.Error(w, "Error occured", http.StatusInternalServerError)
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Error retriving deleted exec", http.StatusInternalServerError)
		log.Println(err)
		return err
	}

	if rowAffected == 0 {
		http.Error(w, "Exec not found", http.StatusNotFound)
		return err
	}
	return nil
}
