// Создание базы данных и спарвочников

package main

// Пакеты
import (
  r "github.com/dancannon/gorethink"
	"fmt"
    "net/http"
    "time"
    "log"
    "strings"
    "encoding/json"
	"io/ioutil"
    // "html/template"
)


// ***************************************************************************************
// Создание таблиц на основе массива наименований
// ***************************************************************************************
func CreateTabs(db string, tables []string) {
     for _, tab := range(tables) {
            CreateTable(db, tab) 
     }
}


// ***************************************************************************************
// Cоздание базы данных
// ***************************************************************************************
func CreateDb(Name string) {
	err:=r.DBCreate(Name).Exec(sessionArray[0])
	
	if err!=nil{
	   fmt.Println(err)
	   return
	}
	log.Println("Cоздана база данных", Name)
}

// ***************************************************************************************
// Cоздание базы данных
// ***************************************************************************************
func DeleteDb(Name string) {
	err:=r.DBDrop(Name).Exec(sessionArray[0])
	
	if err!=nil{
	   fmt.Println(err)
	   return
	}
	
	log.Println("База удалена", Name)
}

// ***************************************************************************************
// Создание таблицы
// ***************************************************************************************
func CreateTable(Db, Name string) {
	 r.DB(Db).TableCreate(Name).Exec(sessionArray[0])
	 log.Println("Cоздана таблица", Name)
}

// ***************************************************************************************
// Добавление в таблицу
// ***************************************************************************************
func InsTable(Db,Tab string, Rec Mst){
	i:= r.InsertOpts{Durability: "soft"}                  // Insrt Optc
	r.DB(Db).Table(Tab).Insert(Rec, i).Exec(sessionArray[0])
}

// ***************************************************************************************
// Insert добавление данных в таблицу
// ***************************************************************************************
func InsTabs(Db, Tab , Recs string){
	 i:= r.InsertOpts{Durability: "soft"}                  // Insrt Optc
     // Rc:=SJ(Recs)
     // err:= r.DB(Db).Table(Tab).Insert(r.JSON(Rc)).Exec(sessionArray[0])
     err:= r.DB(Db).Table(Tab).Insert(SJ(Recs), i).Exec(sessionArray[0])
     if err!=nil{
        fmt.Println("Error:", err)
     }

     log.Println("Данные добавлены в таблицу", Tab)
}


/********************************************************************************************************************************
 *
 * Проверка наличия документа в таблице
 * Если нет документа = 0
 *
 *********************************************************************************************************************************/
func Sys_CountDocument(D,T string) int {
	var response int
	res, err := r.DB(D).Table(T).Count().Run(sessionArray[0])
	defer res.Close()

	// Error
	if err != nil {
		log.Println(err)
		panic("No document ...")
	}

	err = res.One(&response)

	// Error
	if err != nil {
		fmt.Println("Document # Absent ")
		log.Println(err)
		panic("No document")
	}
	return response
}


/********************************************************************************************************************************
 *
 * Проверка наличия документа в таблице
 * Если нет документа возвращаем сообщение
 * Sys_check_doc ("test","Docmove","id")
 *
 *********************************************************************************************************************************/
func Sys_check_doc(DataBase, Tablename, IdDoc string) int {
	var response int
	res, err := r.DB(DataBase).Table(Tablename).Get(IdDoc).Count().Run(sessionArray[0])

	if err != nil {
		log.Println(err)
		panic("No document ...")
	}

	err = res.One(&response)
	defer res.Close()

	// Error
	if err != nil {
		// fmt.Println("Document # Absent ")
		// log.Println(err)
		panic("No document")
		return -1
	}
	return response
}

/********************************************************************************************************************************
 *
 * Возврат стрингового поля
 * Sys_returnstr(DataBase, Tablename, Field, IdDoc string) string
 *
 *********************************************************************************************************************************/
func Sys_returnstr(DataBase, Tablename, Field, IdDoc string) string {

	var response string
	res, err := r.DB(DataBase).Table(Tablename).Get(IdDoc).Field(Field).Run(sessionArray[0])

	if err != nil {
		panic("No document ...")
	}
	err = res.One(&response)
	if err != nil {
		return ""
	}

	defer res.Close()
	return response
}

