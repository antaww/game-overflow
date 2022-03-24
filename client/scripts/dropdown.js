const username = document.querySelector(".username")
const dropdown = document.querySelector(".dropdown")

username.addEventListener('click', event => {
    console.log("username clicked")
    dropdown.classList.toggle("block")
})