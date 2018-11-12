package go_ibm_db

import (
    "database/sql"
)

type Pool struct{
availablePool   map[string] *sql.DB
usedPool        map[string] *sql.DB
}

//Pconnect will return the pool instance

func Pconnect( )(*Pool){
    p:=&Pool{
        availablePool: make(map[string] *sql.DB,100),
		usedPool     : make(map[string] *sql.DB,100),
        }
    return p
}


type Db *sql.DB
var db Db
var con string
var val Db

//Open will check for the connection in the pool
//If not opens a new connection and stores in the pool

func (p *Pool) Open(Connstr string) *sql.DB{
    if val,ok:=p.availablePool[Connstr];ok{
	    p.usedPool[Connstr]=p.availablePool[Connstr];
		delete(p.availablePool,Connstr)
        return val
    }else{
        dbb, err:=sql.Open("go_ibm_db", Connstr)
        if err != nil{
            return nil
        }
        p.usedPool[Connstr]=dbb;
        db=dbb
        con=Connstr
    }
    return db
}

//Close will make the connection available for the next release

func (p *Pool) Close(){
    p.availablePool[con]=p.usedPool[con]
	delete(p.usedPool,con)
}

//Release will remove all the connections in the pool

func (p *Pool) Release(){
    if(p.availablePool != nil){
        for _,dbCon := range p.availablePool{
            dbCon.Close()
        }
    }else{
        p.availablePool=nil
    }
}
