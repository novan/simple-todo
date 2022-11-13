package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Todos []Todo

var mainDB *sql.DB

func main() {

	db, errOpenDB := sql.Open("sqlite3", "todo.db")
	checkErr(errOpenDB)
	mainDB = db

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/todos", getAll)
	e.GET("/todos/:id", getByID)
	e.PUT("/todos/:id", updateByID)

	// Start server
	err := e.Start(":12345")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func getAll(c echo.Context) error {
	rows, err := mainDB.Query("SELECT * FROM todos")
	checkErr(err)
	var todos Todos
	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID, &todo.Name)
		checkErr(err)
		todos = append(todos, todo)
	}
	return c.JSON(http.StatusOK, todos)
}

func getByID(c echo.Context) error {
	id := c.Param("id")
	stmt, err := mainDB.Prepare(" SELECT * FROM todos where id = ?")
	checkErr(err)
	rows, errQuery := stmt.Query(id)
	checkErr(errQuery)
	var todo Todo
	for rows.Next() {
		err = rows.Scan(&todo.ID, &todo.Name)
		checkErr(err)
	}
	return c.JSON(http.StatusOK, todo)
}

func updateByID(c echo.Context) error {
	name := c.FormValue("name")
	id := c.Param("id")
	var todo Todo
	ID, _ := strconv.ParseInt(id, 10, 0)
	todo.ID = ID
	todo.Name = name
	stmt, err := mainDB.Prepare("UPDATE todos SET name = ? WHERE id = ?")
	checkErr(err)
	result, errExec := stmt.Exec(todo.Name, todo.ID)
	checkErr(errExec)
	rowAffected, errLast := result.RowsAffected()
	checkErr(errLast)
	if rowAffected > 0 {
		return c.JSON(http.StatusOK, todo)
	} else {
		return c.String(http.StatusOK, fmt.Sprintf("{row_affected=%d}", rowAffected))
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