/********************************************************************************************************************************
 *
 * Update Document
 * Обновление документа в таблице
 *
 *********************************************************************************************************************************/
func DocUpdateHandler(w http.ResponseWriter, r *http.Request) {
	remPartOfURL := r.URL.Path[len("/upd/"):]
	DocUpdate(remPartOfURL)
	fmt.Fprintf(w, "<h1 style='color:#FFF;'>Данные были обновлены</h1> \n %s!", remPartOfURL)
}

/*********************************************************************************************************************************
 *
 * Document Update
 * Обновление документа с записью в лог таблицу
 *
 *********************************************************************************************************************************/
func DocUpdate(IdDoc string) {
	 r.DB(DBN).Table(TBN).Get(IdDoc).Update(Mst{"OperationName": "Документ обнавлен"}).Run(sessionArray[0])
	 
}

/********************************************************************************************************************************
  Удаление документа
*********************************************************************************************************************************/
func shoutDelete(w http.ResponseWriter, r *http.Request) {
	remPartOfURL := r.URL.Path[len("/del/"):]
	// stJson:=remPartOfURL)
	DeleteDoc(remPartOfURL)
	fmt.Printf("Documet is Number # %s deleted ! ", remPartOfURL)                    // Вывод в окно сервиса
	fmt.Fprintf(w, "<h1>Data is deleted</h1> \n %s!", strings.ToUpper(remPartOfURL)) // Вывод в окно браузера
}


/********************************************************************************************************************************
 *
 * Delete Document form table
 *
 *********************************************************************************************************************************/
func DeleteDoc(NumberDocKey string) {
	rd := r.DeleteOpts{Durability: "soft", ReturnChanges: false}
	r.DB(DBN).Table(TBN).Get(NumberDocKey).Delete(rd).RunWrite(sessionArray[0])
	fmt.Println("Document is Deleted... ")
}

/********************************************************************************************************************************
 *
 *  Удаление документов в таблице
 *  таблица и поле Mx
 *  http://127.0.0.1:5555/api/idx/
 *  Всегда возвращает структуру - обязательно необходимо использовать структуру при возврате даже если
 *  возвращается одно поле (например Int)
 *
 *********************************************************************************************************************************/
func DeleteDocuments(NameTable string) {
	rd := r.DeleteOpts{Durability: "soft", ReturnChanges: false}
	r.DB(DBN).Table(NameTable).Delete(rd).RunWrite(sessionArray[0])

}

/********************************************************************************************************************************
 *
 *
 * Добавление нового документа в таблицу в формате JSON используя метод BODY
 *
 * Если m []*Mst - для вставки нескольких документов в запросе
 *      m Mst    - втсавка предполагает только один документ
 *
 *
 * Примечание : все файлы должны быть в кодировке UTF-8 без BOM
 * Пример вставки одного документа
 *     curl -X PUT  http://10.0.3.24:5555/api/docadd/  -d {\"Number\":\"documents\",\"Samples\":\"R-020030031\"}
 *
 * Пример строки для отправки тестового документа
 *     curl -X POST http://10.0.3.24:5555/api/docadd/ -T "test.json" -H "Content-Type: application/json; charset=utf-8;"
 *
 * Пример JSON файла : test.json (именно такой структуры)
 *	   [
 *	     {"ID":"Привмер 10000"},
 *	     {"ID":"sssddd2 200000"},
 *	     {"ID":"sssddd333 30000"},
 *	     {"Date":"2014-02-01T00:23:34"}
 *	   ]
 *
 *********************************************************************************************************************************/
