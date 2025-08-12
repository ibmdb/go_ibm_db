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

const (
    primaryURL   = "https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/"
    fallbackURL  = "https://github.com/ibmdb/db2drivers/raw/main/clidriver/"
)

func fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}

func dirExists(path string) bool {
    info, err := os.Stat(path)
    return err == nil && info.IsDir()
}

func downloadFile(url string, filename string) error {
    // Create HTTP GET request
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("Failed to fetch %s: %w", url, err)
    }
    defer resp.Body.Close()

    // Check for non-200 status
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Failed to download %s: HTTP %d", url, resp.StatusCode)
    }

    // Create destination file
    out, err := os.Create("../../" + filename)
    if err != nil {
        return fmt.Errorf("Failed to create file %s: %w", filename, err)
    }
    defer out.Close()

    // Copy data
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return fmt.Errorf("Failed to save file %s: %w", filename, err)
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
    fmt.Printf("Extracting with tar -xzf %s -C %s", clidriver, targetDirectory)
    out, err := exec.Command("tar", "xzf", "./../../" + clidriver, "-C", targetDirectory).Output()

    fmt.Println(string(out))
    if err != nil {
        fmt.Println("Error while running tar: " + err.Error())
        return err
    } else {
        fmt.Println("Extracted clidriver path = " + targetDirectory + "/clidriver")
        fmt.Println("\nNow run below commands:")
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
    fmt.Printf("Extracting with gunzip %s", clidriver)
    gunzipOut, err := exec.Command("gunzip", "./../../"+clidriver).Output()

    fmt.Println(string(gunzipOut))
    if err != nil {
        fmt.Println("Error while running gunzip: " + err.Error())
        return err
    }

    clidriver = strings.TrimRight(clidriver, ".gz")

    fmt.Printf("Extracting with tar -xvf %s -C %s", clidriver, targetDirectory)
    tarOut, err := exec.Command("tar", "xvf", "./../../"+clidriver, "-C", targetDirectory).Output()

    fmt.Println(string(tarOut))
    if err != nil {
        fmt.Println("Error while running tar: " + err.Error())
        return err
    } else {
        fmt.Println("Extracted clidriver path = " + targetDirectory + "/clidriver")
        fmt.Println("\nNow run below commands:")
        fmt.Println("export IBM_DB_HOME=" + targetDirectory + "/clidriver")
        fmt.Println("export CGO_CFLAGS==\"-I$IBM_DB_HOME/include\"")
        fmt.Println("export CGO_LDFLAGS==\"-L$IBM_DB_HOME/lib\"")
    }

    return nil
}

func checkInstalledDriver(validateout string) bool {
    var line string

    scanner := bufio.NewScanner(strings.NewReader(validateout))

    for scanner.Scan() {
        line = scanner.Text()

        if strings.Contains(line, "Install") {
            fields := strings.Split(line, " ")
            //fmt.Println(fields[7])
            input1 := fields[7][0:len(fields[7])]
            fmt.Println("Found installed Db2 Client as", input1)
            if checkIncludeDir(input1) {
                fmt.Println("Please set IBM_DB_HOME to", input1)
                return true
            }
        }
    }
    return false
}

func checkIncludeDir(clidriver string) bool {
    if dirExists(clidriver + "/include") {
        //fmt.Println("clidriver/include folder exists in the path", clidriver)
        if dirExists(clidriver + "/lib") {
            //fmt.Println("clidriver/lib folder exists in the path", clidriver)
            if fileExists(clidriver + "/include/sqlcli.h") {
                return true
            } else {
                fmt.Println(clidriver + "/include/sqlcli.h file does not exist.")
            }
        } else {
            fmt.Println(clidriver + "/lib directory does not exist.")
        }
    } else {
        fmt.Println(clidriver + "/include directory does not exist.")
    }
    return false
}

