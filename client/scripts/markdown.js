window.addEventListener('load', () => {
	let users;
	if (localStorage.getItem('mention-users')) {
		users = JSON.parse(localStorage.getItem('mention-users'));
	} else {
		fetch('/users', {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json',
			},
		}).then(r => r.text()).then(r => {
			users = JSON.parse(r).map(u => {
				return {name: u.username, userId: u.id, id: `@${u.username}`};
			});

			localStorage.setItem('mention-users', JSON.stringify(users));
		}).catch(console.error);
	}

	const markdownElements = document.querySelectorAll('.markdown');
	markdownElements.forEach(element => {
		let innerHTML = element.innerHTML.replace(/&lt;/g, '<').replace(/&gt;/g, '>');

		innerHTML = innerHTML.replace(/@(\w+)/g, (match, username) => {
			const user = users.find(u => u.name === username);
			if (user) {
				const url = `/profile?id=${user.userId}`;
				return `<a class="mention" href="${url}" title="${document.location.origin + url}">${match}</a>`;
			}
			return match;
		}).trimEnd().trimStart();

		element.innerHTML = marked.parse(innerHTML);
	});
});