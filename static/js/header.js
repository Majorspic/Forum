const hamburgerButton = document.querySelector(".nav-toggler")
const navigation = document.querySelector("#tag-column")

hamburgerButton.addEventListener("click", toggleNav)

function toggleNav(){
    hamburgerButton.classList.toggle("active")
    navigation.classList.toggle("active")
}

const addPostButton = document.querySelector("#add-post")
const loginButton= document.querySelector(".sign-in-button")
const profilButton = document.querySelector("#user-button")

function checkLoggedIn() {
    let cookies = document.cookie.split(";");
  
    for (let i = 0; i < cookies.length; i++) {
      let cookie = cookies[i].trim();
  
      // Vérifiez si le cookie sécurisé de connexion existe
      if (cookie.startsWith("auth=")) {
        return true;
      }
    }
  
    return false;
  }
  
  // Utilisation de la fonction
  if (checkLoggedIn()) {
    addPostButton.style.display = "block"
    loginButton.style.display = "none"
    profilButton.href = "/profil"
    console.log("L'utilisateur est connecté.");
  } else {
    loginButton.style.display = "block"
    addPostButton.style.display = "none"
    profilButton.href = "/login"
    console.log("L'utilisateur n'est pas connecté.");
  }