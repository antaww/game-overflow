const accept = document.querySelector('.accept-cookies');
const decline = document.querySelector('.decline-cookies');

accept.addEventListener('click', async () => {
	await fetch('/cookies?accept=true', {
		method: 'POST',
		credentials: 'same-origin',
	});
	document.querySelector('.cookies').remove();
});

decline.addEventListener('click', async () => {
	await fetch('/cookies?accept=false', {
		method: 'POST',
		credentials: 'same-origin',
	});
	document.querySelector('.cookies').remove();
});