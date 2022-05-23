const avatarInput = document.querySelector('#avatar-input');
const avatarPreview = document.querySelector('#avatar-preview');
const confirmationInput = document.querySelector('#re-confirm-input');

function addWarningChangingCookies() {
	const cookiesInput = document.querySelector('#use-cookies');
	const actualValue = cookiesInput.checked;

	cookiesInput.addEventListener('change', () => {
		if (cookiesInput.checked !== actualValue) {
			cookiesInput.parentElement.querySelector('.warning').classList.remove('hidden');
		} else {
			cookiesInput.parentElement.querySelector('.warning').classList.add('hidden');
		}
	});
}

function handleLinks() {
	let addLinkElement = document.querySelector('.add-link');
	let removeLinkElement = document.querySelectorAll('.remove-link');

	function addLink(e) {
		if (!e.currentTarget.classList.contains('add-link')) return;

		const linksContainer = document.querySelector('.links');
		let links = document.querySelectorAll('.link');

		const linkName = addLinkElement.parentElement.querySelector('input[name="link-name"]');
		const linkUrl = addLinkElement.parentElement.querySelector('input[name="link-url"]');

		if (!linkName.value) {
			linkName.setCustomValidity('Please enter a name for the link.');
			linkName.reportValidity();
			return;
		} if (!linkUrl.value) {
			linkUrl.setCustomValidity('Please enter a URL for the link.');
			linkUrl.reportValidity();
			return;
		}

		linkName.setCustomValidity('');
		linkUrl.setCustomValidity('');

		const newLink = addLinkElement.parentElement.cloneNode(true);
		addLinkElement.querySelector('i').classList.replace('fa-plus', 'fa-trash');
		addLinkElement.classList.replace('add-link', 'remove-link');
		addLinkElement.addEventListener('click', removeLink);
		newLink.querySelectorAll('input').forEach((input) => {
			input.value = ''
			input.setCustomValidity('');
		});
		if (links.length >= 5) return;

		linksContainer.appendChild(newLink);
		addLinkElement = newLink.querySelector('.add-link');
		addLinkElement.addEventListener('click', addLink);
	}

	function removeLink(e) {
		if (!e.currentTarget.classList.contains('remove-link')) return;

		const linksContainer = document.querySelector('.links');
		let links = document.querySelectorAll('.link');

		if (links.length <= 1) return;

		const link = e.currentTarget.parentElement;
		linksContainer.removeChild(link);
		links = document.querySelectorAll('.link');

		if (links.length === 4) {
			const newLinkElement = addLinkElement.parentElement.cloneNode(true);
			newLinkElement.querySelector('i').classList.replace('fa-trash', 'fa-plus');
			newLinkElement.querySelector('button').classList.replace('remove-link', 'add-link');
			newLinkElement.querySelectorAll('input').forEach((input) => {
				input.value = ''
				input.setCustomValidity('');
			});

			linksContainer.appendChild(newLinkElement);
			addLinkElement = newLinkElement.querySelector('.add-link');
			addLinkElement.addEventListener('click', addLink);
		}

		addLinkElement = linksContainer.querySelector('.add-link');
		addLinkElement.addEventListener('click', addLink);

		removeLinkElement = linksContainer.querySelectorAll('.remove-link');
	}

	addLinkElement?.addEventListener('click', addLink);
	removeLinkElement?.forEach((element) => element.addEventListener('click', removeLink));
}

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
	const isInFirefox = navigator.userAgent.toLowerCase().indexOf('firefox') > -1;

	const color = document.querySelector(`.color-wrapper.customisable ${isInFirefox ? 'input' : ''}`);
	const defaultColor = document.querySelector(`.color-wrapper.default ${isInFirefox ? 'input' : ''}`);

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
	addWarningChangingCookies();
	handleLinks();
	setPasswordConfirmation();
	selectDefaultColor();
	updateAvatarPreview();
});