const modals = document.querySelectorAll('.myModal');
const btn = document.querySelectorAll('.modalProfileBtn');
const span = document.querySelectorAll('.close');

btn.forEach(element => {
	element.addEventListener('click', e => {
		let clicked = e.target;
		const modal = [...modals].find(modal => {
			return clicked.parentElement.parentElement.parentElement.querySelector('.topic-user p').innerText === modal.querySelector('.modal-name').innerText;
		});
		modal.classList.add('modal-display');
		console.log('modal opened');
	});
});

span.forEach(element => {
	element.addEventListener('click', () => {
		modals.forEach(modal => {
			modal.classList.remove('modal-display');
			console.log('modal closed');
		});
	});
});

document.addEventListener('click', e => {
	modals.forEach(modal => {
		if (modal.classList.contains('modal-display') && !e.target.classList.contains('avatar')) {
			if (e.target !== modal.querySelector('.modal-content')) {
				modal.classList.remove('modal-display');
				console.log('modal closed by outside');
			}
		}
	});
});
