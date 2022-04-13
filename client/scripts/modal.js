var modals = document.querySelectorAll(".myModal");
var btn = document.querySelectorAll(".modalProfileBtn");
var span = document.querySelectorAll(".close");


btn.forEach(function (element) {
    element.onclick = function (e) {
        let clicked = e.target;
        const modal = [...modals].find(modal => {
            return clicked.parentElement.parentElement.parentElement.querySelector(".topic-user p").innerText === modal.querySelector(".modal-name").innerText;
        });
        modal.classList.add("modal-display");
        console.log("modal opened");
    };
});

span.forEach(function (element) {
    element.onclick = function () {
        modals.forEach(function (modal) {
            modal.classList.remove("modal-display");
            console.log("modal closed");
        });
    };
});

document.addEventListener("click", function (e) {
    modals.forEach(function (modal) {
        if (modal.classList.contains("modal-display") && !e.target.classList.contains("avatar")) {
            if (e.target !== modal.querySelector(".modal-content")) {
                modal.classList.remove("modal-display");
                console.log("modal closed by outside");
            }
        }
    });
});
