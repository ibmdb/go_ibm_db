package main

import (
    "fmt"
	a "github.com/ibmdb/go_ibm_db"
	"database/sql"
	)
	
var con ="HOSTNAME=host;PORT=number;DATABASE=name;UID=username;PWD=password"
	
func Createconnection()(db *sql.DB){
    db,_=sql.Open("go_ibm_db",con)
    return db
	}


func Createtable( ) error{
    db, err:=sql.Open("go_ibm_db", con)
    db.Exec("DROP table rocket")
    _,err=db.Exec("create table rocket(a int)")
	_,err=db.Exec("create table rocket1(a int)")
    if err != nil{
        return err
    }
    return nil
}
func Insert() error{
    db,err:=sql.Open("go_ibm_db",con)
    _,err =db.Exec("insert into rocket values(1)")
	if err != nil{
        return err
    }
    return nil
}

func Drop() error{
    db,err:=sql.Open("go_ibm_db",con)
    _,err =db.Exec("drop table rocket1")
	if err != nil{
        return err
    }
    return nil
}

func Prepare() error{
db,_:=sql.Open("go_ibm_db",con)
_,err:=db.Prepare("select * from rocket")
if err!=nil{
return err
}
return nil
}

func Query() error{
db,_:=sql.Open("go_ibm_db",con)
st,_:=db.Prepare("select * from rocket")
_,err:=st.Query()
if err != nil{
return err
} 
return nil
}

func Scan() error{
db,_:=sql.Open("go_ibm_db",con)
st,_:=db.Prepare("select * from rocket")
rows,err:=st.Query()
for rows.Next(){
var a string
err = rows.Scan(&a)
if err != nil{
return err
}
}
return nil
}

func Next() error{
db,_:=sql.Open("go_ibm_db",con)
st,_:=db.Prepare("select * from rocket")
rows,err:=st.Query()
for rows.Next(){
var a string
err = rows.Scan(&a)
if err != nil{
return err
}
}
return nil
}

func Columns() error{
db,_:=sql.Open("go_ibm_db",con)
st,_:=db.Prepare("select * from rocket")
rows,_:=st.Query()
_,err := rows.Columns()
if err != nil{
return err
}
for rows.Next(){
var a string
_ = rows.Scan(&a)
}
return nil
}



func Queryrow() error{
a:=1
var uname int
db,err:=sql.Open("go_ibm_db",con)
err=db.QueryRow("select a from rocket where a=?",a).Scan(&uname)
if err != nil{
return err
}
return nil
}

func Begin() error{
db,err:=sql.Open("go_ibm_db",con)
_,err=db.Begin()
if err != nil{
return err
}
return nil
}

func Commit() error{
    db,err:=sql.Open("go_ibm_db",con)
    bg,err:=db.Begin()
	db.Exec("DROP table u")
	_,err=bg.Exec("create table u(id int)")
	err=bg.Commit()
	if err!=nil{
	return err
	}
	return nil
}

func Close()(error){
    db,_:=sql.Open("go_ibm_db",con)
    err:=db.Close()
	if err!=nil{
	return err
	}
	return nil
	}

func PoolOpen() (int){
    pool:=a.Pconnect("PoolSize=50")
    db:=pool.Open(con)
    if(db == nil){
        return 0
    }else {
        return 1
    } 
}
	
func main(){
result:=Createconnection()
if(result != nil){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}

result1:=Createtable()
if(result1 == nil){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}

result2:=Insert()
if(result2 == nil){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}
result3:=Drop()
if(result3 == nil){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}
result4:=Prepare()
if(result4 == nil){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}
result5:=Query()
if(result5 == nil){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}

result6:=Scan()
if(result6 == nil){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}
result7:=Next()
if(result7 == nil){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}

result8:=Columns()
if(result8 == nil){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}

result9:=Queryrow()
if(result9 == nil){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}
result10:=Begin()
if(result10== nil){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}
result11:=Commit()
if(result11 == nil){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}
result12:=Close()
if(result12 == nil){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}
result13:=PoolOpen()
if(result13 ==1){
fmt.Println("Pass")
} else {
fmt.Println("fail")
}
}