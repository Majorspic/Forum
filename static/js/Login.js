function validateLogin() {
    let email = document.getElementById("email").value;
    let password = document.getElementById("password").value;

    // Les mots de passe correspondent, affichage des données dans la console
    console.log("Email:", email);
    console.log("Password:", password);

    // Envoi des données à la route Go nommée "/login/handler"
    fetch("/login/handler", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            email: email,
            password: password
        })
    })
    .then(response => {
        if (response.ok) {
            window.location.href = "/profil"; // Rediriger vers l'URL de redirection
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
