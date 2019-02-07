package main
 
import (
    "net/http"
    "github.com/gorilla/mux"
    //"fmt"
)
 
var router = mux.NewRouter()

func main() {
    router.HandleFunc("/", LoginPageHandler) // GET (Homepage)
 
    router.HandleFunc("/index", IndexPageHandler) // GET
    router.HandleFunc("/login", LoginHandler).Methods("POST")
 
    router.HandleFunc("/register", RegisterPageHandler).Methods("GET")
    router.HandleFunc("/register", RegisterHandler).Methods("POST")
 
    router.HandleFunc("/logout", LogoutHandler).Methods("POST")

    router.HandleFunc("/home", HomeHandler).Methods("POST")
    
    http.Handle("/", router)
    http.ListenAndServe(":8888", nil) //port
}