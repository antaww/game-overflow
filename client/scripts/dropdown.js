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
    dropdown.classList.replace('rotate', 'rotate-reverse');

    if (dropdownCategories.classList.contains('block')) {
        setTimeout(() => dropdownCategories.classList.toggle('none'), timeout);
    } else {
        dropdownCategories.classList.toggle('none');
    }

    dropdownCategories.classList.replace('block', 'block-reverse');
});

