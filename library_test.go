package main

import (
	"testing"
)

func TestCreateTables(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.CreateTables()
	if err != nil {
		t.Errorf("can't create tables")
	}
	lib.init()
}

func TestAddBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	all := []struct {
		title string
		author string
		ISBN string
		p int
	}{
		{"b03","a03","115",0},
		{"b04","a04","116",1},
	}
	for _, nw := range all {
		var p int
		p,err := lib.AddBook(nw.title,nw.author,nw.ISBN)
		if err != nil {
			t.Errorf("can't AddBook")
		}
		if p != nw.p {
			t.Errorf("p was incorrect, got: %d, want: %d.", p, nw.p)
		}
	}
}

func TestDeleteBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	all := []struct {
		ISBN string
		EXPLAINATION string
		pattern int
		p int
	}{
		{"117","lost",0,4},
		{"112","lost",0,0},
	}
	for _, nw := range all {
		var p int
		p,err := lib.DeleteBook(nw.ISBN,nw.EXPLAINATION,nw.pattern)
		if err != nil {
			t.Errorf("can't AddBook")
		}
		if p != nw.p {
			t.Errorf("p was incorrect, got: %d, want: %d.", p, nw.p)
		}
	}
}

func TestAddStudent(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	all := []struct {
		name string
		sid string
		pss string
		p int
	}{
		{"aaa","s01","paa",0},
		{"bbb","s04","pbb",1},
	}
	for _, nw := range all {
		var p int
		p,err := lib.AddStudent(nw.name,nw.sid,nw.pss)
		if err != nil {
			t.Errorf("can't AddStudent")
		}
		if p != nw.p {
			t.Errorf("p was incorrect, got: %d, want: %d.", p, nw.p)
		}
	}
}


func TestQueryBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	all := []struct {
		value string
		pattern string
		p string
	}{
		{"112","ISBN"," cnt:1 ISBN=112 author=a01 title=b01 TOTAL_NUM=1 AVAILABLE_NUM=1"},
		{"a01","author"," cnt:1 ISBN=112 author=a01 title=b01 TOTAL_NUM=1 AVAILABLE_NUM=1 cnt:2 ISBN=114 author=a01 title=b02 TOTAL_NUM=4 AVAILABLE_NUM=1"},
		{"b01","title"," cnt:1 ISBN=112 author=a01 title=b01 TOTAL_NUM=1 AVAILABLE_NUM=1 cnt:2 ISBN=113 author=a02 title=b01 TOTAL_NUM=2 AVAILABLE_NUM=0"},
	}
	for _, nw := range all {
		var p string
		p,err := lib.QueryBook(nw.value,nw.pattern)
		if err != nil {
			t.Errorf("can't AddStudent")
		}
		if p != nw.p {
			t.Errorf("p was incorrect, got: %s, want: %s.", p, nw.p)
		}
	}
}


func TestBorrowBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	all := []struct {
		sid string
		ISBN string
		p int
	}{
		{"s01","113",0},
		{"s02","113",1},
	}
	for _, nw := range all {
		var p int
		p,err := lib.BorrowBook(nw.sid,nw.ISBN)
		if err != nil {
			t.Errorf("can't BorrowBook")
		}
		if p != nw.p {
			t.Errorf("p was incorrect, got: %d, want: %d.", p, nw.p)
		}
	}
}

func TestQueryHistory(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	all := []struct {
		sid string
		p string
	}{
		{"s01"," cnt:1 SID=s01 ISBN=115 BORROW_DATE=2020-01-03 DUE_DATE=2020-05-01 RETURN_DATE=NULL cnt:2 SID=s01 ISBN=113 BORROW_DATE=2020-04-01 DUE_DATE=2020-05-01 RETURN_DATE=NULL cnt:3 SID=s01 ISBN=114 BORROW_DATE=2020-04-02 DUE_DATE=2020-05-02 RETURN_DATE=NULL"},
		{"s02"," cnt:1 SID=s02 ISBN=114 BORROW_DATE=2020-04-01 DUE_DATE=2020-05-01 RETURN_DATE=NULL cnt:2 SID=s02 ISBN=114 BORROW_DATE=2020-04-01 DUE_DATE=2020-05-01 RETURN_DATE=NULL"},
		
	}
	for _, nw := range all {
		var p string
		p,err := lib.QueryHistory(nw.sid)
		if err != nil {
			t.Errorf("can't AddStudent")
		}
		if p != nw.p {
			t.Errorf("p was incorrect, got: %s, want: %s.", p, nw.p)
		}
	}
}

func TestQueryDueDate(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	all := []struct {
		sid string
		ISBN string
		p string
	}{
		{"s01","114","2020-05-02"},
		{"s02","115","error"},
	}
	for _, nw := range all {
		var p string
		p,err := lib.QueryDueDate(nw.sid,nw.ISBN)
		if err != nil {
			t.Errorf("can't QueryDueDate")
		}
		if p != nw.p {
			t.Errorf("p was incorrect, got: %s, want: %s.", p, nw.p)
		}
	}
}

func TestExtendDueDate(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	all := []struct {
		sid string
		ISBN string
		p int
	}{
		{"s01","115",0},
		{"s03","113",1},
	}
	for _, nw := range all {
		var p int
		p,err := lib.ExtendDueDate(nw.sid,nw.ISBN)
		if err != nil {
			t.Errorf("can't ExtendDueDate")
		}
		if p != nw.p {
			t.Errorf("p was incorrect, got: %d, want: %d.", p, nw.p)
		}
	}
}

func TestQueryOverDue(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	all := []struct {
		sid string
		p string
	}{
		{"s01"," cnt:1 ISBN=113 DUE_DATE=2020-05-01 ExtendableTimes=3 cnt:2 ISBN=114 DUE_DATE=2020-05-02 ExtendableTimes=3 cnt:3 ISBN=115 DUE_DATE=2020-05-01 ExtendableTimes=0"},
		{"s02"," cnt:1 ISBN=114 DUE_DATE=2020-05-01 ExtendableTimes=3 cnt:2 ISBN=114 DUE_DATE=2020-05-01 ExtendableTimes=3"},
	}
	for _, nw := range all {
		var p string
		p,err := lib.QueryOverDue(nw.sid)
		if err != nil {
			t.Errorf("can't QueryDueDate")
		}
		if p != nw.p {
			t.Errorf("p was incorrect, got: %s, want: %s.", p, nw.p)
		}
	}
}

func TestReturnBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	all := []struct {
		sid string
		ISBN string
		p int
	}{
		{"s01","112",1},
		{"s01","113",0},
	}
	for _, nw := range all {
		var p int
		p,err := lib.ReturnBook(nw.sid,nw.ISBN)
		if err != nil {
			t.Errorf("can't ReturnBook")
		}
		if p != nw.p {
			t.Errorf("p was incorrect, got: %d, want: %d.", p, nw.p)
		}
	}
}

