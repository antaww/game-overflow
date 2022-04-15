const username = document.querySelector('.username_dropdown');
const dropdown = document.querySelector('.dropdown');
const categories = document.querySelector('.categories');
const dropdownCategories = document.querySelector('.dropdown-categories');
const dropdownArrow = document.querySelector('.dropdown-arrow');
let timeout = 300; //temps en ms (doit être identique à la valeur de l'animation css)

username?.addEventListener('click', event => {
    console.log('username clicked');
    if (dropdown.classList.contains('block')) {
        dropdown.classList.remove('block');
        dropdown.classList.add('block-reverse');
        setTimeout(() => dropdown.classList.toggle('none'), timeout);
    } else {
        dropdown.classList.remove('block-reverse');
        dropdown.classList.add('block');
        dropdown.classList.toggle('none');
    }
});

categories.addEventListener('click', event => {
    console.log('categories clicked');
    // note: maybe not working, test in next commit
    if (dropdownArrow.classList.contains("rotate")) {
        dropdownArrow.classList.remove("rotate")
        dropdownArrow.classList.add("rotate-reverse")
    } else {
        dropdownArrow.classList.remove("rotate-reverse")
        dropdownArrow.classList.add("rotate")
    }

    if (dropdownCategories.classList.contains('block')) {
        dropdownCategories.classList.remove('block');
        dropdownCategories.classList.add('block-reverse');
        setTimeout(() => dropdownCategories.classList.toggle('none'), timeout);
    } else {
        dropdownCategories.classList.remove('block-reverse');
        dropdownCategories.classList.add('block');
        dropdownCategories.classList.toggle('none');
    }
});

