function resizeTextArea() {
    const textArea = document.querySelector('textarea');
    textArea.addEventListener('input', () => {
        textArea.style.height = 'auto';
        textArea.style.height = textArea.scrollHeight + 'px';
    });
}

function displayMarkdown() {
    const textArea = document.querySelector('textarea');
    const markdown = document.querySelector('#markdown');
    textArea.addEventListener('input', () => {
        markdown.innerHTML = marked(textArea.value);
    });
}

function handleMarkdownButtons() {
    const boldButton = document.querySelector('#toolbar-bold');
    const italicButton = document.querySelector('#toolbar-italic');
    const underlineButton = document.querySelector('#toolbar-underline');
    const strikeButton = document.querySelector('#toolbar-strikethrough');
    const linkButton = document.querySelector('#toolbar-link');
    const imageButton = document.querySelector('#toolbar-image');
    const codeButton = document.querySelector('#toolbar-code');
    const quoteButton = document.querySelector('#toolbar-quote');
    const listButton = document.querySelector('#toolbar-list');
    const codeBlockButton = document.querySelector('#toolbar-code-block');
    const tableButton = document.querySelector('#toolbar-table');

    document.querySelectorAll('.message-toolbar-btn').forEach(el => {
        el.addEventListener('mousedown', e => e.preventDefault());
    });

    const textArea = document.querySelector('textarea');
    boldButton.addEventListener('click', () => {
        const selection = textArea.value.substring(textArea.selectionStart, textArea.selectionEnd);
        const index = textArea.innerText.indexOf(selection);

        if (selection.length > 0) {
            if (/^\*\*[\s\S]*\*\*$/.test(selection)) {
                textArea.value = textArea.value.replace(selection, selection.substring(2, selection.length - 2));
                textArea.selectionStart = index + 1;
                textArea.selectionEnd = index + selection.length - 2;
            } else {
                const newText = `**${selection}**`;
                textArea.value = textArea.value.replace(selection, newText);
                textArea.selectionStart = index + 1;
                textArea.selectionEnd = index + newText.length + 1;
            }
        } else {
            textArea.value = textArea.value.substring(0, textArea.selectionStart) + '**' + textArea.value.substring(textArea.selectionStart) + '**';
        }
    });

    italicButton.addEventListener('click', () => {
        const selection = textArea.value.substring(textArea.selectionStart, textArea.selectionEnd);
        const index = textArea.innerText.indexOf(selection);

        if (selection.length > 0) {
            if (/^\*[\s\S]*\*$/.test(selection)) {
                textArea.value = textArea.value.replace(selection, selection.substring(1, selection.length - 1));
                textArea.selectionStart = index + 1;
                textArea.selectionEnd = index + selection.length - 2;
            } else {
                const newText = `*${selection}*`;
                textArea.value = textArea.value.replace(selection, newText);
                textArea.selectionStart = index + 2;
                textArea.selectionEnd = index + newText.length + 1;
            }
        } else {
            textArea.value = textArea.value.substring(0, textArea.selectionStart) + '*' + textArea.value.substring(textArea.selectionStart) + '*';
        }
    });

    underlineButton.addEventListener('click', () => {
        const selection = textArea.value.substring(textArea.selectionStart, textArea.selectionEnd);
        const index = textArea.innerText.indexOf(selection);

        if (selection.length > 0) {
            if (/^\_[\s\S]*\_$/.test(selection)) {
                textArea.value = textArea.value.replace(selection, selection.substring(1, selection.length - 1));
                textArea.selectionStart = index + 1;
                textArea.selectionEnd = index + selection.length - 2;
            } else {
                const newText = `_${selection}_`;
                textArea.value = textArea.value.replace(selection, newText);
                textArea.selectionStart = index + 1;
                textArea.selectionEnd = index + newText.length + 1;
            }
        } else {
            textArea.value = textArea.value.substring(0, textArea.selectionStart) + '_' + textArea.value.substring(textArea.selectionStart) + '_';
        }
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
    handleMarkdownButtons();
    checkFields();
    splitTags();
    document.addEventListener('input', checkFields);
};

//create a function that separates tags in the input field when the user presses space or comma
function splitTags() {
    //querySelector input
    const input = document.querySelector('#tags');
    input.addEventListener('keypress', (e) => {
        //if the user presses space or comma
        if (/[,\s]/.test(e.key)) {
            //split the input field into an array
            const tags = [...document.querySelectorAll('.tag')].map(tag => tag.innerText);
            //if input.value is already in tags : input.value = ''
            //else : tags.push(input.value)
            const value = input.value.replace(/[\s,]+/g, '');
            if (value.length === 0) {
                input.value = '';
                return;
            }

            if (tags.includes(value)) {
                input.value = '';
                return;
            }
            console.log(tags)

            //clear the input field
            //loop through the tags array
            //create a new div element
            const div = document.createElement('div');
            //add the class 'tag' to the div
            div.className = 'tag';
            //add the text from the tags array to the div
            div.innerHTML = `${value}<i class="fa-solid fa-xmark tag-cross"></i>`;
            //add the div to the tags div
            input.parentElement.appendChild(div);
            //clear the input field
            input.value = '';
        }

    });

}

document.addEventListener("click", (e) => {
    if (e.target.closest('.tag-cross')) {
        e.target.parentElement.remove();
    }
});
