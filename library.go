package main

import (
	"fmt"

	// mysql connector
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	User     = "root"
	Password = "hyc021800590"
	DBName   = "ass3"
)

type Library struct {
	db *sqlx.DB
}

func (lib *Library) ConnectDB() {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	if err != nil {
		panic(err)
	}
	lib.db = db
	//fmt.Println("cn!")
}

// CreateTables created the tables in MySQL
func mustExecute(db *sqlx.DB, SQLs []string) {
	for _, s := range SQLs {
		_, err := db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
}

func (lib *Library) init() {
	mustExecute(lib.db,[]string{"delete from ADMINISTRATORS","delete from DELETE_RECORDS","delete from BORROW_RECORDS","delete from BOOKS","delete from STUDENTS"});
	var s1,s2,s3,s4 string
	s1 = "INSERT INTO BOOKS(ISBN,author,title,TOTAL_NUM,AVAILABLE_NUM) VALUES(\"112\",\"a01\",\"b01\",1,1),(\"113\",\"a02\",\"b01\",2,0),(\"114\",\"a01\",\"b02\",4,1),(\"115\",\"a03\",\"b03\",4,3)"
	s2 = "INSERT INTO STUDENTS(SID,SNAME,PASSWORD) VALUES(\"s01\",\"ab\",\"abab\"),(\"s02\",\"ac\",\"acac\"),(\"s03\",\"ad\",\"adad\")"
	s3 = "INSERT INTO BORROW_RECORDS(SID,ISBN,BORROW_DATE,DUE_DATE,ExtendableTimes) VALUES(\"s01\",\"113\",'2020-04-01','2020-05-01',3),(\"s01\",\"114\",'2020-04-02','2020-05-02',3),(\"s01\",\"115\",'2020-01-03','2020-05-01',0),(\"s02\",\"114\",'2020-04-01','2020-05-01',3),(\"s02\",\"114\",'2020-04-01','2020-05-01',3),(\"s03\",\"113\",'2020-04-28','2020-05-28',3)"
	s4 = "INSERT INTO ADMINISTRATORS(AID,PASSWORD) VALUES(\"a01\",\"haha\")"
	mustExecute(lib.db,[]string{s1,s2,s3,s4});
}

func (lib *Library) CreateTables() error {
	var s1,s2,s3,s4,s5 string
	s1 = "CREATE TABLE IF NOT EXISTS BOOKS (ISBN char(20) NOT NULL,author CHAR(32) NOT NULL,title CHAR(32) NOT NULL,TOTAL_NUM SMALLINT NOT NULL,AVAILABLE_NUM SMALLINT NOT NULL,PRIMARY KEY(ISBN));"
	s2 = "CREATE TABLE IF NOT EXISTS STUDENTS (SID CHAR(20) NOT NULL,SNAME CHAR(8) NOT NULL,PASSWORD CHAR(30) NOT NULL,PRIMARY KEY(SID));"
	s3 = "CREATE TABLE IF NOT EXISTS BORROW_RECORDS (SID CHAR(20) NOT NULL,ISBN char(20) NOT NULL,BORROW_DATE DATE NOT NULL,DUE_DATE DATE NOT NULL,RETURN_DATE DATE,ExtendableTimes SMALLINT NOT NULL,FOREIGN KEY(ISBN) REFERENCES BOOKS(ISBN),FOREIGN KEY(SID) REFERENCES STUDENTS(SID));"
	s4 = "CREATE TABLE IF NOT EXISTS DELETE_RECORDS (ISBN char(20) NOT NULL,EXPLAINATION CHAR(30) NOT NULL,DELETE_DATE DATE NOT NULL);"
	s5 = "CREATE TABLE IF NOT EXISTS ADMINISTRATORS (AID CHAR(20) NOT NULL,PASSWORD CHAR(30) NOT NULL,PRIMARY KEY(AID));"
	mustExecute(lib.db,[]string{s1,s2,s3,s4,s5});
        return nil
}

// AddBook add a book into the library
func (lib *Library) AddBook(title, author, ISBN string) (int,error) {
	var s1,s2,s3 string
	var p int
	rows, err := lib.db.Query("SELECT * FROM BOOKS WHERE ISBN =\""+ISBN+"\"")
	if err != nil {
		panic(err)
	}
 	if rows.Next() {
		s2 = "UPDATE BOOKS SET TOTAL_NUM=TOTAL_NUM+1 WHERE ISBN = \""+ISBN+"\""
		s3 = "UPDATE BOOKS SET AVAILABLE_NUM=AVAILABLE_NUM+1 WHERE ISBN = \""+ISBN+"\""
		mustExecute(lib.db,[]string{s2,s3})
		p=0
		fmt.Println("The Total_num of this book increased by 1.");
		
	} else {
		s1 = "INSERT INTO BOOKS(ISBN,author,title,TOTAL_NUM,AVAILABLE_NUM) VALUES(\""+ISBN+"\",\""+author+"\",\""+title+"\",1,1)"
		mustExecute(lib.db,[]string{s1})
		p=1
		fmt.Println("The new book has been successfully added to the library.");
		
	}
	return p,nil
}

func (lib *Library) DeleteBook(ISBN,EXPLAINATION string,pattern int) (int,error) {
	var s1,s2 string
	var p int
	rows, err := lib.db.Query("SELECT TOTAL_NUM,AVAILABLE_NUM FROM BOOKS WHERE ISBN =\""+ISBN+"\"")
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var tnum,num int
		num = 0
		rows.Scan(&tnum,&num)
		if pattern==0 {
			if tnum==num {
				fmt.Println("Error : Books with this ISBN have not been lent.")
				p=0
			} else {
				s1 = "UPDATE BOOKS SET TOTAL_NUM = TOTAL_NUM -1 , AVAILABLE_NUM = AVAILABLE_NUM - 1 WHERE ISBN =\""+ISBN+"\""
				s2 = "INSERT DELETE_RECORDS(ISBN,EXPLAINATION,DELETE_DATE) VALUES(\""+ISBN+"\",\""+EXPLAINATION+"\",CURRENT_DATE())"
				mustExecute(lib.db,[]string{s1,s2});
				p=1
				fmt.Println("Successfully delete the book.")
			}
		} else {
			if num==0 {
				fmt.Println("Error : Books with this ISBN have all been lent.")
				p=2
			} else {
				s1 = "UPDATE BOOKS SET TOTAL_NUM = TOTAL_NUM -1  WHERE ISBN =\""+ISBN+"\""
				s2 = "INSERT INTO DELETE_RECORDS(ISBN,EXPLAINATION,DELETE_DATE) VALUES(\""+ISBN+"\",\""+EXPLAINATION+"\",CURRENT_DATE())"
				mustExecute(lib.db,[]string{s1,s2});
				fmt.Println("Successfully delete the book.")
				p=3
			}
		}
	} else {
		fmt.Println("Error : There is no such book in the library.")
		p=4
	}
	return p,nil
}

