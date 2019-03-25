package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func downloadFile(filepath string, url string) error {
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

func unzipping(sourcefile string) {
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

func linux_untar(clidriver string) {
	out, _ := exec.Command("tar", "xvzf", clidriver).Output()
	fmt.Println(string(out[:]))
}

func main() {
	var cliFileName string
	var url string
	var i int
	_, a := os.LookupEnv("IBM_DB_DIR")
	_, b := os.LookupEnv("IBM_DB_HOME")
	_, c := os.LookupEnv("IBM_DB_LIB")
	out, _ := exec.Command("go", "env", "GOPATH").Output()
	str := strings.TrimSpace(string(out))
	path := filepath.Join(str, "/src/github.com/ibmdb/go_ibm_db/installer/clidriver")
	if !(a && b && c) {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if runtime.GOOS == "windows" {
				i = 1
				const wordsize = 32 << (^uint(0) >> 32 & 1)
				if wordsize == 64 {
					cliFileName = "ntx64_odbc_cli.zip"
				} else {
					cliFileName = "nt32_odbc_cli.zip"
				}
				fmt.Printf("windows\n")
				fmt.Println(cliFileName)
			} else if runtime.GOOS == "linux" {
				i = 2
				if runtime.GOARCH == "ppc64le" {
					const wordsize = 32 << (^uint(0) >> 32 & 1)
					if wordsize == 64 {
						cliFileName = "ppc64le_odbc_cli.tar.gz"
					}
				} else if runtime.GOARCH == "ppc" {
					const wordsize = 32 << (^uint(0) >> 32 & 1)
					if wordsize == 64 {
						cliFileName = "ppc64_odbc_cli.tar.gz"
					} else {
						cliFileName = "ppc32_odbc_cli.tar.gz"
					}
				} else if runtime.GOARCH == "amd64" {
					const wordsize = 32 << (^uint(0) >> 32 & 1)
					if wordsize == 64 {
						cliFileName = "linuxx64_odbc_cli.tar.gz"
					} else {
						cliFileName = "linuxia32_odbc_cli.tar.gz"
					}
				} else if runtime.GOARCH == "390" {
					const wordsize = 32 << (^uint(0) >> 32 & 1)
					if wordsize == 64 {
						cliFileName = "s390x64_odbc_cli.tar.gz"
					} else {
						cliFileName = "s390_odbc_cli.tar.gz"
					}
				}
				fmt.Printf("linux\n")
				fmt.Println(cliFileName)
			} else if runtime.GOOS == "aix" {
				i = 2
				const wordsize = 32 << (^uint(0) >> 32 & 1)
				if wordsize == 64 {
					cliFileName = "aix64_odbc_cli.tar.gz"
				} else {
					cliFileName = "aix32_odbc_cli.tar.gz"
				}
				fmt.Printf("aix\n")
				fmt.Printf(cliFileName)
			} else if runtime.GOOS == "sunos" {
				i = 2
				if runtime.GOARCH == "i86pc" {
					const wordsize = 32 << (^uint(0) >> 32 & 1)
					if wordsize == 64 {
						cliFileName = "sunamd64_odbc_cli.tar.gz"
					} else {
						cliFileName = "sunamd32_odbc_cli.tar.gz"
					}
				} else if runtime.GOARCH == "SUNW" {
					const wordsize = 32 << (^uint(0) >> 32 & 1)
					if wordsize == 64 {
						cliFileName = "sun64_odbc_cli.tar.gz"
					} else {
						cliFileName = "sun32_odbc_cli.tar.gz"
					}
				}
				fmt.Printf("Sunos\n")
				fmt.Printf(cliFileName)
			} else if runtime.GOOS == "darwin" {
				i = 2
				const wordsize = 32 << (^uint(0) >> 32 & 1)
				if wordsize == 64 {
					cliFileName = "macos64_odbc_cli.tar.gz"
				}
				fmt.Printf("darwin\n")
				fmt.Printf(cliFileName)
			} else {
				fmt.Println("not a known platform")
				os.Exit(1)
			}
			fileUrl := "https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/" + cliFileName
			fmt.Println(url)
			fmt.Println("Downloading...")
			err := downloadFile(cliFileName, fileUrl)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("download successful")
			}
			if i == 1 {
				unzipping(cliFileName)
			} else {
				linux_untar(cliFileName)
			}
		} else {
			fmt.Println("Clidriver Already exits in the directory")
		}
	}
}
