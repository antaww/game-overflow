const submit = document.querySelector('#settings-form > button');
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
        console.log(r)

        if (r.success) {
            confirmationInput.setCustomValidity('');
        } else {
            confirmationInput.value = '';
            confirmationInput.setCustomValidity('Incorrect password');
        }
    });
}

function setPasswordConfirmation() {
    // submit.addEventListener('click', () => checkPassword());
    confirmationInput.addEventListener('keypress', () => checkPassword());
}

window.onload = () => {
    setPasswordConfirmation();
};