func (lib *Library) AddStudent(name,sid,PASSWORD string) (int,error) {
	var s1 string
	var p int
	rows, err := lib.db.Query("SELECT * FROM STUDENTS WHERE SID =\""+sid+"\"")
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		fmt.Println("Error : Student account of this SID already exists.")
		p=0
	} else {
		s1 = "INSERT INTO STUDENTS(SID,SNAME,PASSWORD) VALUES(\""+sid+"\",\""+name+"\",\""+PASSWORD+"\")"
		fmt.Println(s1)
		mustExecute(lib.db,[]string{s1})
		fmt.Println("Student account added successfully.")
		p=1
	}
	return p,nil
}

func (lib *Library) QueryBook(value string,pattern string) (string,error) {
	var ISBN,author,title string
	var TOTAL_NUM,AVAILABLE_NUM int
	var s1 string
	var p string
	rows, err := lib.db.Query("SELECT ISBN,author,title,TOTAL_NUM,AVAILABLE_NUM FROM BOOKS WHERE "+pattern+"=\""+value+"\"")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var cnt int = 0
	for rows.Next() {
		err := rows.Scan(&ISBN,&author, &title,&TOTAL_NUM,&AVAILABLE_NUM);
		if err != nil {
			panic(err)
		}
		cnt++
		s1 = fmt.Sprintf("cnt:%d ISBN=%s author=%s title=%s TOTAL_NUM=%d AVAILABLE_NUM=%d",cnt,ISBN,author,title,TOTAL_NUM,AVAILABLE_NUM)
		fmt.Println(s1);
		p = p + " " +  s1
	}
	if cnt==0 {
		fmt.Println("There is no such book in the library.");
	}
	return p,nil
}

