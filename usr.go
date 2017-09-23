package main

// Пакеты
import (
   "net/http"
 
)

// 
// Добавление организации
// http://localhost:5555/vacation/
// 
func UsrAdd(w http.ResponseWriter, req *http.Request) {
     db_Add("Sys","User", w,req)
}


func UsrClear(w http.ResponseWriter, req *http.Request) {
     
}

