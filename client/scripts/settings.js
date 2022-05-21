const avatarInput = document.querySelector('#avatar-input');
const avatarPreview = document.querySelector('#avatar-preview');
const confirmationInput = document.querySelector('#re-confirm-input');

function checkPassword() {
	fetch('/confirm-password', {
		method: 'PUT',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({
			password: confirmationInput.value,
		}),
	}).then(r => r.json()).then(r => {
		if (r.success) {
			confirmationInput.setCustomValidity('');
		} else {
			confirmationInput.value = '';
			confirmationInput.setCustomValidity('Incorrect password');
		}
	});
}

function selectDefaultColor() {
	const color = document.querySelector('.color-wrapper.customisable');
	const defaultColor = document.querySelector('.color-wrapper.default');

	color.addEventListener('click', () => {
		defaultColor.classList.remove('selected');
		color.classList.add('selected');
	});

	defaultColor.addEventListener('click', e => {
		e.preventDefault();
		color.classList.remove('selected');
		defaultColor.classList.add('selected');
		color.value = defaultColor.value;
	});
}

function resizeTextArea() {
	const textArea = document.querySelector('textarea');
	textArea.addEventListener('input', () => {
		textArea.style.height = 'auto';
		textArea.style.height = textArea.scrollHeight + 'px';
	});
}

function setPasswordConfirmation() {
	confirmationInput.addEventListener('keypress', checkPassword);
}

function updateAvatarPreview() {
	avatarInput.addEventListener('change', () => {
		const reader = new FileReader();
		const file = avatarInput.files[0];

		reader.addEventListener('load', () => {
			avatarPreview.src = reader.result;
		});

		if (file) reader.readAsDataURL(file);
	});
}

window.addEventListener('load', () => {
	resizeTextArea();
	setPasswordConfirmation();
	selectDefaultColor();
	updateAvatarPreview();
});