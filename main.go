package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render
var db *mgo.Database

const (
	hostname   string = "localhost:27017"
	dbName     string = "demo_todo"
	collection string = "todo"
	port       string = ":9000"
)

type (
	todoModel struct {
		ID        bson.ObjectId `bson:"id,omitempty"`
		Title     string        `bson:"title"`
		Completed bool          `bson:"completed"`
		CreatedAt time.Time     `bson:"createdAt"`
	}

	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed string    `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func todoHandlers() http.Handler {
	routerGroup := chi.NewRouter()
	routerGroup.Group(func(r chi.Router) {
		r.Get("/", fetchTodos)
		r.Post("/", createTodos)
		r.Put("{id}", updateTodo)
		r.Delete("{id}", deleteTodo)
	})
	return routerGroup
}

func init() {
	rnd = renderer.New()
	sess, err := mgo.Dial(hostname)
	checkErr(err)
	sess.SetMode(mgo.Monotonic, true)
	db = sess.DB(dbName)
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", homeHandler)
	router.Mount("/todo", todoHandlers)

	srv := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Println("Listening on port", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen:%s\n", err)
		}
	}()
}
