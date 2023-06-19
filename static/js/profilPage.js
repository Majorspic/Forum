const changePassword = document.querySelector("#change-password")
const changeUsername = document.querySelector("#change-username")
const popupPassword = document.querySelector("#popup-change-password")
const popupUsername = document.querySelector("#popup-change-username")
const background = document.querySelector("#blurry-background")
const close = document.querySelectorAll(".close")

changePassword.addEventListener("click",function(){
    popupPassword.style.display = "block"
    background.style.display = "block"
   
})

changeUsername.addEventListener("click",function(){
    console.log("fkefjf")
    popupUsername.style.display = "block"
    background.style.display = "block"
   
})

close.forEach(element => {
    element.addEventListener("click",function(){
        popupPassword.style.display = "none"
        popupUsername.style.display = "none"
        background.style.display = "none"
       
    })
});

function changIdent() {
    event.preventDefault();
    let username = document.getElementById("username").value;

    // Les mots de passe correspondent, affichage des données dans la console
    console.log("username:", username);

    // Envoi des données à la route Go nommée "/login/handler"
    fetch("/login/handler", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            username: username,
        })
    })
    .then(response => {
        if (response.ok) {
            alert("Votre Username a bien été changé") // Rediriger vers l'URL de redirection
        } else {
            response.json().then(data => {
                alert(data.message);
            });
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert("Une erreur s'est produite lors de l'envoi des données");
    });

    return false; // Empêche le formulaire de se soumettre normalement
}

function changPass() {
    event.preventDefault();

    let password = document.getElementById("password").value;

    // Les mots de passe correspondent, affichage des données dans la console
    console.log("password:", password);

    // Envoi des données à la route Go nommée "/login/handler"
    fetch("/login/handler", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            password: password
        })
    })
    .then(response => {
        if (response.ok) {
            alert("Votre mot de passe a bien été changé") // Rediriger vers l'URL de redirection
        } else {
            response.json().then(data => {
                alert(data.message);
            });
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert("Une erreur s'est produite lors de l'envoi des données");
    });

    return false; // Empêche le formulaire de se soumettre normalement
}
