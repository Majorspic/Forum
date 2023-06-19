package Forum

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"text/template"

	"github.com/markbates/goth/providers/google"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/gorilla/sessions"
	_ "github.com/markbates/goth"
	_ "github.com/markbates/goth/providers/google"

	"github.com/mattn/go-sqlite3"
)

type PostForm struct {
	Id          int
	IdParent    int
	Username    string
	Category    int
	Title       string
	Description string
	Date        string
	nbtoxic     string
}

const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05"
)

type Game struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func simulateError() {
	// Code qui provoque une erreur
	panic("Une erreur s'est produite")
}

func Home(w http.ResponseWriter, r *http.Request) {
	// hangmandata = extractJson()
	t, err := template.ParseFiles("./page/Home.html", "./template/header.html", "./template/popup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
	e := sqlite3.SQLITE_LIMIT_LENGTH
	fmt.Println(e)
}

func Md5Hash(input string) string {
	// Convertir la chaîne d'entrée en un tableau de bytes
	byteString := []byte(input)

	// Générer le hash md5 à partir du tableau de bytes
	hash := md5.Sum(byteString)

	// Convertir le hash en une chaîne de caractères hexadécimale
	hashString := hex.EncodeToString(hash[:])

	return hashString
}

func Register(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./page/Register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
	fmt.Println(r.Method)
}

func PostPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./page/post.html", "./template/header.html", "./template/comment.html")
	UrlIdPost := r.FormValue("id")
	var dataPost PostForm
	if UrlIdPost != "" {
		// Recuperer les données du post en bdd
		dataPost = getPostData(UrlIdPost)
		fmt.Println(dataPost)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, dataPost)
}

func ProfilPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./page/profilPage.html", "./template/header.html", "./template/likedPost.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./page/createPost.html", "./template/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func MyData(w http.ResponseWriter, r *http.Request) {
	// Vérification de la méthode HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Lecture du corps de la requête
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Affichage des données reçues
	fmt.Println("Request Body:", string(body))

	// Réponse de succès
	var data struct {
		Mail     string `json:"mail"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err != nil {
		fmt.Println("Erreur lors de l'encryotage du Mdp")
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	hashedPassword := Md5Hash(data.Password)

	fmt.Println("Data password", data.Password, "Mot de passe hash", hashedPassword)

	data.Password = string(hashedPassword)

	fmt.Println("Mot de passe normal", data.Password, "mot de passe encryptée", hashedPassword)

	fmt.Println(data.Username, data.Password, data.Mail)
	db, err := sql.Open("sqlite3", "./BDD/BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableQuery := `
			CREATE TABLE IF NOT EXISTS User (
				ID INTEGER PRIMARY KEY AUTOINCREMENT,
				Mail TEXT,
				Username TEXT,
				Password TEXT
			)
		`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	query := "SELECT COUNT(*) FROM User WHERE Mail = ?"
	var count int
	err = db.QueryRow(query, data.Mail).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		errorMessage := "Email déjà utilisé"
		errorResponse := ErrorResponse{Message: errorMessage}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	query2 := "SELECT COUNT(*) FROM User WHERE Username = ?"
	var count2 int
	err = db.QueryRow(query2, data.Username).Scan(&count2)
	if err != nil {
		log.Fatal(err)
	}
	if count2 > 0 {
		errorMessage := "Le pseudo est déjà pris"
		errorResponse := ErrorResponse{Message: errorMessage}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Création de la table "User" si elle n'existe pas déjà

	// Insertion des données dans la table "User"
	insertQuery := `
			INSERT INTO User (Mail, Username, Password)
			VALUES (?, ?, ?)
		`
	_, err = db.Exec(insertQuery, data.Mail, data.Username, data.Password)
	if err != nil {
		log.Fatal(err)
	}

	// Réponse HTTP réussie
	cookie := http.Cookie{
		Name:  "auth",
		Value: data.Password,
		Path:  "/",
	}

	http.SetCookie(w, &cookie)

	// Redirect to the profile page
	http.Redirect(w, r, "/profil", http.StatusFound)

}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("0215394687501684")
	store = sessions.NewCookieStore(key)
)

func Secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

func Login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

var (
	googleOauthConfig *oauth2.Config
)

const (
	googleClientID     = "551950133746-vt5v153ch2p8c55pu9dktsqpt4i4arvq.apps.googleusercontent.com"
	googleClientSecret = "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"
	googleRedirectURL  = "http://localhost:8080/callback"
)

// NewRouter crée et configure les routes
func NewRouter() http.Handler {
	router := http.NewServeMux()

	googleOauthConfig = &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  googleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	// Route racine
	router.HandleFunc("/Handle", HandleMain)

	// Route de connexion Google
	router.HandleFunc("/Hlogin", HandleGoogleLogin)

	// Route de rappel Google
	router.HandleFunc("/callback", HandleGoogleCallback)

	return router
}

func HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Go to <a href=\"/Hlogin\">login with Google</a>")
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {

	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, "Une erreur c'est produit lors de l'authentification", http.StatusInternalServerError)
		return
	}

	client := googleOauthConfig.Client(oauth2.NoContext, token)

	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	userInfo := struct {
		Email string `json:"email"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	rand.Seed(time.Now().UnixNano())

	username := "Guest" + strconv.Itoa(rand.Intn(9999-1000)+1000)
	email := userInfo.Email
	password := token.AccessToken

	fmt.Println("Nom d'utilisateur          ", username, "MDP     ", password, "email            ", email)
	db, err := sql.Open("sqlite3", "./BDD/BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS User (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Mail TEXT,
			Username TEXT,
			Password TEXT
		)
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	query := "SELECT COUNT(*) FROM User WHERE Mail = ?"
	var count int
	err = db.QueryRow(query, email).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		// Compte deja crée
		fmt.Println("Email deja utilisée")

		cookie := http.Cookie{
			Name:  "auth",
			Value: password,
			Path:  "/",
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/profil", http.StatusFound)

		return
	}

	query2 := "SELECT COUNT(*) FROM User WHERE Username = ?"
	var count2 int
	err = db.QueryRow(query2, username).Scan(&count2)
	if err != nil {
		log.Fatal(err)
	}
	if count2 > 0 {
		fmt.Println("Le Pseudo est deja pris !")
		fmt.Fprintf(w, "Pseudo Deja utilisée !")
		return
	}

	insertQuery := `
		INSERT INTO User (Mail, Username, Password)
		VALUES (?, ?, ?)
	`
	_, err = db.Exec(insertQuery, email, username, password)
	if err != nil {
		log.Fatal(err)
	}

	// Utilisez l'e-mail de l'utilisateur pour effectuer d'autres opérations, par exemple, vérifier si l'utilisateur est déjà enregistré sur votre site.
	cookie := http.Cookie{
		Name:  "auth",
		Value: password,
		Path:  "/",
	}

	http.SetCookie(w, &cookie)

	// Redirect to the profile page
	http.Redirect(w, r, "/profil", http.StatusFound)

}

var (
	clientID     = "5244ff87cf424b64bbb8"
	clientSecret = "b3f718c5af1c1ceb05256d667755b3bc23f3ef91"
	redirectURL  = "http://localhost:8080/github/callback"
)

var oauthConf = &oauth2.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	RedirectURL:  redirectURL,
	Scopes:       []string{"user:email"}, // Inclure la portée user:email
	Endpoint:     github.Endpoint,
}

