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
    const buttons = document.querySelectorAll('.edit-comment');
    const editedBtn = document.querySelectorAll('.send-edited-comment');

    buttons.forEach(btn => {
        btn.addEventListener('click', e => {
            let clicked = e.target;

            const clickedParent = clicked.parentElement.parentElement;
            // get sub-message element
            const id = clicked.closest('.sub-post');
            // get text element
            const text = id.querySelector("p.posts-content");
            // set it to modifiable
            text.setAttribute("contenteditable", "true");
            clickedParent.querySelector('.send-edited-form').querySelector('.send-edited-comment').classList.toggle('no-display');
            clicked.classList.toggle('no-display');

            console.log("after toggle");
        });
    });

    editedBtn.forEach(btn => {
        btn.addEventListener('click', e => {
            let clicked = e.target;
            console.log(clicked);

            const clickedParent = clicked.parentElement.parentElement;
            // get sub-message element
            const id = clicked.closest('.sub-post');
            // get text element
            const text = id.querySelector("p.posts-content");
            // set it to modifiable
            text.setAttribute("contenteditable", "false");
            clickedParent.querySelector('.send-edited-form').querySelector('.send-edited-comment').classList.toggle('no-display');
            clicked.classList.toggle('no-display');

            console.log("after toggle");
        });
    });

}


window.onload = () => {
    handleLikes();
    editMessage();
};