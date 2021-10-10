package src

import (
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
	"home-counter/src/config"
	"home-counter/src/handlers"
	"home-counter/src/middlewares"
	"home-counter/src/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// TODO затащить gin
func App() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func (c chan os.Signal) {
		<- c
		models.DB.Close()
		log.Printf("DB connections closed")
		os.Exit(0)
	}(c)

	defer func() {
		if err := recover(); err != nil {
			if err == config.ErrorConfig {
				log.Printf("Config error")
			} else {
				log.Printf("Panic happened. %s. Recovering server", err)
				App()
			}
		}
	}()

	defer func() {
		models.DB.Close()
		log.Printf("DB connections closed")
	}()

	authRouter := httprouter.New()
	authRouter.GET("/auth/callback", handlers.CallbackHandler)
	authRouter.GET("/auth/login", handlers.LoginHandler)
	authRouter.GET("/auth/logout", handlers.LogoutHandler)

	userRouter := mux.NewRouter()
	userRouter.HandleFunc("/user", handlers.UserDataHandler).Methods("GET")
	userRouter.HandleFunc("/user/tariffs", handlers.UserTariffsHandler).Methods("POST")
	userRouter.HandleFunc("/user/meters", handlers.CreateOrUpdateUserMetersHandler).Methods("POST")
	userRouter.HandleFunc("/user/count-meters", handlers.UserMeterCountHandler).Methods("GET")
	userHandler := middlewares.AuthMiddleware(userRouter)

	coreRouter := mux.NewRouter()
	coreRouter.Handle("/user", userHandler)
	coreRouter.Handle("/user/tariffs", userHandler)
	coreRouter.Handle("/user/meters", userHandler)
	coreRouter.Handle("/user/count-meters", userHandler)
	coreRouter.Handle("/auth/login", authRouter)
	coreRouter.Handle("/auth/logout", authRouter)
	coreRouter.Handle("/auth/callback", authRouter)
	coreRouter.HandleFunc("/", handlers.Counter).Methods("GET")

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(coreRouter)

	s := &http.Server{
		Addr:           ":" + config.Config.Port,
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Server is starting %s", config.Config.AppUrl)
	log.Fatal(s.ListenAndServe())
}
