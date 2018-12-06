package go_ibm_db

import (
    "database/sql"
	"fmt"
	"time"
	"strconv"
	"strings"
)

type DBP struct{
sql.DB
con string
n time.Duration
}

type Pool struct{
availablePool   map[string] []*DBP
usedPool        map[string] []*DBP
PoolSize        int
}

var b *Pool
var ConnMaxLifetime,PoolSize int

//Pconnect will return the pool instance
func Pconnect(PoolSize string) (*Pool){
    var Size int
    count:=len(PoolSize)
	if count>0{
        opt:=strings.Split(PoolSize,"=")
	    if(opt[0] == "PoolSize"){    
		    Size,_=strconv.Atoi(opt[1])
	    }else{
	        fmt.Println("Not a valid parameter")
	    }
	} else {
	    Size=100
	}
    p:=&Pool{
	    availablePool: make(map[string] []*DBP),
		usedPool     : make(map[string] []*DBP),
		PoolSize     : Size,
    }
	b=p
	
    return p
} 

var Psize int

//Open will check for the connection in the pool
//If not opens a new connection and stores in the pool

func (p *Pool) Open(Connstr string,options ...string)(*DBP){
    var Time time.Duration
    count := len(options)
    if count>0{
        for i:=0;i<count;i++{
		    opt:=strings.Split(options[i],"=")
			if(opt[0] == "SetConnMaxLifetime"){
			    ConnMaxLifetime,_=strconv.Atoi(opt[1])
				Time=time.Duration(ConnMaxLifetime) * time.Second
			} else {
			    fmt.Println("not a valid parameter")
			}
		}
    } else {
        Time=30*time.Second
    }
    if(Psize<p.PoolSize){
	    Psize=Psize+1;
        if val,ok:=p.availablePool[Connstr];ok{
	        if(len(val) > 1){
	            dbpo:=val[0]
                copy(val[0:],val[1:])
		        val[len(val)-1]=nil
		        val=val[:len(val)-1]
		        p.availablePool[Connstr]=val
		        p.usedPool[Connstr]=append(p.usedPool[Connstr],dbpo)
				dbpo.SetConnMaxLifetime(Time)
		        return dbpo
		    }else{
		        dbpo:=val[0]
		        p.usedPool[Connstr]=append(p.usedPool[Connstr],dbpo)
		        delete(p.availablePool,Connstr)
				dbpo.SetConnMaxLifetime(Time)
		        return dbpo
		    }
        }else{
	        db,err:=sql.Open("go_ibm_db",Connstr)
		    if err != nil{
		        return nil
		    }
		    dbi:=&DBP{
		        DB  :*db,
			    con :Connstr,
				n   :Time,
		    }
		    p.usedPool[Connstr]=append(p.usedPool[Connstr],dbi)
            dbi.SetConnMaxLifetime(Time)			
            return dbi	    
		}
	} else {
	    db,err:=sql.Open("go_ibm_db",Connstr)
		if err != nil{
		    return nil
		}
		dbi:=&DBP{
		    DB  :*db,
			con :Connstr,
		}
		return dbi
	}
}

//Close will make the connection available for the next release

func (d *DBP) Close(){
    Psize=Psize-1
	var pos int
	i:=-1
	if valc,okc:=b.usedPool[d.con];okc{
	    if(len(valc) > 1){
	        for _,b:=range valc{
		        i=i+1
			    if b == d{
			        pos = i
			    }
		    }
		    dbpc:=valc[pos]
		    copy(valc[pos:], valc[pos+1:])
		    valc[len(valc)-1] = nil
		    valc = valc[:len(valc)-1]
			b.usedPool[d.con]=valc
		    b.availablePool[d.con]=append(b.availablePool[d.con],dbpc)
	    } else {
		    dbpc := valc[0]
            b.availablePool[d.con]=append(b.availablePool[d.con],dbpc)
	        delete(b.usedPool,d.con)
	    }
		go d.Timeout()
	} else {
	    d.DB.Close()
	}
}

//Timeout for closing the connection in pool

func (d *DBP) Timeout(){
    var pos int
	i:=-1
	select {
	case <-time.After(d.n):
	    if valt,okt:=b.availablePool[d.con];okt{
            if(len(valt) > 1){
	            for _,b:=range valt{
		            i=i+1
			        if b == d{
			            pos = i
			        }
		        }
		        dbpt:=valt[pos]
		        copy(valt[pos:], valt[pos+1:])
		        valt[len(valt)-1] = nil
		        valt = valt[:len(valt)-1]
				b.availablePool[d.con]=valt
			    dbpt.DB.Close()
		    }else{
		        dbpt:=valt[0]
			    dbpt.DB.Close()
		        delete(b.availablePool,d.con)
		    }
        }
	}
}

//Release will close all the connections in the pool

func (p *Pool) Release(){
    if(p.availablePool != nil){
	    for _,vala := range p.availablePool{
		    for _,dbpr := range vala{
			    dbpr.DB.Close()
			}
		}
		p.availablePool=nil
	}
	if(p.usedPool != nil){
	    for _,valu := range p.usedPool{
		    for _,dbpr := range valu{
			    dbpr.DB.Close()				
			}
		}
		p.usedPool=nil
	}
}

// Display will print the  values in the map 

func (p *Pool) Display(){
    fmt.Println(p.availablePool)
    fmt.Println(p.usedPool)
    fmt.Println(p.PoolSize)
}
