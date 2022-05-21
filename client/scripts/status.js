const isActiveCheckerDelay = 5 * 1000;
const mouseInactiveTime = 30 * 1000;
let mouseMove = false;
let inactiveTimeout;

function sendUserStatus(status) {
	fetch('http://localhost:8091/is-active', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		mode: 'cors',
		body: JSON.stringify({
			isOnline: status,
			session: localStorage.getItem('session'),
		}),
	}).catch(console.error);
}

function setSelfOnline() {
	sendUserStatus(mouseMove);

	const id = document.querySelector('#user-id').getAttribute('data-user-id');

	const elements = [...document.querySelectorAll('.circle')].filter(element => element.getAttribute('data-user-id') === id);

	elements.forEach((item) => {
		if (mouseMove) {
			item.classList.add('connected');
			item.classList.remove('disconnected');
		} else {
			item.classList.add('disconnected');
			item.classList.remove('connected');
		}
	});
}

function setUsersOnline() {
	const users = document.querySelectorAll('.topic-user');
	const uniqueUsers = new Set([...users].map(user => user.innerText));

	const body = JSON.stringify({
		users: [...uniqueUsers],
	});

	fetch('http://localhost:8091/users-active', {
		method: 'PUT',
		headers: {
			'Content-Type': 'application/json',
		},
		mode: 'cors',
		body: body,
	}).then(response => response.json())
		.then(data => {
			data?.forEach(user => {
				const users = document.querySelectorAll('.topic-user');
				let userElements = [...users].filter(e => e.innerText === user.username).map(e => e.parentElement.querySelector('.circle'));
				userElements.forEach(userElement => {
					if (!userElement) return;

					if (user.isOnline) {
						userElement.classList.remove('disconnected');
						userElement.classList.add('connected');
					} else {
						userElement.classList.remove('connected');
						userElement.classList.add('disconnected');
					}
				});
			});
		})
		.catch(error => console.error('Error:', error));
}

document.addEventListener('mousemove', () => {
	clearTimeout(inactiveTimeout);
	mouseMove = true;

	inactiveTimeout = setTimeout(() => mouseMove = false, mouseInactiveTime);
});

window.addEventListener('beforeunload', () => {
	const session = document.cookie.match(new RegExp('(^| )' + 'session' + '=([^;]+)'));
	const body = {
		isOnline: false,
		session: session[2],
	};
	navigator.sendBeacon('http://localhost:8091/is-active', JSON.stringify(body));
});

window.addEventListener('load', () => {
	sendUserStatus(true);
	setUsersOnline();

	setInterval(() => {
		setSelfOnline();
		setUsersOnline();
	}, isActiveCheckerDelay);
});