func (lib *Library) BorrowBook(sid,ISBN string) (int,error) {
	var s1,s2 string
	var p int
	rows, err := lib.db.Query("SELECT * FROM STUDENTS WHERE SID =\""+sid+"\"")
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		rows, err := lib.db.Query("SELECT * FROM BORROW_RECORDS WHERE SID =\""+sid+"\" AND ISNULL(RETURN_DATE) AND DUE_DATE < CURRENT_DATE()")
		if err != nil {
			panic(err)
		}
		var cnt int = 0
		for rows.Next() {
			cnt++	
		}
		if cnt>=3 {
			fmt.Println("Your account is suspended because you have more than 3 overdue books.");
			p=0
		} else {
			rows, err := lib.db.Query("SELECT AVAILABLE_NUM FROM BOOKS WHERE ISBN =\""+ISBN+"\"")
			if err != nil {
				panic(err)
			}
			if rows.Next() {
				var num int
				err := rows.Scan(&num)
				if err != nil {
					panic(err)
				}
				if num <= 0 {
					fmt.Println("This book is not available for borrowing at present.");
					p=1
				} else {
					s1 = "UPDATE BOOKS SET AVAILABLE_NUM = AVAILABLE_NUM-1 WHERE ISBN = \""+ISBN+"\""
					s2 = "INSERT INTO BORROW_RECORDS(SID,ISBN,BORROW_DATE,DUE_DATE,ExtendableTimes) VALUES( \""+sid+"\",\""+ISBN+"\",CURRENT_DATE,date_add(CURRENT_DATE(),interval 30 day),3)"
					mustExecute(lib.db,[]string{s1,s2});
					fmt.Println("Successfully borrow this book.")
					p=2
				}
			} else {
				fmt.Println("Error : There is no such book in the library.");
				p=3
			}
		}
	} else {
		fmt.Println("No student account with this ID.")
		p=4
	}
	return p,nil
}

func (lib *Library) QueryHistory(sid string) (string,error) {
	var p,s1 string
	rows, err := lib.db.Query("SELECT * FROM BORROW_RECORDS WHERE SID =\""+sid+"\"ORDER BY BORROW_DATE")
	if err != nil {
		panic(err)
	}
	var cnt int = 0
	for rows.Next() {
		var SID,ISBN,BORROW_DATE,DUE_DATE,RETURN_DATE string
		var ExtendableTimes int
		cnt++
		rows.Scan(&SID,&ISBN,&BORROW_DATE,&DUE_DATE,&RETURN_DATE,&ExtendableTimes)
		if RETURN_DATE=="" {
			RETURN_DATE = "NULL"		
		}
		s1 = fmt.Sprintf("cnt:%d SID=%s ISBN=%s BORROW_DATE=%s DUE_DATE=%s RETURN_DATE=%s",cnt,SID,ISBN,BORROW_DATE,DUE_DATE,RETURN_DATE)
		fmt.Println(s1)
		p = p + " " + s1
	}
	if cnt == 0 {
		fmt.Println("No borrowing record with this ID.")
	}
	return p,nil
}

func (lib *Library) QueryDueDate(sid,ISBN string) (string,error) {
	var p string
	rows, err := lib.db.Query("SELECT DUE_DATE FROM BORROW_RECORDS WHERE SID =\""+sid+"\" AND ISBN =\""+ISBN+"\" AND ISNULL(RETURN_DATE) ORDER BY DUE_DATE")
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var DUE_DATE string
		rows.Scan(&DUE_DATE)
		fmt.Println("The deadline of this book is",DUE_DATE)
		p=DUE_DATE
	} else {
		fmt.Println("No record of borrowing this book with this ID.");
		p="error"
	}
	return p,nil
}