func AddDocumenttoDatabase(w http.ResponseWriter, req *http.Request) {

	// Формат времени RFC3339Nano
	var CurTm = time.Now().Format(time.RFC3339Nano)

	// JSON формат времени
	// CurT, _ := time.Now().MarshalJSON()
	// Для вставки в таблицу необходимо привести к символьному значению  string(CurT)
	reads, errt := ioutil.ReadAll(req.Body)
	if errt != nil {return}
	defer req.Body.Close()

	DocumentStr := string(reads)

	// Формирование для нескольких документов в Json файле документа
	var m []*Mst // Автоматически подходит для всех форматов Json

	// Формирование для одного документа
	// {}
	//m  := make(map[string]interface{})
	errj := json.Unmarshal([]byte(reads), &m)
	met := req.Method


	// Обработка ошибок
	if errj != nil {
	   log.Println(errj)
	   panic("\n \n Problem With Load Document! \n \n")
	}

	// Alianceid:=m["ID_OWNER_CORP"]
	// Maxseq:=MaxID(m["ID_OWNER_CORP"])

	// Добавление документа c дополнительными системными полями
	rk, _ := r.DB("test").Table("temp").Insert(r.Expr(m).Merge(Mst{"HDF_TIME": CurTm}).Merge(Mst{"HDF_SEQ": 2}).Merge(Mst{"Method": met})).Run(sessionArray[0])
	defer rk.Close()

	fmt.Printf("Documet  # %s Был вставлен ! метод %s \n ", DocumentStr, met)

	// fmt.Fprintf("%s", string(reads))
	w.WriteHeader(200)
}


// ***************************************************************************************
// Record to log table
// ***************************************************************************************
func Logos(Description, Warn string){
	 Io:= r.InsertOpts{Durability: "soft"}                  // Insrt Optc
	 Tm:= time.Now().Format("2006-02-01 14:05:00")
     Dt:= Mst{"Description": Description, "Warn": Warn, "Time": Tm} 
     er:= r.DB("System").Table("Log").Insert(Dt, Io).Exec(sessionArray[0])  
     if er!=nil {
     	log.Println("Err insert to log table.", er)
     }
}


// Универсальный способ добавления
func db_Add(DbName, TabName string, w http.ResponseWriter, req *http.Request) {

    // Формирование для нескольких документов в Json файле документа
    // m  := make(map[string]interface{})
	var m []*Mst // Автоматически подходит для всех форматов Json
    i:= r.InsertOpts{Durability: "soft"}       
    u:= req.Header.Get("Uid")                  // Получение секретного ключа
    k:= req.Header.Get("Uid")
	t:= time.Now().Format(time.RFC3339Nano)    // Формат времени RFC3339Nano
	l:= req.Method                             // Method 

    // Тоько пост метод
    if l!="POST" {
       w.Header().Set("Fj-Error",  "Only Post method.")	
       return
    }

	// JSON формат времени
	// CurT, _ := time.Now().MarshalJSON()
	// Для вставки в таблицу необходимо привести к символьному значению  string(CurT)
	reads, errt := ioutil.ReadAll(req.Body)
	if errt != nil {
	   return
	}
	defer req.Body.Close()

    // Нужен тоько для  отбражения получаемых документов
    // в случае контроля данных
	// DocumentStr := string(reads)
	
	// Формирование для одного документа
	errj := json.Unmarshal([]byte(reads), &m)

	// Обработка ошибок
	if errj != nil {
	   log.Println(errj)
	   panic("\n \n Problem With Load Document! \n \n")
	}

    go func (){
		// Добавление документа c дополнительными системными полями
		er := r.DB(DbName).Table(TabName).Insert(r.Expr(m).Merge(Mst{"Timeinsert": t, "Seq": u, "Method": l, "Uid": u, "Ir":k}),i).Exec(sessionArray[0])
		    
		// Обработка ошибок
		if er != nil {
		   es:=	er.Error()
		   w.Header().Set("Error : ", es)	
		   log.Println("ERR LOAD TO DATABASE : " + es)
		   return
	    }else{
           // Возврат в head информации
	       w.Header().Set("Fj",  "Insert is OK.")
	    }
	}()
	
    w.Header().Set("Uid", u)
	// fmt.Printf("Documet  # %s Был вставлен ! метод %s \n ", DocumentStr, met)
	// fmt.Fprintf("%s", string(reads))
	// w.WriteHeader(200)
}
