//-*- mode: go -*-

package main

import (
    "flag"
    "fmt"
    "gopkg.in/gcfg.v1"
    "log"
    "net/url"
    "os"
    "os/user"
    "path"
    "strings"
    "syscall"
    "os/exec"
    "io/ioutil"
)

type Options struct {
    Runtime struct {
        Mountpoint string
    }
}

type AppContext struct {
    Application string
    File string
    FullPath string
    Project string
    Uuid string
}

var configFile string
var uri string

var Logger *log.Logger;

func init(){
    flag.StringVar(&configFile, "file", "/etc/threepio.ini", "ConfigFile file for threepio")
    flag.StringVar(&configFile, "f", "/etc/threepio.ini", "ConfigFile file for threepio (Shorthand)")

    flag.StringVar(&uri, "uri", "threepio+prelude:///some/path?id=12345", "Project URI; see docs")
    flag.StringVar(&uri, "u", "threepio+prelude:///some/path?id=12345", "Project URI; see docs (shorthand)")

    usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }

    file, _ := os.OpenFile(path.Join(usr.HomeDir, ".threepio.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    Logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func getMountPointFromOptionsFile() (string){
    var options Options
    err := gcfg.ReadFileInto(&options, configFile)
    if err != nil {
        log.Fatalf("Failed to parse gcfg data: %s", err)
    }

    return options.Runtime.Mountpoint
}


func parseUri(uri string) (AppContext) {
    var appContextTest AppContext

    urlObj, err := url.Parse(uri)
    if err != nil {
        log.Fatal( err )
    }

    queryObj, err := url.ParseQuery(urlObj.RawQuery)
    if err != nil {
        log.Fatal( err )
    }

    schemeSplit := strings.Split(urlObj.Scheme, "+")

    appContextTest.Application = schemeSplit[len(schemeSplit)-1]
    appContextTest.Project = mutate(urlObj.Path)
    appContextTest.Uuid = queryObj.Get("uuid")

    return appContextTest
}

func createDirIfMissing(fullPath string){
    err := os.MkdirAll(fullPath, 0755)
    if err != nil {
        Logger.Fatal( err )
    }
}

func launch(appContext AppContext)(AppContext){
    binary, lookErr := exec.LookPath("open")
    if lookErr != nil {
        Logger.Fatal(lookErr)
    }

    args := []string{"open", path.Join(appContext.FullPath, appContext.File)}

    env := os.Environ()

    execErr := syscall.Exec(binary, args, env)
    if execErr != nil {
        Logger.Fatal(execErr)
    }
    return appContext
}

func mutate(s string) (s_mux string) {
    s_mux = strings.Replace(s, "/", "", 1)
    s_mux = strings.Replace(s_mux, " ", "_", -1)

    return
}

func inferFilename(appContext AppContext)(AppContext) {
    var suffix string

    switch appContext.Application {
    case "prelude": suffix = "plproj"
    case "premiere": suffix = "prproj"
    }

    appContext.File = fmt.Sprintf("%s.%s", appContext.Project, suffix)
    return appContext
}

func getAppContext(uri string, mount string)(AppContext) {
    var appContext = parseUri(uri)
    appContext = inferFilename(appContext)
    appContext.FullPath = path.Join(mount, appContext.Uuid)
    return appContext
}


func createAbobeProjectFilesIfNotExist(path string, fileContents []byte)(error error) {
    if _, error = os.Stat(path); error != nil {
        if os.IsNotExist(error) {
            error = ioutil.WriteFile(path, fileContents, 0644)
        }
    }
    return
}


func main(){
    flag.Parse()

    var mount = getMountPointFromOptionsFile()
    var appContext = getAppContext(uri, mount)

    Logger.Printf("Launching %s on path %s to edit %s with assets from project %s", appContext.Application, appContext.File, appContext.File , appContext.Uuid)

    createDirIfMissing(path.Join(mount, appContext.Uuid))

    if file, error := ioutil.ReadFile("/Applications/Threepio.app/Contents/MacOs/empty.plproj"); error == nil {
        createAbobeProjectFilesIfNotExist(path.Join(appContext.FullPath, appContext.Project)+".plproj",file)
    }
    if file, error := ioutil.ReadFile("/Applications/Threepio.app/Contents/MacOs/empty.prproj"); error == nil {
        createAbobeProjectFilesIfNotExist(path.Join(appContext.FullPath, appContext.Project) + ".prproj", file)
    }

    launch(appContext)
}
