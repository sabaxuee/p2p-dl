package main

import (
    "fmt"
    "net/http"
    //"strings"
    "log"
    "strconv"
    "io"
    "io/ioutil"
    "os"
)

type p2pMux struct{
}

/*Router Handler*/
func (p *p2pMux)ServeHTTP(w http.ResponseWriter, r *http.Request){

    fmt.Println(r.Method,r.URL.Path)

    switch r.URL.Path{
    case "/":
        sayhello(w,r)
    case "/pull":
        pull(w,r)
    case "/p2p":
        p2p(w,r)
    default:
        //http.ServeFile(w,r,r.URL.Path[1:])
        fmt.Fprintf(w, "Hello you are visiting path %s, but it doesn't exist.\n", r.URL.Path)
    }
}

/*pull parses filename and target IP from HTTP GET method, and start downloading routine. */
func pull(w http.ResponseWriter, r *http.Request){
    if r.Method == "GET"{
        r.ParseForm()
        src := r.Form.Get("src")
        file := r.Form.Get("f")
        if len(src) > 0 && len(file) > 0 {
            uri := "http://" + r.Form.Get("src") + "/p2p?file=" + r.Form.Get("f")
            fmt.Fprintln(w, uri)
            //n,err := download(src, f)
            go dl(src, file)
        }else{
            fmt.Fprintf(w,"invalid arguments.")
            return
        }
        fmt.Fprintf(w,"file will be download.")
    }
}

func dl(src, f string){
    n, err := download(src, f)
    if err != nil{
        fmt.Println("[%d bytes returned.]",n ,err)
    }
}

/*download routine, supports resuming broken downloads.*/
func download(src , file string)(int64, error){
    url := "http://" + src + "/p2p?file=" + file
    fmt.Printf("we are going to download %s\n", url)
    out, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        return 0, err
    }
    defer out.Close()

    stat, err := out.Stat()
    if err != nil {
        return 0, err
    }
    out.Seek(stat.Size(), 0)

    req, err := http.NewRequest("GET", url, nil)
    req.Header.Set("User-Agent","go-downloader")
    /* Set download starting position with 'Range' in HTTP header*/
    req.Header.Set("Range", "bytes=" + strconv.FormatInt(stat.Size(), 10) + "-")
    fmt.Printf("%v bytes had already been downloaded.\n", stat.Size())

    resp, err := http.DefaultClient.Do(req)

    /*Save response body to file only when HTTP 2xx received. TODO*/
    if err != nil || (resp != nil && resp.StatusCode / 100 != 2) {
        if resp != nil{
            fmt.Println("http status code:", resp.StatusCode, err)
            body, _ := ioutil.ReadAll(resp.Body)
            fmt.Println("response Body:", string(body))
        }
        return 0, err
    }
    defer resp.Body.Close()

    n, err := io.Copy(out, resp.Body)
    if err != nil {
        return 0, err
    }
    fmt.Printf("%d bytes downloaded.",n)
    return n, nil
}

/* p2p responses downloading request, large file blocking implentmented by http.ServeFile()*/
func p2p(w http.ResponseWriter, r *http.Request){
    fmt.Println("p2p...")
    r.ParseForm()
    file := r.Form.Get("file")
    http.ServeFile(w,r,"p2p/"+ file)
}

/*web main page*/
func sayhello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello world!")
}


func main() {
    mux := &p2pMux{}
    /*
    http.HandleFunc("/", sayhello)
    http.HandleFunc("/pull", pull)
    http.NotFoundHandler()
    */

    port := "0.0.0.0:9090"

    /*
    if len(os.Args) == 2{
        port = os.Args[1] + port
    }else if len(os.Args) == 3{
        port = os.Args[1] + ":" + os.Args[2]
    }
    */

    fmt.Printf("listening on %s\n", port)
    err := http.ListenAndServe(port, mux)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
    fmt.Println("Hello, world!")
}
