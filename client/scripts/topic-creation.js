function resizeTextArea() {
    const textArea = document.querySelector('textarea');
    textArea.addEventListener('input', () => {
        textArea.style.height = 'auto';
        textArea.style.height = textArea.scrollHeight + 'px';
    });
}

function checkFields() {
    const title = document.querySelector('#title');
    const content = document.querySelector('#content');
    const category = document.querySelectorAll('input[type="radio"]');
    const submit = document.querySelector('#btn-submit');

    let valid = true;
    if (title.value.length < 1) valid = false;
    if (content.value.length < 1) valid = false;

    let checked = false;
    for (let i = 0; i < category.length; i++) {
        if (category[i].checked) checked = true;
    }

    if (!checked) valid = false;

    submit.disabled = !valid;
    submit.querySelector('i').className = valid ? 'fas fa-check' : 'fas fa-times';
}

window.onload = () => {
    resizeTextArea();

    checkFields();
    document.addEventListener('input', checkFields);
};