func (lib *Library) ExtendDueDate(sid,ISBN string) (int,error) {
	var s1 string
	var p int
	rows, err := lib.db.Query("SELECT ExtendableTimes FROM BORROW_RECORDS WHERE SID =\""+sid+"\" AND ISBN =\""+ISBN+"\" AND ISNULL(RETURN_DATE) ORDER BY ExtendableTimes DESC,DUE_DATE ASC")
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var ExtendableTimes int
		rows.Scan(&ExtendableTimes)
		if ExtendableTimes==0 {
			fmt.Println("The duedate of this book cannot be extended anymore.");
			p=0
		} else {
			s1 = "UPDATE BORROW_RECORDS SET ExtendableTimes=ExtendableTimes - 1 , DUE_DATE=date_add(DUE_DATE,interval 30 day) WHERE SID =\""+sid+"\" AND ISBN =\""+ISBN+"\" ORDER BY ExtendableTimes DESC,DUE_DATE ASC limit 1"
			mustExecute(lib.db,[]string{s1})
			fmt.Println("Successfully extended due date.")
			p=1
		}
	}
	return p,nil
}

func (lib *Library) QueryOverDue(sid string) (string,error) {
	var p,s1 string
	var cnt int = 0
	rows, err := lib.db.Query("SELECT ISBN,DUE_DATE,ExtendableTimes FROM BORROW_RECORDS WHERE sid =\""+sid+"\" AND DUE_DATE<CURRENT_DATE() AND ISNULL(RETURN_DATE)")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		cnt++
		var ISBN,DUE_DATE string
		var ExtendableTimes int
		rows.Scan(&ISBN,&DUE_DATE,&ExtendableTimes)
		s1 = fmt.Sprintf("cnt:%d ISBN=%s DUE_DATE=%s ExtendableTimes=%d",cnt,ISBN,DUE_DATE,ExtendableTimes)
		fmt.Println(s1)
		p = p + " " + s1
	}
	if cnt==0 {
		fmt.Println("There are no overdue books with this ID.")
	}
	return p,nil
}

func (lib *Library) ReturnBook(sid,ISBN string) (int,error) {
	var s1 string
	var p int
	rows, err := lib.db.Query("SELECT * FROM BORROW_RECORDS WHERE SID =\""+sid+"\" AND ISBN =\""+ISBN+"\" AND ISNULL(RETURN_DATE)")
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		s1 = "UPDATE BORROW_RECORDS SET RETURN_DATE = CURRENT_TIME() WHERE SID =\""+sid+"\" AND ISBN =\""+ISBN+"\" AND ISNULL(RETURN_DATE) ORDER BY DUE_DATE limit 1"
		mustExecute(lib.db,[]string{s1});
		fmt.Println("Successfully return the book.")
		p=0
	} else {
		fmt.Println("There are no borrowed books with this ID and ISBN.")
		p=1
	}
	return p,nil
}

func (lib *Library) CheckStudent(sid,pss string) bool {
	var p bool = false
	rows, err := lib.db.Query("SELECT PASSWORD FROM STUDENTS WHERE SID=\""+sid+"\"")
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var PASSWORD string
		rows.Scan(&PASSWORD)
		if pss==PASSWORD {
			p = true
		}
	}
	return p
}

func (lib *Library) CheckAdministrator(aid,pss string) bool {
	var p bool = false
	rows, err := lib.db.Query("SELECT PASSWORD FROM ADMINISTRATORS WHERE AID=\""+aid+"\"")
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var PASSWORD string
		rows.Scan(&PASSWORD)
		if pss==PASSWORD {
			p = true
		}
	}
	return p
}

// etc...

