let msToSeconds = 1000;
let mouseMove = false;

function setSelfOnline() {
    const element = document.querySelector(".circle");
    if (mouseMove) {
        element.classList.add("connected");
        element.classList.remove("disconnected");
    } else {
        element.classList.add("disconnected");
        element.classList.remove("connected");
    }

    sendUserStatus(mouseMove);
}

function sendUserStatus(status) {
    fetch("http://localhost:8091/is-active", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        mode: 'cors',
        body: JSON.stringify({
            isOnline: status
        })
    })
        .catch(error => console.error('Error:', error));
}

function setUsersOnline() {
    const users = document.querySelectorAll(".topic-user");
    const uniqueUsers = new Set([...users].map(user => user.innerText));

    const body = JSON.stringify({
        users: [...uniqueUsers]
    });

    fetch("http://localhost:8091/users-active", {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        mode: 'cors',
        body: body
    }).then(response => response.json())
        .then(data => {
            data?.forEach(user => {
                const users = document.querySelectorAll(".topic-user");
                let userElements = [...users].filter(e => e.innerText === user.username).map(e => e.parentElement.querySelector(".circle"));
                userElements.forEach(userElement => {
                    if (!userElement) return;

                    if (user.isOnline) {
                        userElement.classList.remove("disconnected");
                        userElement.classList.add("connected");
                    } else {
                        userElement.classList.remove("connected");
                        userElement.classList.add("disconnected");
                    }
                });
            });
        })
        .catch(error => console.error('Error:', error));
}

document.addEventListener("mousemove", () => {
    clearTimeout(timeout);
    mouseMove = true;
    let mouseInactiveTime = 30 * msToSeconds; //in seconds
    timeout = setTimeout(() => {
        mouseMove = false;
    }, mouseInactiveTime);
});

window.addEventListener('beforeunload', (e) => {
    const session = document.cookie.match(new RegExp('(^| )' + 'session' + '=([^;]+)'));
    const body = {
        isOnline: false,
        sessionId: session[2]
    };
    navigator.sendBeacon("http://localhost:8091/is-active", JSON.stringify(body));
});

document.addEventListener("DOMContentLoaded", () => {
    let IsActiveCheckerDelay = 5 * msToSeconds; //in seconds
    sendUserStatus(true);
    setUsersOnline();

    setInterval(() => {
        setSelfOnline();
        setUsersOnline();
    }, IsActiveCheckerDelay);
});