export function follow(btn) {
    btn.addEventListener('click', e => {
        const clicked = e.target;
        const id = clicked.dataset.id;

        if (clicked.classList.contains('follow-btn')) {
            fetch('/follow', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({id})
            }).then(res => {
                if (res.status === 200) {
                    clicked.innerHTML = '<i class="fa-solid fa-xmark"></i> Unfollow';
                    clicked.classList.add('unfollow-btn');
                    clicked.classList.remove('follow-btn');
                }
            });
        } else if (clicked.classList.contains('unfollow-btn')) {
            fetch('/unfollow', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({id: id})
            }).then(res => {
                if (res.status === 200) {
                    clicked.innerText = 'Follow';
                    clicked.classList.add('follow-btn');
                    clicked.classList.remove('unfollow-btn');
                }
            });
        }
    });
}