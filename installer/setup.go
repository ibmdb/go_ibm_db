package main

import (
    "fmt"
    "runtime"
	"io"
	"net/http"
	"os"
	"archive/zip"
	 "path/filepath"
    "strings"
	"log"
)

func DownloadFile(filepath string, url string) error {
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}
func Unzip(src string, dest string) ([]string, error) {

    var filenames []string

    r, err := zip.OpenReader(src)
    if err != nil {
        return filenames, err
    }
    defer r.Close()

    for _, f := range r.File {
        rc, err := f.Open()
        if err != nil {
            return filenames, err
        }
        defer rc.Close()
        fpath := filepath.Join(dest, f.Name)
        if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
            return filenames, fmt.Errorf("%s: illegal file path", fpath)
        }

        filenames = append(filenames, fpath)
        if f.FileInfo().IsDir() {

            // Make Folder
            os.MkdirAll(fpath, os.ModePerm)

        } else {

            // Make File
            if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
                return filenames, err
            }

            outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return filenames, err
            }

            _, err = io.Copy(outFile, rc)

            // Close the file without defer to close before next iteration of loop
            outFile.Close()

            if err != nil {
                return filenames, err
            }

        }
    }
    return filenames, nil
}

func main() {
var cliFileName string
var url string

_,a :=os.LookupEnv("IBM_DB_DIR")
_,b :=os.LookupEnv("IBM_DB_HOME")
_,c :=os.LookupEnv("IBM_DB_LIB")
 if(!(a && b && c)){
    if runtime.GOOS == "aix" {
      const wordsize = 32 << (^uint(0) >> 32 & 1)
	  if wordsize==64 {
        cliFileName = "aix64_odbc_cli.tar.gz"
    } else {
         cliFileName = "aix32_odbc_cli.tar.gz"
    }
	fmt.Printf("aix\n")
	fmt.Printf(cliFileName)
	}else if runtime.GOOS == "linux"{
	 if runtime.GOARCH == "ppc64le" {
	 const wordsize = 32 << (^uint(0) >> 32 & 1)
	 if wordsize==64{
	 cliFileName= "ppc64le_odbc_cli.tar.gz"
	 }
	 }else if runtime.GOARCH == "ppc" {
	 const wordsize = 32 << (^uint(0) >> 32 & 1)
	 if wordsize==64{
	 cliFileName="ppc64_odbc_cli.tar.gz"
	 }else{
	 cliFileName="ppc32_odbc_cli.tar.gz"
	 }
	 }else if runtime.GOARCH == "amd64" {
	  const wordsize = 32 << (^uint(0) >> 32 & 1)
	 if wordsize==64{
	 cliFileName="linuxx64_odbc_cli.tar.gz"
	 }else{
	 cliFileName="linuxia32_odbc_cli.tar.gz"
	 }
	 }else if runtime.GOARCH == "390" {
	  const wordsize = 32 << (^uint(0) >> 32 & 1)
	 if wordsize==64{
	 cliFileName="s390x64_odbc_cli.tar.gz"
	 }else{
	 cliFileName="s390_odbc_cli.tar.gz"
	 }
	 }
	 fmt.Printf("linux\n")
     fmt.Printf(cliFileName)
	 }else if runtime.GOOS=="windows"{
	 const wordsize = 32 << (^uint(0) >> 32 & 1)
	  if wordsize==64 {
        cliFileName = "ntx64_odbc_cli.zip"
    } else {
         cliFileName = "nt32_odbc_cli.zip"
    }
	fmt.Printf("windows\n")
	fmt.Printf(cliFileName)
	 }else if  runtime.GOOS=="darwin"{
	 const wordsize = 32 << (^uint(0) >> 32 & 1)
	  if wordsize==64 {
        cliFileName = "ntx64_odbc_cli.zip"
    }
	fmt.Printf("darwin\n")
	fmt.Printf(cliFileName)
	}else if runtime.GOOS =="sunos"{
	 if runtime.GOARCH == "i86pc"{
	const wordsize = 32 << (^uint(0) >> 32 & 1)
	  if wordsize==64 {
	  cliFileName = "sunamd64_odbc_cli.tar.gz"
	  }else{
     cliFileName = "sunamd32_odbc_cli.tar.gz"
	 }
	 }else if runtime.GOARCH == "SUNW"{
	const wordsize = 32 << (^uint(0) >> 32 & 1)
	  if wordsize==64 {
	  cliFileName = "sun64_odbc_cli.tar.gz"
	  }else{
     cliFileName = "sun32_odbc_cli.tar.gz"
	 }
	 }
	 fmt.Printf("Sunos\n")
	 fmt.Printf(cliFileName)
	 }else{
	fmt.Println("not a known platform")
	}
	fileUrl:= "https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/" + cliFileName
	fmt.Println(url)
	fmt.Println("Downloading...")
	err:=DownloadFile(cliFileName,fileUrl)
	if err!=nil{
	fmt.Println(err)
	}else{
	fmt.Printf("download successful")
	}
	files, err := Unzip("ntx64_odbc_cli.zip", "ntx64_odbc_cli")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))
}
}	