func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprintf(w, "Échec de l'échange du code d'autorisation : %v", err)
		return
	}

	client := oauthConf.Client(context.Background(), token)

	// Obtenir les informations de l'utilisateur connecté
	userInfoURL := "https://api.github.com/user"
	response, err := client.Get(userInfoURL)
	if err != nil {
		fmt.Fprintf(w, "Échec de la requête pour obtenir les informations de l'utilisateur : %v", err)
		return
	}
	defer response.Body.Close()

	var userData struct {
		Email string `json:"email"`
		Login string `json:"login"`
	}

	err = json.NewDecoder(response.Body).Decode(&userData)
	if err != nil {
		fmt.Fprintf(w, "Échec de la lecture des données utilisateur : %v", err)
		return
	}

	// Obtenir l'email de l'utilisateur à partir de l'API GitHub
	emailsURL := "https://api.github.com/user/emails"
	response, err = client.Get(emailsURL)
	if err != nil {
		fmt.Fprintf(w, "Échec de la requête pour obtenir l'email de l'utilisateur : %v", err)
		return
	}
	defer response.Body.Close()

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}

	err = json.NewDecoder(response.Body).Decode(&emails)
	if err != nil {
		fmt.Fprintf(w, "Échec de la lecture des adresses email : %v", err)
		return
	}

	// Récupérer l'email principal de l'utilisateur (s'il existe)
	var email string
	for _, e := range emails {
		if e.Primary {
			email = e.Email
			break
		}
	}

	password := token.AccessToken

	db, err := sql.Open("sqlite3", "./BDD/BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS User (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Mail TEXT,
			Username TEXT,
			Password TEXT
		)
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	query := "SELECT COUNT(*) FROM User WHERE Mail = ?"
	var count int
	err = db.QueryRow(query, email).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		fmt.Println("Email deja utilisée")
		cookie := http.Cookie{
			Name:  "auth",
			Value: password,
			Path:  "/",
		}

		http.SetCookie(w, &cookie)

		// Redirect to the profile page
		http.Redirect(w, r, "/profil", http.StatusFound)
		return
	}

	query2 := "SELECT COUNT(*) FROM User WHERE Username = ?"
	var count2 int
	err = db.QueryRow(query2, userData.Login).Scan(&count2)
	if err != nil {
		log.Fatal(err)
	}
	if count2 > 0 {
		fmt.Println("Le Pseudo est deja pris !")
		fmt.Fprintf(w, "Pseudo Deja utilisée !")
		return
	}

	insertQuery := `
		INSERT INTO User (Mail, Username, Password)
		VALUES (?, ?, ?)
	`
	_, err = db.Exec(insertQuery, email, userData.Login, password)
	if err != nil {
		log.Fatal(err)
	}

	// Utilisez l'email et le nom d'utilisateur selon vos besoins

	cookie := http.Cookie{
		Name:  "auth",
		Value: password,
		Path:  "/",
	}

	http.SetCookie(w, &cookie)

	// Redirect to the profile page
	http.Redirect(w, r, "/profil", http.StatusFound)

}

func Profil(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprintf(w, "Le cookie 'monCookie' n'a pas été trouvé")
			return
		}
		fmt.Fprintf(w, "Erreur lors de la récupération du cookie : %v", err)
		return
	}

	value := cookie.Value
	fmt.Fprintf(w, "Valeur du cookie 'monCookie' : %s", value)
}

type User struct {
	ID       int
	Mail     string
	Username string
	Password string
}

func GetPasswordByEmail(email string) (string, error) {
	// Établir une connexion à la base de données
	db, err := sql.Open("sqlite3", "./BDD/BDD.db")
	if err != nil {
		return "", fmt.Errorf("erreur lors de la connexion à la base de données: %v", err)
	}
	defer db.Close()

	// Vérifier que la connexion à la base de données est réussie
	err = db.Ping()
	if err != nil {
		return "", fmt.Errorf("erreur lors de la vérification de la connexion à la base de données: %v", err)
	}

	// Effectuer une requête pour récupérer le mot de passe de l'utilisateur
	query := "SELECT Password FROM User WHERE Mail = ?"
	row := db.QueryRow(query, email)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("aucun utilisateur trouvé avec l'e-mail %s", email)
		}
		return "", fmt.Errorf("erreur lors de la récupération du mot de passe: %v", err)
	}

	return hashedPassword, nil
}

func Loginhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Affichage des données reçues
		fmt.Println("Request Body:", string(body))

		// Réponse de succès
		var data struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Failed to decode JSON data", http.StatusInternalServerError)
			return
		}

		// Récupérer le mot de passe haché de la base de données
		hashedPassword, err := GetPasswordByEmail(data.Email)
		fmt.Println("Le mot de passe hashé dans la BDD est", hashedPassword)
		if err != nil {
			errorMessage := "L'email  est invalide"
			errorResponse := ErrorResponse{Message: errorMessage}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		fmt.Println("Le mot de passe hashé dans la BDD est", hashedPassword)

		mdptemp := data.Password
		if err != nil {
			fmt.Println("Erreur lors de l'encryotage du Mdp")
		}

		fmt.Println("Mdp Temps = ", mdptemp, Md5Hash(mdptemp))
		// mdphash := Md5Hash(mdptemp)
		// fmt.Println("Mdphash = ", mdphash, "Hashedpassword = ", hashedPassword)

		// Comparer le mot de passe fourni avec le mot de passe haché
		if hashedPassword != Md5Hash(mdptemp) {
			errorMessage := "Le mot de passe est invalide"
			errorResponse := ErrorResponse{Message: errorMessage}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		// Authentification réussie
		cookie := http.Cookie{
			Name:  "auth",
			Value: hashedPassword,
			Path:  "/",
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/profil", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func Loginpage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./page/Login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
	fmt.Println(r.Method)
}

func UpdateUsername(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode de la requête est POST
	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("auth")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprintf(w, "Le cookie 'monCookie' n'a pas été trouvé")
			return
		}
		fmt.Fprintf(w, "Erreur lors de la récupération du cookie : %v", err)
		return
	}

	password := cookie.Value
	fmt.Fprintf(w, "Valeur du cookie 'monCookie' : %s", password)

	// Lire le corps de la requête
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du corps de la requête", http.StatusInternalServerError)
		return
	}

	// Structure pour stocker les données JSON
	var data struct {
		Username string `json:"newusername"`
	}

	// Décoder les données JSON du corps de la requête
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Erreur lors de la conversion des données JSON", http.StatusBadRequest)
		return
	}

	// Vérifier que les champs requis sont présents
	if data.Username == "" {
		http.Error(w, "Pseudo Requis", http.StatusBadRequest)
		return
	}

	// Établir une connexion à la base de données
	db, err := sql.Open("sqlite3", "./BDD/BDD.db")
	if err != nil {
		http.Error(w, "Erreur lors de la connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Vérifier que la connexion à la base de données est réussie
	err = db.Ping()
	if err != nil {
		http.Error(w, "Erreur lors de la vérification de la connexion à la base de données", http.StatusInternalServerError)
		return
	}

	// Requête pour récupérer l'ID de l'utilisateur basé sur le mot de passe
	query := "SELECT ID FROM User WHERE Password = ?"
	row := db.QueryRow(query, password)

	var userID int
	err = row.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Aucun utilisateur trouvé avec le mot de passe fourni", http.StatusBadRequest)
			return
		}
		http.Error(w, "Erreur lors de la récupération de l'ID de l'utilisateur", http.StatusInternalServerError)
		return
	}

	// Préparer la requête de mise à jour
	updateQuery := "UPDATE User SET Username = ? WHERE ID = ?"
	stmt, err := db.Prepare(updateQuery)
	if err != nil {
		http.Error(w, "Erreur lors de la préparation de la requête de mise à jour", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Exécuter la requête de mise à jour
	_, err = stmt.Exec(data.Username, userID)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution de la requête de mise à jour", http.StatusInternalServerError)
		return
	}

	// Répondre avec succès
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Le nom d'utilisateur a été mis à jour avec succès"))

}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode de la requête est POST
	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Lire le corps de la requête
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du corps de la requête", http.StatusInternalServerError)
		return
	}

	// Structure pour stocker les données JSON
	var data struct {
		Password string `json:"password"`
	}

	// Décoder les données JSON du corps de la requête
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Erreur lors de la conversion des données JSON", http.StatusBadRequest)
		return
	}

	// Vérifier que le champ requis est présent
	if data.Password == "" {
		http.Error(w, "Le champ 'password' est requis", http.StatusBadRequest)
		return
	}

	// Établir une connexion à la base de données
	db, err := sql.Open("sqlite3", "./BDD/BDD.db")
	if err != nil {
		http.Error(w, "Erreur lors de la connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Vérifier que la connexion à la base de données est réussie
	err = db.Ping()
	if err != nil {
		http.Error(w, "Erreur lors de la vérification de la connexion à la base de données", http.StatusInternalServerError)
		return

	}
	cookie, err := r.Cookie("auth")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprintf(w, "Le cookie 'monCookie' n'a pas été trouvé")
			return
		}
		fmt.Fprintf(w, "Erreur lors de la récupération du cookie : %v", err)
		return
	}

	cook := cookie.Value

	var userID int
	query := "SELECT ID FROM User WHERE Password = ? LIMIT 1"
	row := db.QueryRow(query, cook)
	err = row.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Aucun utilisateur correspondant au mot de passe initial fourni", http.StatusNotFound)
		} else {
			http.Error(w, "Erreur lors de la récupération de l'ID de l'utilisateur", http.StatusInternalServerError)
		}
		return
	}

	// Requête pour mettre à jour le mot de passe de l'utilisateur
	updateQuery := "UPDATE User SET Password = ? WHERE ID = ?"
	stmt, err := db.Prepare(updateQuery)
	if err != nil {
		http.Error(w, "Erreur lors de la préparation de la requête de mise à jour", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Exécuter la requête de mise à jour
	_, err = stmt.Exec(data.Password, userID)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution de la requête de mise à jour", http.StatusInternalServerError)
		return
	}

	// Répondre avec succès
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Le mot de passe a été mis à jour avec succès"))
}

