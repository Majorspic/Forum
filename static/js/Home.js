$(document).ready(function () {
    // Charger les jeux depuis le fichier JSON
    $.getJSON("./static/ress/apijeux.json", function (data) {
        // Convertir JSON en tableau
        var games = data;

        // Créer un objet pour stocker les correspondances entre noms et IDs
        var nameToIdMap = {};

        // Remplir l'objet nameToIdMap avec les données du JSON
        for (var i = 0; i < games.length; i++) {
            var gameId = games[i].id;
            var gameName = games[i].name.toLowerCase();
            nameToIdMap[gameName] = gameId;
        }

        // jquery l'autocomplete
        $("#forum-search-category").autocomplete({
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
                displayPostByCateg(selectedId, "");
            }
        });
    });

    displayPostByCateg(0, "lastcreatpost");
});



function displayPostByCateg(categId, orderMode) {
    console.log("display Categ : " + categId);
    if (orderMode == "") {
        orderMode = jQuery("#term_key option:selected").val();
    }
    console.log("orderMode : " + orderMode)
    // chargement des donnée en bdd selon la category
    jQuery.ajax({
        url: "/GetDbPosts",
        method: "POST",
        data: { cat: categId, order: orderMode },
        dataType: "json",
        success: function (response) {
            console.log("retour success posthome : ")
            console.log(response);
            displayPosts(response);
        },
        error: function (xhr, status, error) {
            console.log(error);
        }
    });
}

function displayPosts(data) {
    // boucle d'affichage des posts
    console.log("entry js displaypost")
    console.log(data)
    for (var key in data) {
        // Verif si la propriété existe dans l'objet
        if (data.hasOwnProperty(key)) {
            // Affiche la clé et la valeur correspondante
            console.log(key + ": " + data[key]);
            console.log(data[key]);
            var headerPost = '<div class="top-side-post"><div class="post-title">' + data[key].Title + '</div>';
            headerPost += '<div class="post-category">' + data[key].Category + '</div></div>';

            var bodyPost = '<div class="post">' + data[key].Description + '</div>';

            var footerPost = '<div class="bottom-side-post">';
            footerPost += '<div class="post-toxicity"><div class="post-toxicity"><img src="../static/images/green-toxic-logo-button.png">' + data[key].nbtoxic + '</div>';
            footerPost += '<div class="decoration-bar"></div>';
            footerPost += '<div class="post-link"><a href="http://localhost:8080/Post?id=' + data[key].Id + '">more</a></div></div>';

            var post = ' <div id="post_' + data[key].Id + '" class="case" >' + headerPost + bodyPost + footerPost + ' </div>';
            jQuery("#fill").append(post);
        }
    }
}

// correction a faire en front 
// + ajout d'un id sur le class case

// - sortir le form du class fill pour ne pas effacer le form quand on réaffiche les posts

// affichage du js
// - margin entre post 
// - decoration bar ? ne s'affiche pas

