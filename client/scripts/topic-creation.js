function resizeTextArea() {
    const textArea = document.querySelector('textarea');
    textArea.addEventListener('input', () => {
        textArea.style.height = 'auto';
        textArea.style.height = textArea.scrollHeight + 'px';
    });
}

window.onload = () => {
    resizeTextArea();
}