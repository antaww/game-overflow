let category = document.querySelector(".topic-category").innerText;
console.log(category);

let options = document.querySelectorAll("#change-topic-category option");
let select = document.querySelector("#change-topic-category");

options.forEach(option => {
    if (option.value === category.toLowerCase()) {
    option.style.display = "none";
    }
});

select?.addEventListener('change', e => {
    if (e.target.value !== "") {
        e.preventDefault();
        const confirmMessage = confirm(`Are you sure you want to change the category to ${e.target.value}?`);
        if (confirmMessage) {
            const url = new URL(window.location.href);
            const topicId = url.searchParams.get('id');
            window.location.href = `/change-category?id=${topicId}&category=${e.target.value}`;
        }
    }
});