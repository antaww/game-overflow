function handleLikes() {
    const likeButtons = document.querySelectorAll('.like-btn');
    const dislikeButtons = document.querySelectorAll('.dislike-btn');

    likeButtons.forEach(button => {
        button.addEventListener('click', event => {
            const postId = event.target.closest('.post').getAttribute('data-post-id');
            const url = `/like?id=${postId}`;

            fetch(url, {
                method: 'PUT',
                mode: 'cors',
            }).then(r => r.json()).then(data => {
                const likes = document.querySelector(`[data-post-id="${postId}"] .points`);
                likes.innerText = data.points;
            });
        });
    });

    dislikeButtons.forEach(button => {
        button.addEventListener('click', event => {
            const postId = event.target.closest('.post').getAttribute('data-post-id');
            const url = `/dislike?id=${postId}`;

            fetch(url, {
                method: 'PUT',
                mode: 'cors',
            }).then(r => r.json()).then(data => {
                const likes = document.querySelector(`[data-post-id="${postId}"] .points`);
                likes.innerText = data.points;
            });
        });
    });
}

function editMessage() {
    const btn = document.querySelectorAll('.edit-comment');

    console.log("before toggle");

    btn.forEach(element => {
        element.addEventListener('click', e => {
            let clicked = e.target;
            console.log(clicked);
            console.log(clicked.closest('.topic-date').innerText);

            let content = clicked.closest('.posts-content');
            let textarea = content.querySelector('.edit-text');

            clicked.parentElement.parentElement.previousSibling.previousSibling.classList.toggle('no-display');


            // console.log(closestPost);
            // console.log(closestText);
            // closestPost.classList.toggle('no-display');
            // closestText.classList.toggle('no-display');
            console.log("after toggle");
        });
    });

}


window.onload = () => {
    handleLikes();
}