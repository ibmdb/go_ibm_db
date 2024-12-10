package main

import (
	"archive/zip"
	"bufio"
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
	out, err := os.Create("../../" + filepath)
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

func Unzipping(sourcefile string, targetDirectory string) {
	reader, err := zip.OpenReader("./../../" + sourcefile)
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
		path := filepath.Join(targetDirectory, f.Name)
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

func linux_untar(clidriver string, targetDirectory string, ibmdbDir string) error {
	fmt.Printf("Extracting with tar -xzf %s -C %s\n", clidriver, targetDirectory)
	out, err := exec.Command("tar", "xzf", "./../../" + clidriver, "-C", targetDirectory).Output()

	fmt.Println(string(out))
	if err != nil {
		fmt.Println("Error while running tar: " + err.Error())
		return err
	} else {
        fmt.Println("clidriver path = " + targetDirectory + "/clidriver")
        fmt.Println("Now run below commands:")
        fmt.Println("export IBM_DB_HOME=" + targetDirectory + "/clidriver")
        fmt.Println("source " + ibmdbDir + "/installer/setenv.sh")
    }

	if runtime.GOOS == "darwin" {
        // Create symlinks for libdb2
        _, _ = exec.Command("ln", "-s", targetDirectory + "/clidriver/lib/libdb2.dylib", targetDirectory + "/libdb2.dylib").Output()
        _, _ = exec.Command("ln", "-s", targetDirectory + "/clidriver/lib/libdb2.dylib", ibmdbDir + "/libdb2.dylib").Output()
        _, _ = exec.Command("ln", "-s", targetDirectory + "/clidriver/lib/libdb2.dylib", ibmdbDir + "/testdata/libdb2.dylib").Output()
    }

	return nil
}

func aix_untar(clidriver string, targetDirectory string) error {
	fmt.Printf("Extracting with gunzip %s \n", clidriver)
	gunzipOut, err := exec.Command("gunzip", "./../../"+clidriver).Output()

	fmt.Println(string(gunzipOut))
	if err != nil {
		fmt.Println("Error while running gunzip: " + err.Error())
		return err
	}

	clidriver = strings.TrimRight(clidriver, ".gz")

	fmt.Printf("Extracting with tar -xvf %s -C %s\n", clidriver, targetDirectory)
	tarOut, err := exec.Command("tar", "xvf", "./../../"+clidriver, "-C", targetDirectory).Output()

	fmt.Println(string(tarOut))
	if err != nil {
		fmt.Println("Error while running tar: " + err.Error())
		return err
	} else {
        fmt.Println("clidriver path = " + targetDirectory + "/clidriver")
        fmt.Println("Now run below commands:")
        fmt.Println("export IBM_DB_HOME=" + targetDirectory + "/clidriver")
        fmt.Println("export CGO_CFLAGS==\"-I$IBM_DB_HOME/include\"")
        fmt.Println("export CGO_LDFLAGS==\"-L$IBM_DB_HOME/lib\"")
	}

	return nil
}

func getinstalledpath(validateout string) {
	var line string

	scanner := bufio.NewScanner(strings.NewReader(validateout))

	for scanner.Scan() {
		line = scanner.Text()

		if strings.Contains(line, "Install") {
			fields := strings.Split(line, " ")
			fmt.Println(fields[7])
			input1 := fields[7][0:len(fields[7])]
			fmt.Println("Clidriver is already present")
			fmt.Println("Please set IBM_DB_HOME to ", input1)
		}
	}
}

func checkincludepath(includepath string) bool {

	if _, err1 := os.Stat(includepath + "/include"); !os.IsNotExist(err1) {
		//fmt.Println("clidriver/include folder exists in the path")
		if _, err2 := os.Stat(includepath + "/lib"); !os.IsNotExist(err2) {
			//fmt.Println("clidriver/lib folder exists in the path")
			return true
		}
	}

	return false
}

func main() {
    var target, ibmdbDir, cliFileName string
    var unpackageType int
    var err11 error
    var out []byte

    fmt.Println("os =",runtime.GOOS, ", processor =", runtime.GOARCH)
    path, ok := os.LookupEnv("DB2HOME")
    if ok {
        fmt.Println("NOTE: Environment variable DB2HOME name is changed to IBM_DB_HOME.")
        fmt.Println("DB2HOME environment variable is set to", path)
    }

    out, err11 = exec.Command("db2cli", "validate").Output()
    if err11 != nil {
        _, ok := os.LookupEnv("IBM_DB_HOME")
        if !ok {
            if runtime.GOOS == "windows" {
                fmt.Println("Please set IBM_DB_HOME and add %IBM_DB_HOME%/bin to PATH and %IBM_DB_HOME%/lib to LIB environment variables after clidriver installation")
            } else if runtime.GOOS == "aix" {
                fmt.Println("Please set IBM_DB_HOME, CGO_CFLAGS, CGO_LDFLAGS and LIBPATH environment variables after clidriver installation")
            } else if runtime.GOOS == "darwin" {
                fmt.Println("Please set IBM_DB_HOME, CGO_CFLAGS, CGO_LDFLAGS and DYLD_LIBRARY_PATH environment variables after clidriver installation")
            } else {
                fmt.Println("Please set IBM_DB_HOME, CGO_CFLAGS, CGO_LDFLAGS and LD_LIBRARY_PATH environment variables after clidriver installation")
            }
        }
    } else {
        path, ok := os.LookupEnv("IBM_DB_HOME")
        if !ok {
            //set IBM_DB_HOME
            getinstalledpath(string(out))
            os.Exit(1)
        } else {
            fmt.Println("clidriver folder exists in the path....", path)
            if checkincludepath(path) {
                os.Exit(1)
            }
        }
    }

    _, setupFile, _, ok := runtime.Caller(0)
    if ok {
        ibmdbDir = filepath.Dir(setupFile) + "/.."
    } else {
        ibmdbDir = "./.."
    }
    target = ibmdbDir + "/.."
    ibmdbDir, _ = filepath.Abs(ibmdbDir)
    target, _ = filepath.Abs(target)

    if len(os.Args) == 2 {
        target = os.Args[1]
    }

	if _, err1 := os.Stat(target + "/clidriver"); !os.IsNotExist(err1) {
		fmt.Println("clidriver folder exists in the path")

		if _, err2 := os.Stat(target + "/clidriver/include"); !os.IsNotExist(err2) {
			//fmt.Println("clidriver/include folder exists in the path")

			if _, err3 := os.Stat(target + "/clidriver/lib"); !os.IsNotExist(err3) {
				fmt.Println("==> Direcotry \"" + target + "/clidriver\" exists.")
				os.Exit(2)
			} else {
				fmt.Println(target+"/clidriver/lib folder does not exist, installing clidriver ....")
			}
		} else {
			fmt.Println(target+"/clidriver/include folder does not exist, installing clidriver ....")
		}
	}

	if runtime.GOOS == "windows" {
		unpackageType = 1
		const wordsize = 32 << (^uint(0) >> 32 & 1)
		if wordsize == 64 {
			cliFileName = "ntx64_odbc_cli.zip"
		} else {
			cliFileName = "nt32_odbc_cli.zip"
		}
		fmt.Println(cliFileName)
	} else if runtime.GOOS == "linux" {
		unpackageType = 2
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
		} else if runtime.GOARCH == "s390x" {
			const wordsize = 32 << (^uint(0) >> 32 & 1)
			if wordsize == 64 {
				cliFileName = "s390x64_odbc_cli.tar.gz"
			} else {
				cliFileName = "s390_odbc_cli.tar.gz"
			}
		}
		fmt.Println(cliFileName)
	} else if runtime.GOOS == "aix" {
		unpackageType = 3
		const wordsize = 32 << (^uint(0) >> 32 & 1)
		if wordsize == 64 {
			cliFileName = "aix64_odbc_cli.tar.gz"
		} else {
			cliFileName = "aix32_odbc_cli.tar.gz"
		}
		fmt.Println(cliFileName)
	} else if runtime.GOOS == "sunos" {
		unpackageType = 2
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
		fmt.Printf(cliFileName)
	} else if runtime.GOOS == "darwin" {
        unpackageType = 2
        const wordsize = 32 << (^uint(0) >> 32 & 1)
        if wordsize == 64 {
            if runtime.GOARCH == "arm64" {
                cliFileName = "macarm64_odbc_cli.tar.gz"
            } else {
                cliFileName = "macos64_odbc_cli.tar.gz"
            }
        }
    } else {
        fmt.Println("Error: Unsupported platform!")
        os.Exit(3)
    }
	fileUrl := downloadUrl(cliFileName)
	fmt.Println("Downloading " + fileUrl)
	err := downloadFile(cliFileName, fileUrl)
	if err != nil {
		fmt.Println("Error while downloading file: " + err.Error())
		os.Exit(4)
	}
	fmt.Println("download successful")

	if unpackageType == 1 {
		Unzipping(cliFileName, target)
	} else if unpackageType == 3 {
		aix_untar(cliFileName, target)
	} else {
		linux_untar(cliFileName, target, ibmdbDir)
	}
}

func downloadUrl(cliFileName string) string {
	downloadUrl, downloadUrlFound := os.LookupEnv("IBM_DB_DOWNLOAD_URL")
	if !downloadUrlFound {
	    clidriverVersion, ok := os.LookupEnv("CLIDRIVER_DOWNLOAD_VERSION")
        if ok {
            if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" && strings.HasPrefix(strings.ToLower(clidriverVersion), "v11") {
		        downloadUrl = "https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/" + cliFileName
            } else if runtime.GOOS == "darwin" && runtime.GOARCH != "arm64" && strings.HasPrefix(strings.ToLower(clidriverVersion), "v12") {
		        downloadUrl = "https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/" + cliFileName
            } else {
		        downloadUrl = "https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/" + clidriverVersion + "/" + cliFileName
            }
        } else {
            if runtime.GOOS == "darwin" {
		        downloadUrl = "https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/" + cliFileName
            } else {
		        downloadUrl = "https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/v11.5.9/" + cliFileName
            }
        }
	}
	return downloadUrl
}
