package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type Tweets struct {
	Id    int
	Tweet string
}

var DbConnection *sql.DB

func main() {
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/tweet_delete/", deleteTweet)
	http.HandleFunc("/tweet/", getPostTweet)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	DbConnection, _ := sql.Open("sqlite3", "./example.sql")
	defer DbConnection.Close()
	cmd := `SELECT * FROM tweets`
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	var body []Tweets
	for rows.Next() {
		var b Tweets
		err := rows.Scan(&b.Id, &b.Tweet)
		if err != nil {
			log.Fatalln(err)
		}
		body = append(body, b)
	}
	t, err := template.ParseFiles("views/index.html")
	if err != nil {
		log.Fatalln(err)
	}
	t.Execute(w, body)
}

func getPostTweet(w http.ResponseWriter, r *http.Request) {
	DbConnection, _ := sql.Open("sqlite3", "./example.sql")
	defer DbConnection.Close()
	v := r.FormValue("tweet")
	cmd := `INSERT INTO tweets(tweet)VALUES(?)`
	DbConnection.Exec(cmd, v)
	http.Redirect(w, r, "/", http.StatusFound)
}

func deleteTweet(w http.ResponseWriter, r *http.Request) {
	DbConnection, _ := sql.Open("sqlite3", "./example.sql")
	defer DbConnection.Close()
	cmd := "DELETE FROM Tweets WHERE id = ?"
	i := r.FormValue("tweet_delete")
	var I int
	I, _ = strconv.Atoi(i)
	DbConnection.Exec(cmd, I)
	http.Redirect(w, r, "/", http.StatusFound)
}
