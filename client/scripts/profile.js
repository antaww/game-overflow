import {follow} from "./follow.js";

window.addEventListener('load', () => {
    const followBtn = document.querySelector('.user-follow');
    if (!followBtn) return;

    window.ClassicEditor.create(document.querySelector('#description'), {
        toolbar: ['heading', '|', 'bold', 'strikethrough', 'italic', 'underline', 'code', '|', 'link', 'blockQuote', 'horizontalLine', '|', 'bulletedList', 'numberedList', '|', 'undo', 'redo'],
    }).then(editor => {
        const content = window.localStorage.getItem('editor') || '';
        if (content) editor.setData(content);

        editor.updateSourceElement();
        editor.on('change:sate', () => editor.updateSourceElement());
    }).catch(console.error);

    follow(followBtn);
});

