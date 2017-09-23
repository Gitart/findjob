package  main
import (
  r "github.com/dancannon/gorethink"
    "encoding/json"
	"time"
	"log"
	"strings"
	"fmt"
	"net/http"
	"io/ioutil"
	  mrn "math/rand"
)



//
// Запись в лог файл
// 
func LOG(C, T string) {
     R  := r.InsertOpts{ReturnChanges: false, Durability:"soft"}
 	 tm := time.Now().Format("2006-01-02T15:04:05.000")
	 I  := Mst{"TITLE": T, "DATETIME": tm, "CORPORATIONNAME": C}
	 Q  := r.DB("System").Table("Log").Insert(I,R)
     Rk, err := Q.Run(sessionArray[0])

	 if err != nil {
	    return
	 }

	 defer Rk.Close()
}


//     TITLE       : Преобразование строки в JSON формат
//     AUTHOR      : Savchenko Arthur
//     DATE TIME   : 25-04-2016 11:23
//     USAGE       : t:=SJ(`[{"":"","":""},{"":"","":""},{"":"","":""}.....{"":"","":""}]`) = JSON
//     DESCRIPTION : Analog => r.JSON()
//     NOTE        : Всегда брать в [] - даже если одиночное значение = [{"":""}]
func SJ(m string) []Mst {
	var T, TE []Mst
	err := json.Unmarshal([]byte(m), &T)

	if err != nil {
	   return TE
	}
	return T
}


// Чтение JSON файла из директории
// http.HandleFunc("/tst/readjsonfile/", Test_ReadJsonFile)                  // Read Json file
func Test_ReadJsonFile(w http.ResponseWriter, rs *http.Request) {

	var response []interface{}
	h := `http://195.128.18.66:5555/static/data/Contractors.json`

	res, err := r.HTTP(h).Run(sessionArray[0])

	// Error
	if err != nil {
	   log.Println("No open table for Import ...")
	}

	err = res.All(&response)

	if err != nil {
	   fmt.Fprintf(w, "%s", strings.ToUpper("404"))
	   w.WriteHeader(204)
	} else {
	   data, _ := json.Marshal(response)
	   w.WriteHeader(200)
	   fmt.Fprintf(w, string(data))
	}
}



//   Title       : Текущее время в формате (YYYY-MM-DD HH:MM:SS)
// 	 Date        : 2015-12-14
func CTM() string {
	 return time.Now().Format("2006-01-02 15:04:05")
}


// Title       : Текущее время в формате (YYYY-MM-DD HH:MM:SS)
// Date        : 2015-12-14
func Ctm(Dt int) string {
	 T:=[]string{"02/01/2006", "02.01.2006", "02-01-2006", "15:04:05", "150405", "02.01.2006 15:04:05", "02/01/2006 15:04:05", "02-01-2006 15:04:05", "2006-01-02T15:04:05"}
     D:=T[Dt]
	 return time.Now().Format(D)
}

// Title       : Текущее время в формате UNixtime Nano
// Date        : 2016-04-25 12:21
func CTU() int64 {
	 return time.Now().UnixNano() / 1000000
}



//   Title       : Cоздание файла c расширением
// 	 Date        : 2016-01-05
func WriteLogFile(NameFile, Ext, Text string) {
	var n string

	if NameFile == "" {
		n = "Log_" + time.Now().Format("20060102150405") + "." + Ext              // Формат
	} else {
		n = "log/" + NameFile + time.Now().Format("20060102150405") + "." + Ext   // отправка лога в директорию лог
		n = NameFile + time.Now().Format("20060102150405") + "." + Ext            // отправка лога в текущую директорию
	}

	t   := []byte(Text)
	err := ioutil.WriteFile(n, t, 0644)

	if err != nil {
		return
		//fmt.Println(err)
	}
}


// Title               : Формирование индексов для таблицы
// Date & Time         : 09.02.2016 12:00
// Описание параметров :
//      Создание множества индексов для выбранной таблицы
//      Первый параметр имя базы
//      Второй параметр имя таблицы
//      Все остальные параметры имя индексов....
//      Sys_create_indeхs("DB","TB","IDX1","IDX2","IDX3")
func Sys_create_indeхs(DB, TB string, IDX ...string) {
	for _, t := range IDX {
		//rk ,err:=r.DB(DB).Table(TB).IndexCreate(t).Run(sessionArray[0])
		//defer rk.Close()
		//if err!=nil{return}
		err := r.DB(DB).Table(TB).IndexCreate(t).Exec(sessionArray[0])
		if err != nil {
		   return
		}
	}
}

// Title               : Просмотр JSON 
// Date & Time         : 12.02.2016 12:00
func jsonPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}



/***************************************************************************************************
 *
 * Generate a random, but human-friendly, string of the given length.
 * Should be possible to read out loud and send in an email without problems.
 *
 * Samples : ubipamus
 *
 ***************************************************************************************************/
func RandomHumanFriendlyString(length int) string {
	const (
		vowels     = "aeiouy"                   // email+browsers didn't like "æøå" too much
		consonants = "bcdfghjklmnpqrstvwxz"
	)

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		if i%2 == 0 {
			b[i] = vowels[mrn.Intn(len(vowels))]
		} else {
			b[i] = consonants[mrn.Intn(len(consonants))]
		}
	}
	return string(b)
}

/***************************************************************************************************
 *
 * Generate a random, but cookie-friendly, string of the given length.
 *
 ***************************************************************************************************/
func RandomCookieFriendlyString(length int) string {
	 const allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	 b := make([]byte, length)
	 for i := 0; i < length; i++ {
	 	 b[i] = allowed[mrn.Intn(len(allowed))]
	 }
	 return string(b)
}