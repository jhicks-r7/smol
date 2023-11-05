package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"text/template"

	"github.com/go-playground/form"
	"github.com/julienschmidt/httprouter"
)

type Config struct {
    ListenerAddress string
    ListenerPort string
}

type CartItem struct {
    ItemNumber string
    Comment string
}

type TemplateData struct {
    CartItems []CartItem
}

type UpdateCommentForm struct {
    ItemNumber int `form:"itemnumber"`
    Comment string `form:"comment"`
}

type DeleteItemForm struct {
    ItemNumber int `form:"itemnumber"`
}

type Application struct {
    Config Config
    CartItems   []CartItem
    FormDecoder *form.Decoder 
}

func GenerateItemNumber() string {
    charset := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
    length := 8

    randStr := make([]rune, length)
    for i := range randStr {
        randStr[i] = charset[rand.Intn(len(charset))] 
    }
    return string(randStr)
}

func (conf *Config) getListener() string {
    return conf.ListenerAddress + ":" + conf.ListenerPort
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
    data := app.GetTemplateDefault(w, r)  

    ts, err := template.ParseFiles("./home.tmpl")
    
    if err != nil {
        fmt.Printf("Failed to parse home template: %v", err)
        os.Exit(1)
    }

    err = ts.Execute(w, data)
    
    if err != nil {
        fmt.Printf("Failed to execute home template: %v", err)
    }
    
}

func (app *Application) decodePostForm(r *http.Request, dst any) error {
    // Call ParseForm() on the request
    err := r.ParseForm()
    if err != nil {
        return err
    }
   
    // Call Decode() on our decoder isntance, passing the target destination as the first parameter
    err = app.FormDecoder.Decode(dst, r.PostForm)
    if err != nil {
        // If we try to use an invalid target destination, the Decode() method
        // Will return an error with the type *form.InvalidDecoderError.. We use
        // errors.As() to check for this and raise a panic rather than returning the error
        var invalidDecoderError *form.InvalidDecoderError

        if errors.As(err, &invalidDecoderError) {
            panic(err)
        }
        return err
    }

    return nil
}

func (app *Application) addItem(w http.ResponseWriter, r *http.Request) { 
    item := CartItem {
        ItemNumber: GenerateItemNumber(),
    } 

    app.CartItems = append(app.CartItems, item)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) DeleteItem(w http.ResponseWriter, r *http.Request) {
   var form DeleteItemForm

    err := app.decodePostForm(r, &form)
    if err != nil {
        app.BadRequestResponse(w, r)
        return
    } 

    if form.ItemNumber > len(app.CartItems) {
        app.BadRequestResponse(w, r)
        return
    }
    
    index := form.ItemNumber
    
    app.CartItems = append(app.CartItems[:index], app.CartItems[index+1:]...)
    
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) BadRequestResponse(w http.ResponseWriter, r *http.Request) {
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)   
}

func (app *Application) UpdateComment(w http.ResponseWriter, r *http.Request) {
    var form UpdateCommentForm

    err := app.decodePostForm(r, &form)
    if err != nil {
        app.BadRequestResponse(w, r)
        return
    }
    
    if form.ItemNumber > len(app.CartItems) {
        app.BadRequestResponse(w, r) 
        return
    }

    if len(form.Comment) > 20 {
        app.BadRequestResponse(w, r)
        return
    }
    
    app.CartItems[form.ItemNumber].Comment = form.Comment 
    http.Redirect(w, r, "/", http.StatusSeeOther) 
     
}

func (app *Application) GetTemplateDefault(w http.ResponseWriter, r *http.Request) TemplateData {
    data := TemplateData {
        CartItems: app.CartItems,
    }
    
    return data
}

func main() {
    var conf Config
    flag.StringVar(&conf.ListenerAddress, "address", "127.0.0.1", "IPv4 Address to Listen On")
    flag.StringVar(&conf.ListenerPort, "port", "8000", "TCP Port to Listen On")
    flag.Parse()
     
    var cartItems []CartItem
    app := Application {
        Config: conf,
        CartItems: cartItems,
        FormDecoder: form.NewDecoder(),
    }
    
    router := httprouter.New()

    router.HandlerFunc(http.MethodGet, "/", app.home)
    router.HandlerFunc(http.MethodGet, "/api/add", app.addItem)     
    router.HandlerFunc(http.MethodPost, "/api/updatecomment", app.UpdateComment)
    router.HandlerFunc(http.MethodPost, "/api/delete", app.DeleteItem)
    
    srv := http.Server {
        Addr: app.Config.getListener(),
        Handler: router, 
    }
    
    fmt.Printf("Starting server on %s\n", conf.getListener())
    err := srv.ListenAndServe()
    fmt.Printf("Error: %v", err)
}
