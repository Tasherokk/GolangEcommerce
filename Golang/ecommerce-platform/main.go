package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Tasherokk/GolangProject.git/controllers"
	"github.com/Tasherokk/GolangProject.git/database"
	"github.com/Tasherokk/GolangProject.git/middleware"
	"github.com/Tasherokk/GolangProject.git/routes"
	"github.com/Tasherokk/GolangProject.git/test"
	"github.com/Tasherokk/GolangProject.git/validators"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {


	initLog()
	database.InitDB()
    defer func() {
		if err := database.DB.Close(); err != nil {
			log.Fatalf("Failed to close the database: %v", err)
		}
	}()
	database.InitRedis()
	defer database.RDB.Close()
	validators.Init();


	// startTest()


	startServer()

	

}





func initLog() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}


func startServer() {
	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)
	routes.RegisterRoutes(router)
	router.HandleFunc("/users/{id}", controllers.GetUserHandler).Methods("GET")
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {http.NotFound(w, r)})	
	router.Handle("/metrics", promhttp.Handler())
	

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router))
}



func startTest() {
	test.Test()
	test.TestDb()
	test.TestDb2()
}





