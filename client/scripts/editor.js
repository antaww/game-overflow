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
							}).then(r => r.json()).then(r => {
								localStorage.setItem('mention-users', JSON.stringify(r));
								users = JSON.parse(r);
							}).catch(console.error);
						}
						users = users.map(u => {
							return {id: `@${u.username}`, name: u.username};
						}).sort((a, b) => a.id.localeCompare(b.id));

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
		editor.on('change:sate', () => editor.updateSourceElement());
	}).catch(console.error);
});