package main

// Пакеты
import (
 r "github.com/dancannon/gorethink"
   "runtime"
   "fmt"
   "net/http"
   // "html/template"
   // "path"
   // "time"
      
)


var (IpPort string = "")

/***************************************************************************************************************************************
 *
 *   Title        : Connection to DB
 *   Initialisation Service
 * 	 Date         : 2015-03-11
 *	 Description  : Initialization DB Connect
 *   Author       : SAVCHENKO ARTHUR 
 *   ☎           : 8-097-5547468
 *
 ****************************************************************************************************************************************/
func Dbini() {

	// Инициализация подключения к базе
	// на той машине где расположен и стартует сервис
	// Для переключения на тестовую машину

	if runtime.GOOS == "windows" {
	   IpPort = "191.101.15.57:28015"                        // Продукшн сервер
	   IpPort = "localhost:28015"                            // Тестовый локальный сервис

        // Connect to Production Server
        production := false
        if production == true { 
           IpPort = "193.111.9.202:28015"
        }

	}  else {
	     IpPort = "loclahost:28015"                       // Локальный ресурс 
	}


	// Сессия подключения
	// MaxActive: 100, IdleTimeout: time.Second * 10, MaxIdle: 10})

	// Для случая когда будет установлен пароль  базы данных
	// ServiceClientKey       = "S0864AA791CE7E7B00RT"
    // AuthKey: DatabaseKey,
    // 

	// session, err := r.Connect(r.ConnectOpts{Address: IpPort, Database: DatabaseName, Password: ServiceClientKey})
	// session, err := r.Connect(r.ConnectOpts{Address: IpPort, Database: DatabaseName, AuthKey: DatabaseKey, Username:"admin", Password: ServiceClientKey})
	session, err := r.Connect(r.ConnectOpts{Address: IpPort, Database: WDB})

    // Закрыть сеесию 
	// defer session.Close()

	// Для случая когда будет установлен пароль для админа + ключ базы данных
	// session, err := r.Connect(r.ConnectOpts{Address: IpPort, Database: DatabaseName, AuthKey: DatabaseKey, Username: "admin", Password: ServiceClientKey})

	// Для случая когда будет установлен пароль для админа
	// session, err := r.Connect(r.ConnectOpts{Address: IpPort, Database: DatabaseName, Username: "admin", Password: ""})

	// Обработка ошибок
	if err != nil {
	   fmt.Println("NO CONNECTION.")
	   return
	}

	// Максимальное количество подключений
	session.SetMaxOpenConns(200)
	session.SetMaxIdleConns(200)
	sessionArray = append(sessionArray, session)
}

// Информационная часть для высветки
func InfoStart(){
	Dbini()

    AdresPort  = AdressService + PortService    // Service work port
	AdresPort  = ":8092"                        // Port Only for anythink IP
	AdresPort  = ":3333"                        // Port Only for anythink IP
	Typeos    := runtime.GOOS                   // Type operation system (LINUX,WINDOWS)
	Versstr   := runtime.Version()              // Version GO
    fmt.Println("Type Os :",  Typeos)
    fmt.Println("Version :",  Versstr)
    fmt.Println("Service : Find Job is Started..")
}

/*******************************************************************************************************************************
 * DATETIME         : 28-07-2015 12:44
 * DESCRIPTION      : Стартовая процедура
 * NOTES            : Запуск сервиса с параметрами
 *******************************************************************************************************************************/
func main() {
  	InfoStart()

    // Routes
	http.HandleFunc("/static/",                StaticPage)  // Статические страницы
	http.HandleFunc("/db/create/",             Start)       // Стартовые операции - создание базы и заполнение спрачоников
	http.HandleFunc("/info/",                  Info)        // Информация о сервисе
	http.HandleFunc("/vacation/",              Vacation)    // Добавление вакансии в форме

    http.HandleFunc("/usr/add/",               UsrAdd)      // Добавление вакансии в форме    
    http.HandleFunc("/db/clear/",              UsrClear)      // Добавление вакансии в форме    

    // Start service for listing
	err:=http.ListenAndServe(AdresPort, nil)                    
	if err!=nil{
		fmt.Println("Conflict Ports.")
		return
	}
}

/********************************************************************************************************************************
 Статические странички
 c установкой разрешений и доступов на операции
 Нужна для организации доступа к CSS и библиотекам
 http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {http.ServeFile(w, r, r.URL.Path[1:])})
 http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {http.FileServer(http.Dir("/static/"))})
 /static/....
 *********************************************************************************************************************************/
func StaticPage(w http.ResponseWriter, r *http.Request) {
	// Allows
	w.Header().Set("Access-Control-Allow-Origin", "*")
	/* Allows
	      if origin := r.Header().Get("Origin"); origin != "" {
		     w.Header().Set("Access-Control-Allow-Origin", origin)
		     w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		     w.Header().Set("Access-Control-Allow-Headers",  "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		     w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
		    }
	*/
	//  File static page
	http.ServeFile(w, r, r.URL.Path[1:])
}









