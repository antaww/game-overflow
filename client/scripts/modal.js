var modals = document.querySelectorAll(".myModal");

var btn = document.querySelectorAll(".modalProfileBtn");

var span = document.querySelectorAll(".close");

// When the user clicks the button, open the modal
btn.forEach(function (element) {
    element.onclick = function (e) {
        let clicked = e.target;
        const modal = [...modals].find(modal => {
            return clicked.parentElement.parentElement.parentElement.querySelector(".topic-user p").innerText === modal.querySelector(".modal-name").innerText;
        });
        modal.style.display = "block";
        modal.style.width = "100%";
        modal.style.height = "100%";
        console.log("modal opened");
    };
});

span.forEach(function (element) {
    element.onclick = function () {
        const modal = [...modals].find(modal => modal.style.display !== "none");
        modal.style.display = "none";
        console.log("modal closed");
    }
})
