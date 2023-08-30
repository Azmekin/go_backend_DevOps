package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/buger/jsonparser"
	"html/template"
	"log"
	"net/http"
	"os"
	"io"
	"fmt"
)

var Ctx = context.Background()

func ConnectRedis() *redis.Client {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
	redis_addr, _ := os.LookupEnv("REDIS_ADDR")
	username, _ := os.LookupEnv("USERNAME")
	password, _ := os.LookupEnv("PASSWORD")
	cert, err := tls.LoadX509KeyPair("cert/redis.crt", "cert/redis.key")
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := os.ReadFile("cert/ca.crt")
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
	if len(query["key"]) < 1 {
		http.Error(w, "404", http.StatusNotFound)
		return
	}
	querySting := query["key"][0]
	val, err := Cli.Get(Ctx, querySting).Result()
	if err != nil {
		http.Error(w, "404, maybe you deleted it?", http.StatusNotFound)
		return
	}
	p := &Page{Title: "View", Key: []byte(querySting), Value: []byte(val)}
	//Cli.Get(Ctx, querySting)
	renderTemplate(w, "view", p)
}




func startHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "403", http.StatusForbidden)
	return
}

func delHandler(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
    	http.Error(w, err.Error(), 422)
    	return
    }
    value,_,_,err := jsonparser.Get(body, "key")
    if err != nil {
    	http.Error(w, err.Error(), 422)
    	return
    }
	Cli.GetDel(Ctx, string(value)) //It works only with string command from task link. It's because there is not Cli.Del
	http.Redirect(w, r, "/get_key?"+"key"+"="+string(value), http.StatusFound)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	data,err := io.ReadAll(r.Body)
	if err != nil {
    	http.Error(w, err.Error(), 422)
    	return
    }
	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
	    err := Cli.Set(Ctx, string(key), string(value), 0).Err()
	    if err != nil {
        		http.Error(w, err.Error(), 422)
        		return nil
        	}
    	return nil
    })

	fmt.Fprintf(w, "Accepted")
}

func main() {

	http.HandleFunc("/get_key", viewHandler)
	http.HandleFunc("/set_key", saveHandler)
	http.HandleFunc("/del_key", delHandler)
	http.HandleFunc("/", startHandler) // url like "/payload/: doesn't work, only "/payload" like task form
	log.Fatal(http.ListenAndServe(":8080", nil))
}
