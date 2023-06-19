const close = document.querySelector("#close")
const comment = document.querySelectorAll(".comment-logo")
const popup = document.querySelector("#comment-popup")
const toxicLogo = document.querySelectorAll(".toxic-logo")
const redThumb = document.querySelectorAll(".red-thumb")
let numberRedThumb = 0
let numberToxic = 0




close.addEventListener("click",function(){
    popup.style.display = "none"
   
})

comment.forEach(element => {
    element.addEventListener("click",function(){
        popup.style.display = "flex"

       
    })
    
});

toxicLogo.forEach(element => {
    element.classList.add("toxic"+numberToxic)
    element.dataset.value = numberToxic
    numberToxic++
})

redThumb.forEach(element => {
    element.classList.add("red-thumb"+numberRedThumb)
    element.dataset.value = numberRedThumb
    numberRedThumb++
})

toxicLogo.forEach(element => {
    number = element.dataset.value 
    const CorrespondingReaction = document.querySelector(".red-thumb"+number)
    element.addEventListener("click", function() {
            element.classList.toggle("toxic-filter");
            if (CorrespondingReaction.classList.contains("red-thumb-filter")){
                CorrespondingReaction.classList.remove("red-thumb-filter")
            }
        
      });
      
})

redThumb.forEach(element => {
    number = element.dataset.value 
    const CorrespondingReaction = document.querySelector(".toxic"+number)
    element.addEventListener("click", function() {
        
            element.classList.toggle("red-thumb-filter");
            if (CorrespondingReaction.classList.contains("toxic-filter")){
                CorrespondingReaction.classList.remove("toxic-filter");
            }
            
      });
})

console.log(redThumb[0].classList)


$(document).ready(function () {

    var idpost
    idpost = jQuery("#idposthidden").val();
    console.log(idpost)


    displayCommentByPostId(idpost);
});



function displayCommentByPostId(Idpost) {
    console.log("display Idpost : " + Idpost);
    // chargement des donnée en bdd selon la category
    jQuery.ajax({
        url: "/GetDbComments",
        method: "POST",
        data: { id: Idpost },
        dataType: "json",
        success: function (response) {
            console.log(response);
            response.forEach(function (comment) {
                console.log(comment);
                displayComments(comment);
            });
        },
        error: function (xhr, status, error) {
            console.log(error);
        }
    });
}

function displayComments(data) {
    // boucle d'affichage des posts
    console.log("entry js displaycomments")
    // for (var key in data) {
    // Vérifie si la propriété existe dans l'objet
    // if (data.hasOwnProperty(key)) {
    // Affiche la clé et la valeur correspondante
    // console.log(key + ": " + data[key]);
    // console.log(data["Username"]);
    var headerComment = '<div id="information-top-comment"><div id = "username-comment" >' + data["Username"] + '</div></div > ';

    var bodyComment = '<div class="post">' + data["Description"] + '</div>';

    var footerComment = '<div class="bottom-side-post"><div id="reaction-case-comment"><div id="left-side-reaction-comment"> ';
    footerComment += '<div class="reaction-section-comment"><img src="../static/images/grey-toxic-logo-button.png" alt="" id="toxic-logo"><div id="toxic-nbr-comment">11</div></div>';
    footerComment += '<div class="reaction-section-comment"><img src="../static/images/greythumb-up-logo.png" alt="" id="red-thumb"><div id="red-thumb-nbr-comment">11</div></div>';
    footerComment += '<div class="reaction-section-comment"><img src="../static/images/comment-logo.png" alt=""><div id="comment-nbr-comment">11</div></div></div>';


    var comment = ' <div id="comment-case">' + headerComment + bodyComment + footerComment + ' </div>';
    jQuery("#fillcomments").append(comment);
    //}
    //}
}


