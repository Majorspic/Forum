function validateForm() {
    let mail = document.getElementById("Mail").value;
    let username = document.getElementById("Username").value;
    let password = document.getElementById("Password").value;
    let confirmPassword = document.getElementById("Confirm password").value;
  
    if (password !== confirmPassword) {
      alert("Password et Confirm password ne sont pas les memes");
      return false;
    }
  
    // Les mots de passe correspondent, affichage des données dans la console
    console.log("Mail:", mail);
    console.log("Username:", username);
    console.log("Password:", password);
  
    // Envoi des données à la route Go nommée "MyData"
    // À remplacer par l'URL appropriée pour votre application Go
    fetch("/MyData", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        mail: mail,
        username: username,
        password: password
      })
    })
    .then(response => {
      if (response.ok) {
        window.location.href = response.url; // Rediriger vers l'URL de redirection
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
  



// alert("Test!");
		// fetch("/checkmailusername", {
		// 	method: "POST",
		// 	headers: {
		// 		"Content-Type": "application/json",
		// 	},
		// 	body: JSON.stringify({ Mail: mail, Username: username }),
		// })
		// .then(response => response.json())
		// .then(data => {
		// 	// Gérez la réponse de la fonction Go ici
		// 	console.log(data);
		// })
		// .catch(error => {
		// 	// Gérez les erreurs ici
		// 	console.error(error);
		// });