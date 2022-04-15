let msToSeconds = 1000;
let mouseMove = false;

document.addEventListener("mousemove", () => {
    clearTimeout(timeout);
    mouseMove = true;
    let mouseInactiveTime = 30 * msToSeconds; //in seconds
    timeout = setTimeout(function () {
        console.log("Mouse is not moving");
        mouseMove = false;
    }, mouseInactiveTime);
});

document.addEventListener("DOMContentLoaded", () => {
    let IsActiveCheckerDelay = 5 * msToSeconds; //in seconds
    setInterval(() => {

        const element = document.querySelector(".circle");
        if (mouseMove) {
            console.log("Mouse is moving");
            element.classList.add("connected");
            element.classList.remove("disconnected");
        } else {
            console.log("Mouse is not moving");

            element.classList.add("disconnected");
            element.classList.remove("connected");
        }

        fetch("http://localhost:8091/IsActive", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            mode: 'cors',
            body: JSON.stringify({
                isOnline: mouseMove
            })
        })
            /*.then(response => response.json())
            .then(data => {
                if (!mouseMove) {
                    data = false;
                }

                console.log(data);


            })*/
            .catch(error => console.error('Error:', error));
    }, IsActiveCheckerDelay);


});



