const username = document.querySelectorAll(".username_dropdown")
const dropdown = document.querySelector(".dropdown")
const categories = document.querySelector(".categories")
const dropdown_categories = document.querySelector(".dropdown_categories")

username.forEach(element => {
    element.addEventListener('click', event => {
        console.log("username clicked")
        dropdown.classList.toggle("block")
    })
})

categories.addEventListener('click', event => {
    console.log("categories clicked")
    dropdown_categories.classList.toggle("block")
})