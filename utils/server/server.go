package server

import (
	"net/http"
	"os"

	userRouter "movies/internal/user/router"

	mux "github.com/gorilla/mux"
	cors "github.com/rs/cors"
)

func Start() error {

	// cors config
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "http://localhost:5000"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "HEAD"},
	})

	// init Router
	r := mux.NewRouter()

	userRouter := userRouter.NewUserRouter(r)
	userRouter.Handle()

	port := ":" + os.Getenv("PORT")
	handler := c.Handler(r)

	if err := http.ListenAndServe(port, handler); err != nil {
		return err
	}
	return nil
}
