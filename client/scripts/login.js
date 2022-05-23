window.addEventListener('load', () => {
	const form = document.querySelector('form');

	form.addEventListener('submit', async (e) => {
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
			alert(result.message);
		}
	});
});