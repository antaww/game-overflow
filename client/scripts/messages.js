const likeButton = document.querySelectorAll('.like-function');

likeButton.forEach(button => {
    button.addEventListener('click', (event) => {
        const postId = event.target.closest('.post').getAttribute('data-post-id');
        const url = `/like?id=${postId}`;
        const method = 'POST';
        fetch(url, {
            method,
            mode: 'cors'
        })
    });
});