func main() {
    var target, ibmdbDir, cliFileName, existingDriver string
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
        // db2cli validate is not working, check for IBM_DB_HOME
        ibmDBHome, found := os.LookupEnv("IBM_DB_HOME")
        if !found  || ibmDBHome == "" || !dirExists(ibmDBHome) {
            // IBM_DB_HOME is not set
            if runtime.GOOS == "windows" {
                fmt.Println("Please set IBM_DB_HOME and add %IBM_DB_HOME%/bin to PATH and %IBM_DB_HOME%/lib to LIB environment variables after clidriver installation")
            } else if runtime.GOOS == "aix" {
                fmt.Println("Please set IBM_DB_HOME, CGO_CFLAGS, CGO_LDFLAGS and LIBPATH environment variables after clidriver installation")
            } else if runtime.GOOS == "darwin" {
                fmt.Println("Please set IBM_DB_HOME, CGO_CFLAGS, CGO_LDFLAGS and DYLD_LIBRARY_PATH environment variables after clidriver installation")
            } else {
                fmt.Println("Please set IBM_DB_HOME, CGO_CFLAGS, CGO_LDFLAGS and LD_LIBRARY_PATH environment variables after clidriver installation")
            }
        } else {
            fmt.Println("IBM_DB_HOME is set to", ibmDBHome)
        }
    } else {
        fmt.Println("db2cli validate command is working.")
        ibmDBHome, found := os.LookupEnv("IBM_DB_HOME")
        if !found  || ibmDBHome == "" || !dirExists(ibmDBHome) {
            // Check install path in validate output
            if checkInstalledDriver(string(out)) {
                os.Exit(1)
            }
        } else {
            fmt.Println("IBM_DB_HOME =", ibmDBHome)
            if checkIncludeDir(ibmDBHome) {
                fmt.Println("IBM_DB_HOME is set and directory exists. Skipping download.")
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
    existingDriver, _ = filepath.Abs(target + "/clidriver")

    if dirExists(existingDriver) {
        fmt.Println("Detected existing directory \"" + existingDriver + "\"")

        if checkIncludeDir(existingDriver) {
            fmt.Println("==> Skipping clidriver download.")
            os.Exit(2)
        } else {
            // Attempt to remove the directory and its contents
            err2 := os.RemoveAll(existingDriver)
            if err2 != nil {
                fmt.Printf("Error deleting directory '%s': %v\n", existingDriver, err2)
            } else {
                fmt.Printf("Directory '%s' and its contents successfully deleted.\n", existingDriver)
            }
            fmt.Println("==> Installing clidriver ...")
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
    } else if runtime.GOOS == "aix" {
        unpackageType = 3
        const wordsize = 32 << (^uint(0) >> 32 & 1)
        if wordsize == 64 {
            cliFileName = "aix64_odbc_cli.tar.gz"
        } else {
            cliFileName = "aix32_odbc_cli.tar.gz"
        }
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
    fileUrl := getDownloadUrl(cliFileName)
    fmt.Println("Downloading " + fileUrl)

    // Try primary
    err := downloadFile(fileUrl, cliFileName)
    if err != nil {
        fmt.Println("Error:", err)
        // Try fallback if IBM_DB_DOWNLOAD_URL is not set.
        downloadUrl, downloadUrlFound := os.LookupEnv("IBM_DB_DOWNLOAD_URL")
        if !downloadUrlFound || downloadUrl == "" {
            fmt.Println("Trying fallback URL...")
            newUrl := strings.ReplaceAll(fileUrl, primaryURL, fallbackURL)
            err = downloadFile(newUrl, cliFileName)
            if err != nil {
                fmt.Println("Fallback download failed:", err)
                fmt.Println("Error while downloading file: " + err.Error())
                os.Exit(4)
            }
        } else {
            os.Exit(5)
        }
    }
    fmt.Println("Download completed.")

    if unpackageType == 1 {
        Unzipping(cliFileName, target)
    } else if unpackageType == 3 {
        aix_untar(cliFileName, target)
    } else {
        linux_untar(cliFileName, target, ibmdbDir)
    }
}

func getDownloadUrl(cliFileName string) string {
    downloadUrl, downloadUrlFound := os.LookupEnv("IBM_DB_DOWNLOAD_URL")
    if !downloadUrlFound || downloadUrl == "" {
        clidriverVersion, ok := os.LookupEnv("CLIDRIVER_DOWNLOAD_VERSION")
        if ok && clidriverVersion != "" {
            if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" && strings.HasPrefix(strings.ToLower(clidriverVersion), "v11") {
                downloadUrl = primaryURL + cliFileName
            } else if runtime.GOOS == "darwin" && runtime.GOARCH != "arm64" && strings.HasPrefix(strings.ToLower(clidriverVersion), "v12") {
                downloadUrl = primaryURL + cliFileName
            } else {
                downloadUrl = primaryURL + clidriverVersion + "/" + cliFileName
            }
        } else {
            downloadUrl = primaryURL + cliFileName
        }
    } else {
        fmt.Println("IBM_DB_DOWNLOAD_URL is set to", downloadUrl)
    }
    return downloadUrl
}
