const topicStatus = document.querySelector('.topic-status');

//if topic status is clicked, pop confirm message saying "Are you sure you want to change status of this topic?"
topicStatus?.addEventListener('click', function (e) {
    e.preventDefault();
    const confirmMessage = confirm('Are you sure you want to close this topic?');
    if (confirmMessage) {
        const url = new URL(window.location.href);
        const topicId = url.searchParams.get('id');
        window.location.href = `/close-topic?id=${topicId}`;
    }
});

const topicStatusClosed = document.querySelector('.topic-status-closed');

topicStatusClosed?.addEventListener('click', function (e) {
    e.preventDefault();
    const confirmMessage = confirm('Are you sure you want to open this topic?');
    if (confirmMessage) {
        const url = new URL(window.location.href);
        const topicId = url.searchParams.get('id');
        window.location.href = `/open-topic?id=${topicId}`;
    }
});

const deleteMessage = document.querySelectorAll('.delete-comment');

// deleteMessage?.addEventListener('click', function (e) {
//     e.preventDefault();
//     const confirmMessage = confirm('Are you sure you want to delete this message ?');
//     if (confirmMessage) {
//         const url = new URL(window.location.href);
//         const idMessage = e.target.closest.getAttribute('idMessage');
//         window.location.href = `/delete-message?idMessage=${idMessage}&id=${url.searchParams.get('id')}`;
//     }
// });

deleteMessage.forEach(element => {
    element.addEventListener('click', e => {
        let clicked = e.target;
        e.preventDefault();
        const confirmMessage = confirm('Are you sure you want to delete this message ?');
        if (confirmMessage) {
            const url = new URL(window.location.href);
            // const idMessage = clicked.closest.getAttribute('idMessage');
            //idMessage gets the value of the idMessage parent attribute of the element that was clicked
            const idMessage = clicked.parentElement.parentElement.getAttribute('idMessage');
            console.log(idMessage);
            window.location.href = `/delete-message?idMessage=${idMessage}&id=${url.searchParams.get('id')}`;
        }
    });
});