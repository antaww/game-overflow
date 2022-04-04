const likeButton = document.querySelectorAll('.like-function');

likeButton.forEach(button => {
    button.addEventListener('click', (event) => {
        const postId = event.target.parentNode.parentNode.parentNode.parentNode.dataset.postid;
        const isLike = event.target.previousElementSibling == null;
        const url = `/posts/${postId}/like`;
        const method = isLike ? 'PUT' : 'DELETE';
        fetch(url, {
            method: method
        }).then(res => {
            if (res.ok) {
                event.target.innerText = isLike ? 'Unlike' : 'Like';
                event.target.previousElementSibling.innerText = isLike ? '1' : '0';
            } else {
                console.log('error');
            }
        })
    });
});

likeButton.forEach(button => {
    button.addEventListener('click', (event) => {
        const messageId = event.target.parentNode.parentNode.parentNode.parentNode.dataset.postid;
        const isLike = event.target.previousElementSibling == null;
        const url = `/posts/${postId}/like`;
        const method = isLike ? 'PUT' : 'DELETE';
        fetch(url, {
            method: method
        }).then(res => {
            if (res.ok) {
                event.target.innerText = isLike ? 'Unlike' : 'Like';
                event.target.previousElementSibling.innerText = isLike ? '1' : '0';
            } else {
                console.log('error');
            }
        })
    });
});

