const topicStatus = document.querySelector('.topic-status');

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

const deleteTopic = document.querySelectorAll('.delete-topic');

deleteTopic.forEach(element => {
    element.addEventListener('click', e => {
        let clicked = e.target;
        e.preventDefault();
        const confirmMessage = confirm('Are you sure you want to delete this topic ?');
        if (confirmMessage) {
            const url = new URL(window.location.href);
            window.location.href = `/delete-topic?id=${url.searchParams.get('id')}`;
        }
    });
});