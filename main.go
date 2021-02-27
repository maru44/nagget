package main

import (
    "net/http"
    "log"
    "fmt"
    "html/template"
    "time"
    "io/ioutil"
    //"os"
    "strconv"
    "encoding/json"
    //"github.com/gorilla/csrf"
    "github.com/gorilla/mux"
    //"github.com/jinzhu/gorm"
    //_ "github.com/mattn/go-sqlite3"

    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

var gen_temp = make(map[string]*template.Template)


type Article struct {
    //gorm.Model
    Id string `json:"Id"`
    Kind string `json:"Kind"`
    Content string `json:"Content"`
    Gender string `json:"Gender"`
}

type Blog struct {
    gorm.Model
    //Id string `json:"Id"`
    Kind string `json:"Kind"`
    Content string `json:"Content"`
    Gender string `json:"Gender"`
    Good int `json:"Good"`
}

type Cookie struct {
    Name string
    Value string

    Path string
    Domain string
    Expires time.Time
    RawExpires string

    MaxAge int
    Secure bool
    HttpOnly bool
    Raw string
    Unparsed []string
}

/*
var blog Blog
var blogs []Blog
*/

var Articles []Article

func main() {
    gen_temp["home"] = loadTemplate("home")
    gen_temp["confirm"] = loadTemplate("confirm")
    gen_temp["artie"] = loadTemplate("artie")
    gen_temp["list"] = loadTemplate("list")

    //db, err := gorm.Open("sqlite3", "db.sqlite3")
    db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
    if err != nil {
        panic("falied")
    }
    //defer db.Close()

    //Migrate
    db.AutoMigrate(&Blog{})

    //初回のみ
    //create
    db.Create(&Blog{Model: gorm.Model {ID: 1}, Kind: "困ったこと", Content: "困ったな", Gender: "男性"})
    db.Create(&Blog{Model: gorm.Model {ID: 2}, Kind: "不便", Content: "不便だった", Gender: "女性"})
    db.Create(&Blog{Model: gorm.Model {ID: 3}, Kind: "不便", Content: "不便だった2", Gender: "女性"})
    p := Blog{Model: gorm.Model {ID: 4}, Kind: "不便", Content: "不便だった3", Gender: "女性"}
    db.Create(&p)
    //db.Create(&Blog{Model: gorm.Model {ID: 4}, Kind: "不便", Content: "不便だった3", Gender: "女性"})
    db.Create(&Blog{Model: gorm.Model {ID: 5}, Kind: "あったらいいな", Content: "欲しい", Gender: "男性"})
    
    /*
    db.Create(&Blog{Model: gorm.Model {ID}})
    if db.NewRecord(p) {
        db.Create(&p)
    }
    */

    Articles = []Article{
        Article{Id: "1", Kind: "困ったこと", Content: "困ったな", Gender: "男性"},
        Article{Id: "2", Kind: "不便", Content: "不便だった", Gender: "女性"},
        Article{Id: "3", Kind: "不便", Content: "不便だった2", Gender: "女性"},
        Article{Id: "4", Kind: "不便", Content: "不便だった3", Gender: "女性"},
        Article{Id: "5", Kind: "あったらいいな", Content: "欲しい", Gender: "男性"},
    }

    // /static/
    //dir, _ := os.Getwd()
    //http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/static/"))))

    r := mux.NewRouter()

    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

    //  "/"  >>  home
    r.HandleFunc("/", home)

    r.HandleFunc("/create", createBlog).Methods("POST")

    r.HandleFunc("/blogs", blogs)

    // iina
    r.HandleFunc("/iina", iinas)

    // komatta
    r.HandleFunc("/komatta", komas)

    // huben
    r.HandleFunc("/huben", hubens)

    r.StrictSlash(true).HandleFunc("/api/detail/{id}", detailApi).Methods("GET")

    // localhost:8000
    //一番下に書く
    //csrf
      //本番環境ではcsrf.Secure(false)が要らない
      //h := csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false))(r)
      //http.ListenAndServe(":8000", h)
    //http.ListenAndServe(":8000", nil)
    http.ListenAndServe(":8000", r)
}

func blogAll() []Blog {
    db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
    if err != nil {
        panic("falied")
    }
    var blogs []Blog
    //db.Limit(3).Order("created_at desc").Find(&blogs)
    db.Order("created_at desc").Find(&blogs)
    return blogs
}

func blogFilter(k string) []Blog {
    db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
    if err != nil {
        panic("failed")
    }
    var blogs []Blog
    db.Where("kind = ?", k).Order("created_at desc").Find(&blogs)
    return blogs
}

func handlerA(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello World from Go.")
}

func home(w http.ResponseWriter, r *http.Request) {
    //w.Header().Set("Set-Cookie", "_gorilla_csrf="+csrf.Token(r))
    if err := gen_temp["home"].Execute(w, struct {
        Header header
        Title string
    }{
        Header : newHeader("ホーム"),// {{ template "header" .Header }}での.Header
        Title: "home",
    }); err != nil {
        log.Printf("failed to execute template: %v", err)
    }
}

func confirm(w http.ResponseWriter, r *http.Request) {
    //kind := r.Form.Get("kinds")
    if err := gen_temp["confirm"].Execute(w, struct {
        Header header
        Kind string
        Gender string
        Content string
    }{
        Header: newHeader("確認画面"),
        Kind: r.FormValue("kinds"),
        Content: r.FormValue("content"),
        Gender: r.FormValue("genderSelect"),
    }); err != nil {
        log.Printf("failed to execute template: %v", err)
    }
}

func artis(w http.ResponseWriter, r *http.Request) {
    //json.NewEncoder(w).Encode(articles)
    if err := gen_temp["artie"].Execute(w, struct {
        Items []Article
        Header header
    }{
        Items: Articles,
        Header: newHeader("記事一覧"),
    }); err != nil {
        log.Printf("failed to execute template: %v", err)
    }
}

func blogs(w http.ResponseWriter, r *http.Request) {
    if err := gen_temp["artie"].Execute(w, struct {
        Header header
        Items []Blog
    }{
        Header: newHeader("Blogs"),
        Items: blogAll(),
    }); err != nil {
        log.Printf("failed to execute template: %v", err)
    }
}

func iinas(w http.ResponseWriter, r *http.Request) {
    if err := gen_temp["list"].Execute(w, struct {
        Header header
        Items []Blog
    }{
        Header: newHeader("あったらいいな"),
        Items: blogFilter("あったらいいな"),
    }); err != nil {
        log.Printf("failed to execute template: %v", err)
    }
}

func komas(w http.ResponseWriter, r *http.Request) {
    if err := gen_temp["list"].Execute(w, struct {
        Header header
        Items []Blog
    }{
        Header: newHeader("困ったこと"),
        Items: blogFilter("困ったこと"),
    }); err != nil {
        log.Printf("failed to execute template: %v", err)
    }
}

func hubens(w http.ResponseWriter, r *http.Request) {
    if err := gen_temp["list"].Execute(w, struct {
        Header header
        Items []Blog
    }{
        Header: newHeader("不便"),
        Items: blogFilter("不便"),
    }); err != nil {
        log.Printf("failed to execute template: %v", err)
    }
}

func detailApi(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
    if err != nil {
        panic("falied")
    }

    blogs := []Blog{}
    db.Find(&blogs)
    for _, blog := range blogs {
        //string to int
        s, err := strconv.Atoi(id)
        cookieStr := "good_" + id
        ss := uint(s)
        if err == nil {
            if blog.ID == ss {
                //if r.Cookie(cookieStr) {}
                coo, ok_ := r.Cookie(cookieStr)
                if ok_ == nil {
                    if coo.Value == "1" {
                        blog.Good -= 1
                        db.Save(&blog)
                        //newGood := blog.Good - 1
                        //db.Model(&blog).Update("Good", newGood)
                        cookie := &http.Cookie{
                            Name: cookieStr,
                            Value: "0",
                            Path: "/",
                        }
                        http.SetCookie(w, cookie)
                    } else {
                        blog.Good += 1
                        db.Save(&blog)
                        cookie := &http.Cookie{
                            Name: cookieStr,
                            Value: "1",
                            Path: "/",
                        }
                        http.SetCookie(w, cookie)
                    }
                } else {
                    blog.Good += 1
                    db.Save(&blog)
                    cookie := &http.Cookie{
                        Name: cookieStr,
                        Value: "1",
                        Path: "/",
                    }
                    http.SetCookie(w, cookie)
                }
                json.NewEncoder(w).Encode(blog)
            }
        }
    }
}


func createBlog(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var blog Blog
    if err:= json.Unmarshal(reqBody, &blog); err!= nil {
        log.Fatal(err)
    }
    db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
    db.Create(&blog)
    responseBody, err := json.Marshal(blog)
    if err != nil {
        log.Fatal(err)
    }
    w.Write(responseBody)
}


func loadTemplate(name string) *template.Template {
    temp, err := template.ParseFiles(
        "templates/" + name + ".html",
        "templates/_header.html",
        "templates/_footer.html",
    )
    if err != nil {
        log.Fatalf("template error: %v", err)
    }
    return temp
}

type header struct {
    Pagetitle string
}

func newHeader(pagetitle string) header {
    return header{Pagetitle: pagetitle}
}