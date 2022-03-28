const username = document.querySelector(".username")
const dropdown = document.querySelector(".dropdown")
const categories = document.querySelector(".categories")
const dropdown_categories = document.querySelector(".dropdown-categories")

username.addEventListener('click', event => {
    console.log("username clicked")
    dropdown.classList.toggle("block")
})

categories.addEventListener('click', event => {
    console.log("categories clicked")
    dropdown_categories.classList.toggle("block")
})