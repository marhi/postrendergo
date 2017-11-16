package main

import (
    "text/template"
    "io/ioutil"
    "log"
    "os"
    "flag"
    "net/http"
    "time"
    "github.com/gorilla/mux"
    "encoding/json"
    "sync"
    "strings"
    "path/filepath"
)

type Config struct {
    Delay         int64
    To            string
    Authorization string
    Output        string
    sync.Mutex
}

type Data struct {
    Content  string `json:"content"`
    Path     string `json:"path"`
    Fullpath string `json:"fullpath"`
}

var c = Config{}

// os.MkDirAll Not working?
func MkdirR(path []string) {

    for k, _ := range (path) {
        p, _ := filepath.Abs(filepath.Join(path[:k+1]...))
        log.Println("Creating: ", p)
        os.Mkdir(p, 0770)
    }

}

func WriteData(d *Data) {
    c.Lock()
    defer c.Unlock()

    // Insecure, fix
    path := strings.Replace(d.Fullpath, "..", "", -1)

    paths := strings.Split(path, "/")
    paths = append([]string{c.Output}, paths...)

    MkdirR(paths)

    outfile := filepath.Join(append(paths, "postrender")...)
    log.Println("File out: ", outfile)

    err := ioutil.WriteFile(outfile, []byte(d.Content), 0660)
    if err != nil {
        log.Println(err)
    }

}

func handlePost(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Access-Control-Allow-Origin", "*")
    defer req.Body.Close()

    decoder := json.NewDecoder(req.Body)
    d := Data{}
    err := decoder.Decode(&d)

    if err != nil {
        res.WriteHeader(http.StatusInternalServerError)
        return
    }

    go WriteData(&d)
    res.WriteHeader(http.StatusOK)
    return
}

func main() {

    delay := flag.Int64("delay", 2000, "Time in ms to wait before sending HTML to the backend")
    to := flag.String("host", "", "Url to which a POST request is made with the HTML, eg. http://127.0.20:1337/renderer")
    auth := flag.String("auth", "", "Auth string for Authorization header")
    listen := flag.String("listen", "", "Port on which to listen for incomming requests, ie: :1337")
    target := flag.String("dir", "", "Directory in which to save the rendered content")

    flag.Parse()
    if len(*to) <= 0 {
        flag.Usage()
        os.Exit(1)
    }

    s, err := os.Stat(*target)
    if err != nil || !s.IsDir() {
        if err != os.ErrNotExist {
            log.Println(err)
            os.Exit(2)
        }
    }

    log.Printf("Outputting to: %s\n", *target)

    // No need for mutex yet, but maybe we'll provide a live config change
    // And where's the fun without mutexes? Mutices? Those lock things.
    c.Lock()
    c.Delay = *delay
    c.Output = *target
    c.To = *to
    c.Authorization = *auth
    c.Unlock()

    tcontent, err := ioutil.ReadFile("template.txt")
    if err != nil {
        log.Fatalln(err)
    }

    templ, err := template.New("prerender.js").Parse(string(tcontent))
    if err != nil {
        log.Fatalln(err)
    }

    err = templ.Execute(os.Stdout, c)
    if err != nil {
        log.Fatalln(err)
    }

    if len(*listen) > 0 {

        r := mux.NewRouter()
        r.HandleFunc("/renderer", handlePost)

        server := http.Server{
            IdleTimeout:  time.Duration(*delay) * time.Millisecond,
            WriteTimeout: time.Duration(*delay) * time.Millisecond,
            ReadTimeout:  time.Duration(*delay) * time.Millisecond,
            Addr:         *listen,
            Handler:      r,
        }

        log.Fatalln(server.ListenAndServe())
    }

}
