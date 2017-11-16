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
)

type Config struct {
    Delay         int64
    To            string
    Authorization string
}

func handlePost(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Access-Control-Allow-Origin", "*")
    res.WriteHeader(http.StatusOK)
    log.Println(req)
    return
}

func main() {

    delay := flag.Int64("delay", 2000, "Time in ms to wait before sending HTML to the backend")
    to := flag.String("host", "", "Url to which a POST request is made with the HTML, eg. http://127.0.20:1337/renderer")
    auth := flag.String("auth", "", "Auth string for Authorization header")
    listen := flag.String("listen", "", "Port on which to listen for incomming requests, ie: :1337")

    flag.Parse()
    if len(*to) <= 0 {
        flag.Usage()
        os.Exit(1)
    }

    c := Config{*delay, *to, *auth}

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
