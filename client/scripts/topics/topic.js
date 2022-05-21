

window.addEventListener('load', () => {
	const category = document.querySelector('.topic-category').innerText;
	const options = document.querySelectorAll('#change-topic-category option');
	const select = document.querySelector('#change-topic-category');

	options.forEach(option => {
		if (option.value === category.toLowerCase()) {
			option.style.display = 'none';
		}
	});

	select?.addEventListener('change', e => {
		if (e.target.value !== '') {
			e.preventDefault();
			const confirmMessage = confirm(`Are you sure you want to change the category to ${e.target.value}?`);
			if (confirmMessage) {
				const url = new URL(window.location.href);
				const topicId = url.searchParams.get('id');
				window.location.href = `/change-category?id=${topicId}&category=${e.target.value}`;
			}
		}
	});

	document.querySelectorAll('.markdown').forEach(markdown => {
		const innerHTML = markdown.innerHTML.replace(/&lt;/g, '<').replace(/&gt;/g, '>');
		markdown.innerHTML = window.marked.parse(innerHTML);
	});
});