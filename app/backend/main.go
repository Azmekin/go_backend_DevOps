package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"html/template"
	"log"
	"net/http"
	"os"
)

var Ctx = context.Background()

func ConnectRedis() *redis.Client {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
	redis_addr, _ := os.LookupEnv("REDIS_ADDR")
	username, _ := os.LookupEnv("USERNAME")
	password, _ := os.LookupEnv("PASSWORD")
	cert, err := tls.LoadX509KeyPair("redis.crt", "redis.key")
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Username: username, // use your Redis user. More info https://redis.io/docs/management/security/acl/
		Password: password,
		TLSConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		},
	})
	return client

}

var Cli = ConnectRedis()

type Page struct {
	Title string
	Key   []byte
	Value []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Key, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Key: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	//p, err := loadPage(title) //сменить на просмотр параметра
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusForbidden)
	//	return
	//}
	querySting := query["key"][0]
	val, err := Cli.Get(Ctx, querySting).Result()
	if err != nil {
		http.Error(w, "404", http.StatusNotFound)
		return
	}
	p := &Page{Title: "View", Key: []byte(querySting), Value: []byte(val)}
	//Cli.Get(Ctx, querySting)
	renderTemplate(w, "view", p)
}
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title} //сменить на создание параметра
	}
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, p)
}

func killHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadPage("Delete")
	if err != nil {
		p = &Page{Title: "Delete"} //сменить на создание параметра
	}
	t, _ := template.ParseFiles("del.html")
	t.Execute(w, p)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "403", http.StatusForbidden)
	return
}

func delHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	Cli.GetDel(Ctx, key) //Work only with string command from task link
	http.Redirect(w, r, "/get_key?"+"key"+"="+key, http.StatusFound)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[len("/save/"):]
	key := r.FormValue("key")
	value := r.FormValue("value")
	//p := &Page{Title: "Save", Key: []byte(body), Value: []byte(body)}
	err := Cli.Set(Ctx, key, value, 0).Err()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/get_key?"+"key"+"="+key, http.StatusFound)
}

func main() {

	http.HandleFunc("/get_key", viewHandler)
	http.HandleFunc("/save_key", editHandler) //form for send request on set_key
	http.HandleFunc("/kill_key", killHandler) //form for send request on del_key
	http.HandleFunc("/set_key", saveHandler)
	http.HandleFunc("/del_key", delHandler)
	http.HandleFunc("/", startHandler) // url like "/payload/: doesn't work, only "/payload" like task form
	log.Fatal(http.ListenAndServe(":8080", nil))
}
