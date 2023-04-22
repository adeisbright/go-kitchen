package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

const (
	HOST = "127.0.0.1"
	PORT = "8080"
)

//Todo Interface
type Todo struct {
	Done  bool
	Title string
}

type App struct {
	Title      string
	Text       string
	Tagline    string
	Activities []Todo
}

// indexHandler is the route handler for a request that hits the home page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate, _ := template.ParseFiles("./templates/index.html")

	menu := App{
		Title:   "Go Kitchen",
		Tagline: "Go Serve Yourself",
		Text: `Go Kitchen is a Serve Yourself Kitchen. You have no
			human to interact with.
			Your entire communication is with the automated systems`,
		Activities: []Todo{
			{
				Done:  true,
				Title: "Build a Go Microservice",
			},
			{
				Done:  false,
				Title: "Encrypt User Data",
			},
			{
				Done:  true,
				Title: "Create a Token in the Metaverse",
			},
		},
	}
	indexTemplate.Execute(w, menu)
}

type User struct {
	Email           string
	Password        string
	ConfirmPassword string
}

type LoginData struct {
	Title string
}

// Handles user registration once they submit their details

func handleRegistration(w http.ResponseWriter, r *http.Request) {

	//Parse the form
	r.ParseForm()
	for key, value := range r.Form {
		fmt.Printf("key %v value %v \n", key, value)
		break
	}

	email := template.HTMLEscapeString(r.Form.Get("email"))
	password := r.Form["password"][0]
	cPassword := r.Form["cpassword"][0]
	fmt.Printf("%v, %v, %v", email, password, cPassword)

	//Handle the JSON and respond
	user := User{
		Email:           email,
		Password:        password,
		ConfirmPassword: cPassword,
	}

	json, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

type ValidUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Error Handling Code
func ErrorHandler(w http.ResponseWriter, message string, statusCode int) {
	// Set the Content type and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Create a map that will later have a message key
	resp := make(map[string]string)

	resp["message"] = message

	// Stringnify or marshal the map
	json, _ := json.Marshal(resp)

	w.Write(json)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		//Retrieve the content type and ensure it is only application/json
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			ErrorHandler(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}

		// Create a reference to the type that is to be used
		user := ValidUser{}

		// Create an Error Variable for unsupported json type
		var unMarshalError *json.UnmarshalTypeError

		// Parse the body of the request
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&user)
		if err != nil {
			if errors.As(err, &unMarshalError) {
				ErrorHandler(w, "Bad Request: Wrong Data provided for"+unMarshalError.Field, http.StatusBadRequest)
			} else {
				ErrorHandler(w, "Bad Request"+err.Error(), http.StatusBadRequest)
			}
			return
		}

		json, _ := json.Marshal(user)
		w.Write(json)

	} else if r.Method == "GET" {
		tmp := template.Must(template.ParseFiles("./templates/login.html"))
		data := LoginData{
			Title: "Login",
		}
		tmp.Execute(w, data)
	}

}

func handleOverviewVisit(w http.ResponseWriter, r *http.Request) {
	overviewTemplate, _ := template.ParseFiles("./templates/overview.html")
	type Event struct {
		Title       string
		Description string
		Address     string
		Date        string
		ToAttend    bool
	}

	type TemplateData struct {
		Title  string
		Events []Event
	}

	data := TemplateData{
		Title: "Welcome",
		Events: []Event{
			{
				Title:       "Bigjara LogN Community Meetup",
				Description: "A meeting for those interested in Bigjara",
				Address:     "KM 25 Lekki Epe Expressway, Dominion Plaza",
				Date:        "April 29th 2023, 10AM",
				ToAttend:    false,
			},
			{
				Title:       "Bigjara Lite BootCamp",
				Description: "This is a two weeks training on Machine Learning",
				Address:     "KM 25 Lekki Epe Expressway, Dominion Plaza",
				Date:        "May 6th 2023, 10AM",
				ToAttend:    true,
			},
		},
	}

	if r.Method == "GET" {
		overviewTemplate.Execute(w, data)
	} else if r.Method == "POST" {
		r.ParseForm()
		for key, value := range r.Form {
			fmt.Printf("key %v value %v \n", key, value)
			//çbreak
		}

		title := template.HTMLEscapeString(r.Form.Get("title"))
		description := r.PostFormValue("description")
		data := TemplateData{
			Title: "Form Submitted",
			Events: []Event{
				{
					Title:       title,
					Address:     "Empty",
					Description: description,
					Date:        "",
					ToAttend:    true,
				},
			},
		}
		//overviewTemplate.Execute(w, data)

		json, _ := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Token", "123456")
		w.WriteHeader(http.StatusOK) //Sends a HTTP respnose header with the given  code as the provided one
		w.Write(json)
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	hasUserName := r.URL.Query().Has("name")
	userName := r.URL.Query().Get("name")
	if hasUserName {
		fmt.Printf("%v , %v", "The user name is ", userName)

	} else {
		fmt.Printf("%v , %v", "This visitor has no name. Then ", "anonymous")
	}

	fmt.Fprintf(w, "You were the f**ked dude he warned us about")
}

func main() {
	fmt.Println("The server is now listening on port 8080")
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/overview", handleOverviewVisit)
	http.HandleFunc("/register", handleRegistration)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/test", testHandler)
	http.ListenAndServe(HOST+":"+PORT, nil)

}