func GetDbComments(w http.ResponseWriter, r *http.Request) {
	// Récupérer la catégorie de la recherche depuis la barre de recherche
	idparent := r.FormValue("id")

	// Appeler la fonction postByCateg pour récupérer les posts de la catégorie spécifiée
	comments := commentsByIdPost(idparent)

	// Parcourir les posts et les afficher avec w.Write
	datajson, err := json.Marshal(comments)
	fmt.Println(comments)
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON :", err)
		simulateError()
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(datajson)
}
func commentsByIdPost(idparent string) []PostForm {
	// Connexion à la base de données
	db, err := sql.Open("sqlite3", "./BDD/BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var query string
	query = fmt.Sprintf("SELECT * FROM posts WHERE IdParent = %s", idparent)
	rows, err := db.Query(query)
	fmt.Println(rows)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Récupération des posts et stockage dans une slice
	var comments []PostForm
	for rows.Next() {
		var comment PostForm
		err := rows.Scan(&comment.Id, &comment.IdParent, &comment.Username, &comment.Category, &comment.Title, &comment.Description, &comment.Date)
		if err != nil {
			log.Fatal(err)
		}
		comments = append(comments, comment)

	}
	fmt.Println(comments)

	// Gestion des éventuelles erreurs lors du parcours des résultats
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return comments
}
func AddPost(w http.ResponseWriter, r *http.Request) {
	// Vérification de la méthode HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Lecture du corps de la requête
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Affichage des données reçues
	fmt.Println("Request Body:", string(body))

	var data struct {
		Category    int    `json:"category"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	err = json.Unmarshal(body, &data)
	fmt.Println("debug test1")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	fmt.Println(data)

	db, err := sql.Open("sqlite3", "./BDD/BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Ligne 258")

	username := getConnectedUsername(r)
	fmt.Println("le username est ", username)
	username = "Lucas"

	fmt.Println("Ligne 316")

	// Création de la table "Posts" si elle n'existe pas déjà
	createTablePostsQuery := `
		CREATE TABLE IF NOT EXISTS Posts (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			IdParent INTEGER ,
			Username Text ,
			Category INTEGER,
			Title TEXT,
			Description TEXT,
			Date DATETIME
		)
	`
	fmt.Println("Ligne 331")
	defer db.Close()

	fmt.Println("Ligne 333")
	_, err = db.Exec(createTablePostsQuery)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ligne 340")

	// Création de la table "Toxic" si elle n'existe pas déjà
	createTableToxicQuery := `
		CREATE TABLE IF NOT EXISTS Toxic (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			IdPost INTEGER,
			Username text
		)
	`
	_, err = db.Exec(createTableToxicQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Création de la table "UnLike" si elle n'existe pas déjà
	createTableUnLikeQuery := `
		CREATE TABLE IF NOT EXISTS UnLike (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			IdPost INTEGER,
			Username text
		)
	`
	_, err = db.Exec(createTableUnLikeQuery)
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	currDate := now.Format(DDMMYYYYhhmmss)

	// Insertion des données dans la table "User"
	insertQuery := `
		INSERT INTO Posts (IdParent, Username, Category, Title, Description, Date)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	zero := 0

	_, err = db.Exec(insertQuery, zero, username, data.Category, data.Title, data.Description, currDate)
	if err != nil {
		log.Fatal(err)
	}

	// Réponse de succès
	w.WriteHeader(http.StatusOK)

	// Réponse HTTP réussie
	fmt.Fprintf(w, "Données insérées avec succès dans la base de données.")
	fmt.Println("Ligne 375")

}
func GetDbPosts(w http.ResponseWriter, r *http.Request) {
	// Récupérer la catégorie de la recherche depuis la barre de recherche
	search := r.FormValue("cat")
	order := r.FormValue("order")

	// Appeler la fonction postByCateg pour récupérer les posts de la catégorie spécifiée
	posts := postsByCateg(search, order)

	// Parcourir les posts et les afficher avec w.Write
	datajson, err := json.Marshal(posts)
	fmt.Println("Data json")
	fmt.Println(string(datajson))
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON :", err)
		simulateError()
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(datajson)
}

type Userr struct {
	Username string
}

func getConnectedUsername(r *http.Request) string {
	cookie, err := r.Cookie("auth")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("Le cookie 'auth' n'a pas été trouvé")
			return ""
		}
		fmt.Printf("Erreur lors de la récupération du cookie : %v", err)
		return ""
	}

	fmt.Println("Ligne 270")

	value := cookie.Value
	fmt.Printf("Valeur du cookie 'monCookie' : %s", value)

	password := value

	db, err := sql.Open("sqlite3", "./BDD/BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Préparez la requête avec un paramètre de placeholder pour le mot de passe
	query := "SELECT Username FROM User WHERE Password = ?"
	rows, err := db.Query(query, password)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Parcourir les résultats de la requête
	for rows.Next() {
		var user Userr
		err := rows.Scan(&user.Username)
		if err != nil {
			log.Fatal(err)
		}
		username := user.Username
		return username
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return "Erreur pas de pseudo associé"
}
func postsByCateg(search string, order string) []PostForm {
	// Connexion à la base de données
	db, err := sql.Open("sqlite3", "./BDD/BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Exécution de la requête pour récupérer les posts correspondants à la catégorie
	switch order {
	case "lastcreatpost":
		order = "Date DESC"
		break
	case "firstcreatpost":
		order = "Date ASC"
		break
	case "liked":
		order = "Id ASC"
		break
	}

	var query string
	// if search == "0" {
	// 	query = fmt.Sprintf("SELECT * FROM Posts WHERE IdParent = 0 ORDER BY %s", order)

	// } else {
	// 	query = fmt.Sprintf("SELECT * FROM Posts WHERE Category = %s ORDER BY %s", search, order)
	// }
	if search == "0" {
		query = fmt.Sprintf("SELECT p.*, COUNT(t.Id) as nbtoxic FROM Posts p LEFT JOIN Toxic t ON t.IdPost=p.Id WHERE p.IdParent = 0 ORDER BY p.%s", order)

	} else {
		query = fmt.Sprintf("SELECT p.*, COUNT(t.Id) as nbtoxic FROM Posts p LEFT JOIN Toxic t ON t.IdPost=p.Id WHERE p.Category = %s ORDER BY p.%s", search, order)
	}

	rows, err := db.Query(query)
	fmt.Println("rows", rows)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Récupération des posts et stockage dans une slice
	var posts []PostForm
	for rows.Next() {
		var post PostForm
		//	err := rows.Scan(&post.Id, &post.IdParent, &post.Username, &post.Category, &post.Title, &post.Description, &post.Date)
		err := rows.Scan(&post.Id, &post.IdParent, &post.Username, &post.Category, &post.Title, &post.Description, &post.Date, &post.nbtoxic)

		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
		fmt.Println("post:", post)
	}

	// Gestion des éventuelles erreurs lors du parcours des résultats
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(posts)
	return posts
}
func getPostData(idPost string) PostForm {
	// Connexion à la base de données
	db, err := sql.Open("sqlite3", "./BDD/BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Executer en SQL
	var query string
	query = fmt.Sprintf("SELECT * FROM posts WHERE Id = %s", idPost)
	row, err := db.Query(query)
	fmt.Println(row)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	// Récupération des posts et stockage dans une slice
	var post PostForm
	for row.Next() {
		err := row.Scan(&post.Id, &post.IdParent, &post.Username, &post.Category, &post.Title, &post.Description, &post.Date)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Gestion des éventuelles erreurs lors du parcours des résultats
	if err = row.Err(); err != nil {
		log.Fatal(err)
	}

	return post
}
