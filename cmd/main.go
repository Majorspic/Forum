package main

import (
	Forum "Forum/internal"
	"fmt"
	"net/http"
)

const port = ":8080"

func main() {

	http.HandleFunc("/", Forum.Home)
	http.HandleFunc("/Register", Forum.Register)
	http.HandleFunc("/Post", Forum.PostPage)
	http.HandleFunc("/createPost", Forum.CreatePost)
	http.HandleFunc("/GetDbComments", Forum.GetDbComments)
	http.HandleFunc("/AddPost", Forum.AddPost)
	http.HandleFunc("/GetDbPosts", Forum.GetDbPosts)
	http.HandleFunc("/MyData", Forum.MyData)
	http.HandleFunc("/secret", Forum.Secret)
	http.HandleFunc("/testlogin", Forum.Login)
	http.HandleFunc("/logout", Forum.Logout)
	http.HandleFunc("/Handle", Forum.HandleMain)
	http.HandleFunc("/Hlogin", Forum.HandleGoogleLogin)
	http.HandleFunc("/callback", Forum.HandleGoogleCallback)
	http.HandleFunc("/Newrouter", Forum.NewRouter().ServeHTTP)
	http.HandleFunc("/github/login", Forum.HandleGitHubLogin)
	http.HandleFunc("/github/callback", Forum.HandleGitHubCallback)
	http.HandleFunc("/login/handler", Forum.Loginhandler)
	http.HandleFunc("/login", Forum.Loginpage)
	http.HandleFunc("/profil", Forum.ProfilPage)

	fmt.Println("(http://localhost:8080) - Server started on port 8080")
	fileServer := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.ListenAndServe(port, nil)
}
