const username = document.querySelector(".username_dropdown")
const dropdown = document.querySelector(".dropdown")
const categories = document.querySelector(".categories")
const dropdown_categories = document.querySelector(".dropdown-categories")
const dropdown_arrow = document.querySelector(".dropdown-arrow")


username?.addEventListener('click', event => {
    console.log("username clicked")
    dropdown.classList.toggle("block")
})

categories.addEventListener('click', event => {
    console.log("categories clicked")
    if (dropdown_arrow.classList.contains("rotate")) {
        dropdown_arrow.classList.remove("rotate")
        dropdown_arrow.classList.add("rotate-reverse")
    } else {
        dropdown_arrow.classList.remove("rotate-reverse")
        dropdown_arrow.classList.add("rotate")
    }
    const timeout = 300; //temps en ms (doit être identique à la valeur de l'animation css)
    if (dropdown_categories.classList.contains("block")) {
        dropdown_categories.classList.remove("block")
        dropdown_categories.classList.add("block-reverse")
        setTimeout(() => dropdown_categories.classList.toggle("none"), timeout)
    } else {
        dropdown_categories.classList.remove("block-reverse")
        dropdown_categories.classList.add("block")
        dropdown_categories.classList.toggle("none")
    }
})

