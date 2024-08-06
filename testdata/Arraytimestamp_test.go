package main

import (
	"fmt"
	"time"
	"strings"
	"testing"
)

func TestTimeStampArray(t *testing.T) {
	if TimeStampArray() != nil {
		t.Error("Error at StringArray")
	}
}

//TimeStampArray function performs inserting float32,float64 datatypes.
func TimeStampArray() error {
        db := Createconnection()
        defer db.Close()

        db.Exec("Drop table arr")
        _, err := db.Exec("create table arr(var1 timestamp, var2 time, var3 date)")
        if err != nil {
                fmt.Println("Exec error: ", err)
                return err
        }
        a := []time.Time{}
        for i := 0; i < 5; i++ {
                a = append(a, time.Now())
                time.Sleep(1 * time.Second)
        }
        st, err := db.Prepare("Insert into arr values(?,?,?)")
        defer st.Close()
        if err != nil {
                fmt.Println("Prepare error: ", err)
                return err
        }
        _, err = st.Query(a, a, a)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []timestamp")
                return err
        }
        return nil
}

