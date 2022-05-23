window.addEventListener('load', () => {
	const element = document.querySelector('.ck-editor');
	const save = element.dataset.save;

	window.ClassicEditor.create(element, {
		toolbar: ['heading', '|', 'bold', 'strikethrough', 'italic', 'underline', 'code', '|', 'link', 'blockQuote', 'horizontalLine', '|', 'bulletedList', 'numberedList', '|', 'undo', 'redo'],
		mention: {
			feeds: [
				{
					marker: '@',
					feed: queryText => {

						let users;
						if (localStorage.getItem('mention-users')) {
							users = JSON.parse(localStorage.getItem('mention-users'));
						} else {
							fetch('/users', {
								method: 'GET',
								headers: {
									'Content-Type': 'application/json',
								},
							}).then(r => r.json()).then(json => {
								users = [...json].map(u => {
									return {name: u.username, userId: u.id, id: `@${u.username}`};
								});

								localStorage.setItem('mention-users', JSON.stringify(users));
							}).catch(console.error);
						}
						users = users.sort((a, b) => a.id.localeCompare(b.id));

						return users.filter(user => user.id.toLowerCase().includes(queryText.toLowerCase()));
					},
					minimumCharacters: 1,
				},
			],
		},
	}).then(editor => {
		const content = window.localStorage.getItem(save) || '';
		if (content) editor.setData(content);

		window.addEventListener('beforeunload', () => {
			window.localStorage.setItem(save, editor.getData());
		});

		editor.updateSourceElement();
		editor.model.document.on('change:data', () => editor.updateSourceElement());
	}).catch(console.error);
});