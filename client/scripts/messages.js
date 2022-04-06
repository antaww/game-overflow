const likeButtons = document.querySelectorAll('.like-btn');
const dislikeButtons = document.querySelectorAll('.dislike-btn');

likeButtons.forEach(button => {
    button.addEventListener('click', (event) => {
        const postId = event.target.closest('.post').getAttribute('data-post-id');
        const url = `/like?id=${postId}`;

        fetch(url, {
            method: 'PUT',
            mode: 'cors'
        }).then(r => r.json()).then(data => {
            const likes = document.querySelector(`[data-post-id="${postId}"] .points`);
            likes.innerText = data.points;
        });
    });
});

dislikeButtons.forEach(button => {
    button.addEventListener('click', (event) => {
        const postId = event.target.closest('.post').getAttribute('data-post-id');
        const url = `/dislike?id=${postId}`;

        fetch(url, {
            method: 'PUT',
            mode: 'cors'
        }).then(r => r.json()).then(data => {
            const likes = document.querySelector(`[data-post-id="${postId}"] .points`);
            likes.innerText = data.points;
        });
    });
});
