const close = document.querySelector("#close")
const blurry = document.querySelector("#blurry-background")
const popup = document.querySelector("#popup")
let userConnected

function checkLoggedIn() {
  let cookies = document.cookie.split(";");

  for (let i = 0; i < cookies.length; i++) {
    let cookie = cookies[i].trim();

   
    if (cookie.startsWith("auth=")) {
      return true;
    }
  }

  return false;
}

// Utilisation de la fonction
if (checkLoggedIn()) {
    userConnected = true
} else {
  userConnected = false
}


close.addEventListener("click",function(){
    popup.style.display = "none"
    blurry.style.display = "none"
})



function afficherPopup() {
    blurry.style.display = "block"
    popup.style.display = "block"
    
  }
  
  
  var delaiAffichagePopup = 2000;
  
  function demarrerDelai() {
    if (!userConnected){
    setTimeout(afficherPopup, delaiAffichagePopup)
    }
  }
  
  window.addEventListener('load', demarrerDelai);

  