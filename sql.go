package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type PostData struct {
	Id      string
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"date"`
}

type ImageServ struct {
	Id    string
	Image []byte `json:"image"`
}

type Db struct {
	Tables    map[string]map[string]string
	DbName    string
	DbPath    string
	PostD     PostData
	TableName string
	FetchInfo string
	ImageS    ImageServ
}

const (
	username = "root"
	password = "root"
	hostname = "localhost:3306"
	dbname   = "eska"
)

func uuid4SQL() string {
	/*
		Генератор уникальных id
	*/

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",

		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)

}

func dbConnection() (*sql.DB, error) {
	var dbFirstName = fmt.Sprintf("%s:%s@tcp(%s)/", username, password, hostname)
	db, err := sql.Open("mysql", dbFirstName)
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, `CREATE DATABASE IF NOT EXISTS eska`)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return nil, err
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return nil, err
	}
	log.Printf("rows affected %d\n", no)

	db.Close()
	db, err = sql.Open("mysql", dsn())
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}
	//defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
	return db, nil
}

func (c Db) AddPost() error {
	/*
		Добавляет пост в БД
	*/

	db, dateBaseError := dbConnection()
	if dateBaseError != nil {
		panic(dateBaseError)
	}

	var ErrorAddInfo error

	records := `INSERT INTO posts VALUES (?, ?, ?, ?)`
	query, prepareError := db.Prepare(records)
	if prepareError != nil {
		ErrorAddInfo = prepareError
	}

	_, execError := query.Exec(c.PostD.Id, c.PostD.Title, c.PostD.Content, c.PostD.Date)
	if execError != nil {
		ErrorAddInfo = execError
	}
	return ErrorAddInfo
}

func (c Db) AddImage() error {
	/*
		Добавляет пост в БД
	*/

	db, dateBaseError := dbConnection()
	if dateBaseError != nil {
		panic(dateBaseError)
	}

	var ErrorAddInfo error

	records := `INSERT INTO post_images VALUES (?, ?)`
	query, prepareError := db.Prepare(records)
	if prepareError != nil {
		ErrorAddInfo = prepareError
	}

	_, execError := query.Exec(c.ImageS.Id, c.ImageS.Image)
	if execError != nil {
		ErrorAddInfo = execError
	}
	return ErrorAddInfo
}

func (c Db) ChangePost() error {
	/*
		Добавляет пост в БД
	*/

	db, dateBaseError := dbConnection()
	if dateBaseError != nil {
		panic(dateBaseError)
	}

	var ErrorAddInfo error
	//
	//records := `UPDATE posts SET Title = ?, Content = ? WHERE Id = ?`
	//_, execError := db.Exec(records, c.PostD.Title, c.PostD.Content, c.PostD.Id)

	//...
	_, execError := db.Exec("update  posts set content = ?, date = ? where id = ?", c.PostD.Title, c.PostD.Content, c.PostD.Id) //ебанутая хуйня которая парашно работает не обращать внимание на логику её тут нет

	if execError != nil {
		ErrorAddInfo = execError
	}
	return ErrorAddInfo
}

func (c Db) removeInfo() (bool, error) {
	/*
		Удаляет данные из БД
	*/

	db, dateBaseError := dbConnection()
	if dateBaseError != nil {
		return false, nil
	}

	var IdDb = c.PostD.Id
	var deleteReq = fmt.Sprintf("DELETE FROM " + c.TableName + " WHERE Id = '" + IdDb + "'")
	_, execError := db.Exec(deleteReq)

	if execError != nil {
		return false, nil
	}

	return true, nil

}

func (c Db) removeInfoImage() (bool, error) {
	/*
		Удаляет данные из БД
	*/

	db, dateBaseError := dbConnection()
	if dateBaseError != nil {
		return false, nil
	}

	var deleteReq = fmt.Sprintf("DELETE FROM " + c.TableName + " WHERE Id = '" + c.ImageS.Id + "'")
	_, execError := db.Exec(deleteReq)

	if execError != nil {
		return false, nil
	}

	return true, nil

}

func (c Db) fetchInfo() ([]any, error) {
	/*
		Выкачивает всю инфу из БД
	*/

	db, dateBaseError := dbConnection()
	if dateBaseError != nil {
		panic(dateBaseError)
	}
	var results []any
	record, queryError := db.Query("SELECT * FROM " + c.TableName)

	if queryError != nil {
		return nil, queryError
	}

	defer func(record *sql.Rows) {
		err := record.Close()
		if err != nil {
			panic(err)
		}

	}(record)

	if c.FetchInfo == "posts" {
		for record.Next() {
			var Id string
			var Title string
			var Content string
			var Date string
			scanError := record.Scan(&Id, &Title, &Content, &Date)

			if scanError != nil {
				return results, scanError
			}

			var user = PostData{Id: Id, Title: Title, Content: Content, Date: Date}

			results = append(results, user)
		}

		return results, nil

	}
	return nil, nil

}

func (c Db) getImageById() ([]byte, error) {
	/*
		Выкачивает всю инфу из БД
	*/

	db, dateBaseError := dbConnection()
	if dateBaseError != nil {
		panic(dateBaseError)
	}
	record, queryError := db.Query("SELECT Image FROM " + c.TableName + " WHERE Id= '" + c.ImageS.Id + "'")

	if queryError != nil {
		return nil, queryError
	}

	defer func(record *sql.Rows) {
		err := record.Close()
		if err != nil {
			panic(err)
		}

	}(record)

	if c.FetchInfo == "post_images" {
		for record.Next() {
			var Image []byte
			scanError := record.Scan(&Image)

			return Image, scanError
		}

	}
	return nil, nil

}

func (c Db) getPostById() ([]byte, error) {
	/*
		Выкачивает всю инфу из БД
	*/

	db, dateBaseError := dbConnection()
	if dateBaseError != nil {
		panic(dateBaseError)
	}
	record, queryError := db.Query("SELECT * FROM " + c.TableName + " WHERE Id= '" + c.PostD.Id + "'")

	if queryError != nil {
		return nil, queryError
	}

	defer func(record *sql.Rows) {
		err := record.Close()
		if err != nil {
			panic(err)
		}

	}(record)

	if c.FetchInfo == "posts" {
		for record.Next() {
			var Id string
			var Title string
			var Content string
			var Date string
			scanError := record.Scan(&Id, &Title, &Content, &Date)

			var data = map[string]any{"title": Title, "content": Content, "id": Id, "date": Date}
			var dataPosts, jsonError = json.MarshalIndent(data, "", "   ")
			if jsonError != nil {
				panic(jsonError)
			}
			return dataPosts, scanError
		}

	}
	return nil, nil

}

func (c Db) createTable() error {
	/*
		Создает таблицу posts
	*/

	db, dateBaseError := dbConnection()
	if dateBaseError != nil {
		panic(dateBaseError)
	}

	var table string

	for tableName, Params := range c.Tables {
		var tableType = "CREATE TABLE IF NOT EXISTS "
		table = tableType + tableName + " (Id CHAR(50) NOT NULL,\n"
		for name, param := range Params {
			table = table + name + " " + param + ",\n"
		}
		table = table[:len(table)-2]
		table = table + ");"
	}
	fmt.Println("СОЗДАНИЕ ТАБЛИЦЫ: ", table)
	query, prepareError := db.Prepare(table)
	if prepareError != nil {
		return prepareError
	}
	_, execError := query.Exec()

	if execError != nil {
		return execError
	}

	return nil
}

func InitImagesDbAdmin() {
	var imageTable = map[string]string{"Image": "BLOB"}
	var tables = map[string]map[string]string{"post_images": imageTable}
	var db = Db{DbName: "eska", TableName: "post_images", FetchInfo: "post_images", Tables: tables}
	Err := db.createTable()

	if Err != nil {
		panic(Err)
	}

}

func InitPostsDbAdmin() {
	var postsTable = map[string]string{"Title": "TEXT", "Content": "TEXT", "Date": "TEXT"}
	var tables2 = map[string]map[string]string{"posts": postsTable}
	var db = Db{DbName: "eska", TableName: "posts", FetchInfo: "posts", Tables: tables2}
	Err := db.createTable()

	if Err != nil {
		panic(Err)
	}

}
