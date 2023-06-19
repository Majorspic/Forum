// Créer un objet pour stocker les correspondances entre noms et IDs
var nameToIdMap = {};


$(document).ready(function () {
  // Charger les jeux depuis le fichier JSON
  $.getJSON("./static/ress/apijeux.json", function (data) {
    // Convertir JSON en tableau
    var games = data;

    // Remplir l'objet nameToIdMap avec les données du JSON
    for (var i = 0; i < games.length; i++) {
      var gameId = games[i].id;
      var gameName = games[i].name.toLowerCase();
      nameToIdMap[gameName] = gameId;
    }
    // jquery l'autocomplete
    $("#Category").autocomplete({
      minLength: 1, // Auto-complétion
      source: function (request, response) {
        var term = request.term.toLowerCase();
        var matches = [];

        // Parcourir les jeux et trouver des correspondances
        for (var gameName in nameToIdMap) {
          if (gameName.indexOf(term) >= 0) {
            matches.push(gameName);
          }
        }

        // Renvoyer les correspondances
        response(matches);
      },
      select: function (event, ui) {
        var selectedName = ui.item.value;
        var selectedId = nameToIdMap[selectedName];
        console.log("L'ID selectionné est le :", selectedId);
      }
    });
  });
});


function validatePost() {

  event.preventDefault(); // Empêcher la redirection par défaut

  let Title = document.getElementById("Title").value;
  let Category = document.getElementById("Category").value;
  // for (var gameName in nameToIdMap) {
  //   if (gameName == Category) {
  var CategoryId = nameToIdMap[Category];
  //   }
  // }
  let Description = document.getElementById("Description").value;

  // Les mots de passe correspondent, affichage des données dans la console
  console.log("Titre :", Title);
  console.log("category:", CategoryId);
  console.log("description:", Description);


  // Envoi des données à la route Go "AddPost"
  fetch("/AddPost", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      title: Title,
      category: CategoryId,
      description: Description
    })
  })
    .then(response => {
      if (!response.ok) {
        throw new Error("Failed to add post");
      }
      // Vous pouvez afficher un message de succès ici si vous le souhaitez
    })
    .catch(error => {
      console.error('Error:', error);
      // Vous pouvez afficher un message d'erreur ici si vous le souhaitez
    });

  return false; // Empêche le formulaire de se soumettre normalement
}