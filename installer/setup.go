package main

import (
    "fmt"
    "runtime"
	"io"
	"net/http"
	"os"
	"os/exec"
	"archive/zip"
	"archive/tar"
	"compress/gzip"
	"path/filepath"
	"strings"
)
//func to Download
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
//func to unzip
func Unzipping(sourcefile string) {
    reader, err := zip.OpenReader(sourcefile)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer reader.Close()
    for _, f := range reader.Reader.File {
        zipped, err := f.Open()
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        defer zipped.Close()
        path := filepath.Join("./", f.Name)
        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
            fmt.Println("Creating directory", path)
        } else {
            writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, f.Mode())
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
            defer writer.Close()
            if _, err = io.Copy(writer, zipped); err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
            fmt.Println("Unzipping : ", path)
        }
    }
}
//Func to Untar
func Untaring(sourcefile string) {
    file, err := os.Open(sourcefile)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer file.Close()
    var fileReader io.ReadCloser = file
    if strings.HasSuffix(sourcefile, ".gz") {
    if fileReader, err = gzip.NewReader(file); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer fileReader.Close()
    }
    tarBallReader := tar.NewReader(fileReader)
    for {
        header, err := tarBallReader.Next()
        if err != nil {
            if err == io.EOF {
                break
            }
            fmt.Println(err)
            os.Exit(1)
        }
        filename := header.Name
        switch header.Typeflag {
            case tar.TypeDir:
                fmt.Println("Creating directory :", filename)
                err = os.MkdirAll(filename, os.FileMode(header.Mode))
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }
            case tar.TypeReg:
                fmt.Println("Untarring :", filename)
                writer, err := os.Create(filename)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }
                io.Copy(writer, tarBallReader)
                err = os.Chmod(filename, os.FileMode(header.Mode))
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }
                writer.Close()
            default:
                fmt.Printf("Unable to untar type : %c in file %s", header.Typeflag, filename)
        }
    }
}

func linux_untar(clidriver string){
 out, _:= exec.Command("tar","xvzf",clidriver).Output()
fmt.Println(string(out[:]))
}


func main() {
    var cliFileName string
    var url string
    i:=0
    _,a :=os.LookupEnv("IBM_DB_DIR")
    _,b :=os.LookupEnv("IBM_DB_HOME")
    _,c :=os.LookupEnv("IBM_DB_LIB")
    if(!(a && b && c)){
    if runtime.GOOS == "aix" {
	    i=2
        const wordsize = 32 << (^uint(0) >> 32 & 1)
	    if wordsize==64 {
        cliFileName = "aix64_odbc_cli.tar.gz"
        } else {
            cliFileName = "aix32_odbc_cli.tar.gz"
        }
	    fmt.Printf("aix\n")
	    fmt.Printf(cliFileName)
	}else if runtime.GOOS == "linux"{
	       i=2
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
	    i=1
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
            cliFileName = "macos64_odbc_cli.tar.gz"
        }
	    fmt.Printf("darwin\n")
	    fmt.Printf(cliFileName)
	}else if runtime.GOOS =="sunos"{
	    i=2
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
	        fmt.Println("download successful")
	    }	
	    if(i == 1){
	        Unzipping(cliFileName)
	    }else if(i==2){
		linux_untar(cliFileName)

		}else {
	        Untaring(cliFileName)
		}
    }
}	