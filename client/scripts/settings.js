const avatarInput = document.querySelector('#avatar-input');
const avatarPreview = document.querySelector('#avatar-preview');
const confirmationInput = document.querySelector('#re-confirm-input');

function checkPassword() {
    fetch('/confirm-password', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            password: confirmationInput.value,
        }),
    }).then(r => r.json()).then(r => {
        if (r.success) {
            confirmationInput.setCustomValidity('');
        } else {
            confirmationInput.value = '';
            confirmationInput.setCustomValidity('Incorrect password');
        }
    });
}

function setPasswordConfirmation() {
    confirmationInput.addEventListener('keypress', () => checkPassword());
}

function updateAvatarPreview() {
    avatarInput.addEventListener('change', () => {
        const file = avatarInput.files[0];
        const reader = new FileReader();

        reader.addEventListener('load', () => {
            console.log(reader);
            avatarPreview.src = reader.result;
        });

        if (file) reader.readAsDataURL(file);
    });
}

window.onload = () => {
    setPasswordConfirmation();
    updateAvatarPreview();
};