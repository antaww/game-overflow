window.addEventListener('load', () => {
    const select = document.querySelector('#feed-sorting');

    select?.addEventListener('change', () => {
        const value = select.value;
        const url = new URL(window.location.href);
        url.searchParams.set('s', value);

        window.location.href = url.href;
    });

    const createTopic = document.querySelector('.create-topic');
	if (!createTopic) return;
    window.addEventListener('scroll', (e) => {
		const limit = document.documentElement.scrollHeight - document.documentElement.clientHeight;
		const scrollPosition = window.scrollY;

        if (scrollPosition >= limit - 150) {
            createTopic.style.bottom = `${230 + (scrollPosition - limit)}px`;
        } else {
            createTopic.style.bottom = '5rem';
        }
    });
});