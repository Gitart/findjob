

package main

// Пакеты
import (
 r "github.com/dancannon/gorethink"
   "fmt"
   "net/http"
   "time"
   "path"
   "html/template"
)




// 
// Добавление организации
// http://localhost:5555/vacation/
// 
func Vacation(w http.ResponseWriter, req *http.Request) {
      
    // Путь к форме
	fp        := path.Join("html", "company.html")
	tmpl, err := template.ParseFiles(fp)
	
	if err != nil {
	   http.Error(w, err.Error(), http.StatusInternalServerError)
	   return
	}


   var ZpRec Mst

    // Zps=Mst{"Rec":"Нет записи"}
    // считывание с формы и сохранение в базу данных
    if req.Method=="POST" {

            Tm          := time.Now().Format(time.RFC3339)
            Tm           = Ctm(2)
       	    Organ       := req.FormValue("organ")       
       	    Contact     := req.FormValue("contact")       
	        Adress      := req.FormValue("adress") 
	        Telphone    := req.FormValue("telphone") 
	        Site        := req.FormValue("site") 
	        Mobile      := req.FormValue("mobtel") 
	        Id          := RandomHumanFriendlyString(10)

             Errs:="Ok"

             if Organ==""||Contact=="" {
                Errs="Error data input"
              }

	       	// Формировнаие формы для встаавки
	        ZpRec = Mst{
	        	         "Organ"  : Organ,
	        	         "Contact": Contact, 
	                     "Adress" : Adress, 
	                     "Telph"  : Telphone,
	                     "Site"   : Site,
	                     "Mobile" : Mobile,
	                     "Time"   : Tm, 
	                     "ID"     : Id,
	                     "Tim"    : time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
	                     "Error"  : Errs,
	                     }


            // Проверка на заполнение основных полей
            if Errs=="Ok" {

			        // Запись в базу
			        er := r.DB("Fj").Table("Organ").Insert(ZpRec).Exec(sessionArray[0])
			              if er != nil {
				             fmt.Println(er)
				             return
			                }      
		                    fmt.Println("Да сохранились")
            }else{
            	            fmt.Println("NO Save Record!")
            }
    }
     
tmpl.Execute(w, ZpRec)



/*
	StructureID := req.FormValue("contact")       
	Cabc        := req.FormValue("adress")   

	
	// Формировнаие формы для встаавки
	ZpRec := Mst{"STRUCTURE": StructureID, "ABCC": Cabc }

	// Запись в базу
	er := r.DB("Barsetka").Table("Log").Insert(ZpRec).Exec(sessionArray[0])
	if er != nil {
		return
	}
	// defer rks.Close()

	book:=Mst{"Inser":"Ok"}
	// Вызов формы из темплейта
	tmpl.Execute(w, book)
*/
}