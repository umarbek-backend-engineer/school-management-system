package router

import (
	"RESTAPI/internal/api/handler"
	"net/http"
)

func Router() *http.ServeMux {

	mux := http.NewServeMux()

	//------------------------------------------------------------------------------------------------------

	mux.HandleFunc("GET /teachers/", http.HandlerFunc(handler.GetTeachersHandler))
	mux.HandleFunc("POST /teachers/", http.HandlerFunc(handler.AddTeacherHandler))
	mux.HandleFunc("PATCH /teachers/", http.HandlerFunc(handler.PatchTeachersHandler))
	mux.HandleFunc("DELETE /teachers/", http.HandlerFunc(handler.DeleteTeachersHandler))
	mux.HandleFunc("DELETE /allteachers/", http.HandlerFunc(handler.DeleteAllTeachersHandler))

	mux.HandleFunc("GET /teachers/{id}", http.HandlerFunc(handler.GetTeacherHandler))
	mux.HandleFunc("PUT /teachers/{id}", http.HandlerFunc(handler.UpdateTeacherHandler))
	mux.HandleFunc("PATCH /teachers/{id}", http.HandlerFunc(handler.PatchTeacherHandler))
	mux.HandleFunc("DELETE /teachers/{id}", http.HandlerFunc(handler.TeacherDeleteHandler))

	mux.HandleFunc("GET /teachers/{id}/students/", http.HandlerFunc(handler.GetStudentsByTeacherID))
	mux.HandleFunc("GET /teachers/{id}/studentcount/", http.HandlerFunc(handler.GetstudnetCountForATeacher))

	//------------------------------------------------------------------------------------------------------

	mux.HandleFunc("GET /students/", http.HandlerFunc(handler.GetStudentsHandler))
	mux.HandleFunc("POST /students/", http.HandlerFunc(handler.AddStudentHandler))
	mux.HandleFunc("PATCH /students/", http.HandlerFunc(handler.PatchStudentsHandler))
	mux.HandleFunc("DELETE /students/", http.HandlerFunc(handler.DeleteStudentsHandler))
	mux.HandleFunc("DELETE /allstudents/", http.HandlerFunc(handler.DeleteAllStudentsHandler))

	mux.HandleFunc("GET /students/{id}", http.HandlerFunc(handler.GetStudentHandler))
	mux.HandleFunc("PUT /students/{id}", http.HandlerFunc(handler.UpdateStudentHandler))
	mux.HandleFunc("PATCH /students/{id}", http.HandlerFunc(handler.PatchStudentHandler))
	mux.HandleFunc("DELETE /students/{id}", http.HandlerFunc(handler.StudentDeleteHandler))

	//------------------------------------------------------------------------------------------------------

	mux.HandleFunc("GET /execs/", http.HandlerFunc(handler.GetExecsHandler))
	mux.HandleFunc("POST /execs/", http.HandlerFunc(handler.AddExecsHandler))
	mux.HandleFunc("PATCH /execs/", http.HandlerFunc(handler.PatchExecsHandler))

	mux.HandleFunc("GET /execs/{id}", http.HandlerFunc(handler.GetExecHandler))
	mux.HandleFunc("PATCH /execs/{id}", http.HandlerFunc(handler.PatchExecHandler))
	mux.HandleFunc("DELETE /execs/{id}", http.HandlerFunc(handler.ExecDeleteHandler))
	mux.HandleFunc("POST /execs/{id}/updatepassword", http.HandlerFunc(handler.UpdatePasswordhandler))

	mux.HandleFunc("POST /execs/login", http.HandlerFunc(handler.LogInHandler))
	mux.HandleFunc("POST /execs/logout", http.HandlerFunc(handler.LogOuthandler))
	mux.HandleFunc("POST /execs/forgotpassword", http.HandlerFunc(handler.ForgotPasswordHandler))
	mux.HandleFunc("POST /execs/resetpassword/reset/{resetcode}", http.HandlerFunc(handler.ResetPasswordHandler))

	//------------------------------------------------------------------------------------------------------

	return mux
}
