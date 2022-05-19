window.addEventListener('load', () => {
	const select = document.querySelector('#feed-sorting');

	select?.addEventListener('change', () => {
		const value = select.value;
		const url = new URL(window.location.href);
		url.searchParams.set('s', value);

		window.location.href = url.href;
	});
});