func main() {
	fmt.Println("Welcome to the Library Management System!")
	var lib Library 
	lib.ConnectDB()
	lib.CreateTables()
	//lib.init()
	var inp,id,pss,mode,pattern,value,author,ISBN,title,explanation,sname,sid string
	var ls int
	mode = "1"
	for true {
		fmt.Println("Please select login mode:")
		fmt.Println("1.Student 2.Administrator 0.Exit The System")
		fmt.Scanln(&inp)
		if inp=="0" {
			break
		} else if inp=="1" {
			fmt.Println("Please enter your Student ID:")
			fmt.Scanln(&id)
			fmt.Println("Please enter your password:")
			fmt.Scanln(&pss)
			var p bool = lib.CheckStudent(id,pss)
			if !p {
				fmt.Println("Login failed.")
			}
			for p {
				fmt.Println("Please select what you want to do:")
				fmt.Println("1.QueryBook 2.BorrowBook 3.QueryHistory 4.QueryDueDate 5.ExtendDueDate 6.QueryOverDue 7.ReturnBook 0.ExitSystem")
				fmt.Scanln(&mode)
				if mode=="0" {
					break;
				}
				switch  mode{ 
					    case "1": 
						fmt.Println("Please select the Keyword pattern:")
						fmt.Println("1.title 2.author 3.ISBN")
						fmt.Scanln(&pattern)
						if pattern=="1" || pattern=="2" || pattern=="3" {
							fmt.Println("Please enter the Keyword:")
							fmt.Scanln(&value)
							if pattern=="1" {
								lib.QueryBook(value,"title")
							} else if pattern=="2" {
								lib.QueryBook(value,"author")
							} else {
								lib.QueryBook(value,"ISBN")
							}
						} else {
							fmt.Println("Error pattern!")
						}
					    case "2": 
						fmt.Println("Please select the ISBN:")
						fmt.Scanln(&ISBN)
						lib.BorrowBook(id,ISBN)
					    case "3": 
						lib.QueryHistory(id)
					    case "4": 
						fmt.Println("Please select the ISBN:")
						fmt.Scanln(&ISBN)
						lib.QueryDueDate(id,ISBN)
					    case "5": 
						fmt.Println("Please select the ISBN:")
						fmt.Scanln(&ISBN)
						lib.ExtendDueDate(id,ISBN)
					    case "6": 
						lib.QueryOverDue(id)
					    case "7": 
						fmt.Println("Please select the ISBN:")
						fmt.Scanln(&ISBN)
						lib.ReturnBook(id,ISBN)
				} 
			}
			
		} else if inp=="2" {
			fmt.Println("Please enter your Administrator ID:")
			fmt.Scanln(&id)
			fmt.Println("Please enter your password:")
			fmt.Scanln(&pss)
			var p bool = lib.CheckAdministrator(id,pss)
			if !p {
				fmt.Println("Login failed.")
			}
			for p {
				fmt.Println("Please select what you want to do:")
				fmt.Println("1.AddBook 2.DeleteBook 3.AddStudent 4.QueryBook 0.ExitSystem")
				fmt.Scanln(&mode)
				if mode=="0" {
					break;
				}
				switch  mode{ 
					    case "1": 
						fmt.Println("Please enter the title")
						fmt.Scanln(&title)
						fmt.Println("Please enter the author")
						fmt.Scanln(&author)
						fmt.Println("Please enter the ISBN")
						fmt.Scanln(&ISBN)
						lib.AddBook(title,author,ISBN)
					    case "2": 
						fmt.Println("Please enter the ISBN:")
						fmt.Scanln(&ISBN)
						fmt.Println("Is this book in the library now?")
						fmt.Println("0.NO 1.Yes")
						fmt.Scanln(&ls)
						fmt.Println("Please enter the explanation:")
						fmt.Scanln(&explanation)
						lib.DeleteBook(ISBN,explanation,ls)
					    case "3": 
						fmt.Println("Please enter the student name")
						fmt.Scanln(&sname)
						fmt.Println("Please enter the Student ID")
						fmt.Scanln(&sid)
						fmt.Println("Please enter the PASSWORD")
						fmt.Scanln(&pss)
						lib.AddStudent(sname,sid,pss)
					    case "4": 
						fmt.Println("Please select the ISBN:")
						fmt.Scanln(&ISBN)
						lib.QueryBook(id,ISBN)
				} 
			}
		}
		if mode=="0" {
			break;
		}
		
	}
}
