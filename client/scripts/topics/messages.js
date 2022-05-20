function handleLikes() {
	const likeButtons = document.querySelectorAll('.like-btn');
	const dislikeButtons = document.querySelectorAll('.dislike-btn');

	likeButtons.forEach(button => {
		button.addEventListener('click', event => {
			if (event.target.classList.contains('like-color')) {
				event.target.classList.remove('like-color');
			} else {
				event.target.classList.add('like-color');
			}

			event.target.closest('.topic-likes').querySelector('.fa-angle-down').classList.remove('dislike-color');

			const postId = event.target.closest('.topic').getAttribute('data-post-id');
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
			if (event.target.classList.contains('dislike-color')) {
				event.target.classList.remove('dislike-color');
			} else {
				event.target.classList.add('dislike-color');
			}

			event.target.closest('.topic-likes').querySelector('.fa-angle-up').classList.remove('like-color');

			const postId = event.target.closest('.topic').getAttribute('data-post-id');
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
			const id = clicked.closest('.sub-post');
			const text = id.querySelector('md-block.topic-content');
			text.setAttribute('contenteditable', 'true');
			clickedParent.querySelector('.send-edited-form').querySelector('.send-edited-comment').classList.toggle('no-display');
			if (!clickedParent.querySelector('.send-edited-comment').classList.contains('no-display')) {
				text.style.backgroundColor = 'rgba(245, 245, 245, 0.35)';
			}
			clicked.classList.toggle('no-display');
		});
	});

	editedBtn.forEach(btn => {
		btn.addEventListener('click', e => {
			let clicked = e.target;
			const clickedParent = clicked.parentElement.parentElement;
			const id = clicked.closest('.sub-post');
			const text = id.querySelector('md-block.topic-content');
			text.setAttribute('contenteditable', 'false');
			text.style.backgroundColor = 'transparent';
			clickedParent.querySelector('.edit-comment').querySelector('.fa-solid').classList.toggle('no-display');
			clicked.classList.toggle('no-display');

			const messageId = e.target.closest('.send-edited-form').dataset.messageId;
			const topicId = e.target.closest('.send-edited-form').dataset.topicId;
			const url = `/edit-message?idMessage=${messageId}&id=${topicId}`;
			fetch(url, {
				method: 'POST',
				mode: 'cors',
				body: text.innerText,
			});
		});
	});
}

window.addEventListener('load', () => {
	handleLikes();
	editMessage();

	const deleteMessage = document.querySelectorAll('.delete-comment');

	deleteMessage.forEach(element => {
		element.addEventListener('click', e => {
			e.preventDefault();
			const confirmMessage = confirm('Are you sure you want to delete this message ?');
			if (confirmMessage) {
				const url = new URL(window.location.href);
				const messageId = element.parentElement.parentElement.dataset.messageId;
				window.location.href = `/delete-message?idMessage=${messageId}&id=${url.searchParams.get('id')}`;
			}
		});
	});
});
