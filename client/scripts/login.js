window.addEventListener('load', () => {
	const form = document.querySelector('form');

	form.addEventListener('submit', async e => {
		e.preventDefault();

		const username = document.querySelector('#username').value;
		const password = document.querySelector('#password').value;

		const data = {
			username,
			password
		};

		const response = await fetch('/login', {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});

		const result = await response.json();
		if (result.success) {
			localStorage.setItem('session', result.session);
			window.location.href = '/';
		} else {
			let errorElement = document.querySelector('.login-error');
			if (!errorElement) {
				errorElement = document.createElement('span');
				errorElement.classList.add('login-error');
				document.querySelector('.input-area').prepend(errorElement);
			}
			errorElement.innerHTML = result.error;
		}
	});
});