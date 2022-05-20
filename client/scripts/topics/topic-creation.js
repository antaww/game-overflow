function checkFields() {
	const title = document.querySelector('#title');
	const content = document.querySelector('.ck-content');
	const category = document.querySelectorAll('input[type="radio"]');
	const submit = document.querySelector('#btn-submit');

	let valid = true;
	if (title.value.length < 1) valid = false;
	if (content.innerText.length < 1) valid = false;

	let checked = false;
	for (let i = 0; i < category.length; i++) {
		if (category[i].checked) checked = true;
	}

	if (!checked) valid = false;

	submit.disabled = !valid;
	submit.querySelector('i').className = valid ? 'fas fa-check' : 'fas fa-times';
	saveBeforeUnload();
}

function saveBeforeUnload() {
	window.localStorage.setItem('title', document.querySelector('#title').value);
	window.localStorage.setItem('category', document.querySelector('input[type="radio"]:checked').value);
	window.localStorage.setItem('tags', tags.join(' '));
}

function retrieveData() {
	const title = window.localStorage.getItem('title');
	const category = window.localStorage.getItem('category');
	const tags = window.localStorage.getItem('tags');

	if (title) document.querySelector('#title').value = title;
	if (category) document.querySelector(`input[value="${category}"]`).checked = true;
	if (tags) {
		tags.split(' ').forEach(createTag);
	}
}

window.onload = () => {
	window.ClassicEditor.create(document.querySelector('#editor'), {
		toolbar: ['heading', '|', 'bold', 'strikethrough', 'italic', 'underline', 'code', '|', 'link', 'blockQuote', 'horizontalLine', '|', 'bulletedList', 'numberedList', '|', 'undo', 'redo'],
	}).then(editor => {
		const content = window.localStorage.getItem('editor') || '';
		if (content) editor.setData(content);

		window.addEventListener('beforeunload', () => {
			window.localStorage.setItem('editor', editor.getData());
		});
		editor.updateSourceElement();
		editor.on('change:sate', () => editor.updateSourceElement());
	}).catch(error => {
		console.error(error);
	});

	retrieveData();
	splitTags();
	checkFields();
	document.addEventListener('input', checkFields);
	document.addEventListener('beforeunload', () => saveBeforeUnload());

};

const tags = [...document.querySelectorAll('.tag')].map(tag => tag.innerText);

function splitTags() {
	const input = document.querySelector('#tags');
	const inputBox = document.querySelector('.create-topic-tags');
	const splitRegex = /[\s,]+/g;

	inputBox.addEventListener('click', () => input.focus());

	function updateTags() {
		const value = input.value.replace(splitRegex, '');
		if (value.length === 0 || tags.includes(value)) {
			input.value = '';
			return;
		}

		if (tags.length >= 8) {
			input.setCustomValidity('You can only have 8 tags');
			input.value = '';
			input.form.reportValidity();
			return;
		}

		createTag(value);
		input.value = '';
	}

	input.addEventListener('keypress', e => {
		if (tags.length >= 8) {
			input.setCustomValidity('You can only have 8 tags');
			input.form.reportValidity();
			input.value = '';
			return;
		}

		if (splitRegex.test(e.key)) updateTags();
	});
	input.addEventListener('change', (e) => updateTags());
}

function createTag(text) {
	const el = document.createElement('span');
	el.className = 'tag';
	el.innerHTML = `${text}<i class="fa-solid fa-xmark tag-cross"></i>`;
	document.querySelector('.tag-list').appendChild(el);
	tags.push(text);
}

document.querySelector('#btn-submit').addEventListener('click', () => {
	const input = document.querySelector('#tags');
	input.style.visibility = 'hidden';
	input.value = tags.join(',');
});

document.addEventListener('click', e => {
	if (e.target.closest('.tag-cross')) {
		e.target.parentElement.remove();
		tags.splice(tags.indexOf(e.target.parentElement.innerText.replace(/<\/?[^>]+(>|$)/g, '')), 1);
	}
});
