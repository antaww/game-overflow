function verifyPasswords() {
    const password = document.querySelector('#password');
    const confirmPassword = document.querySelector('#confirm-password');
    const confirm = document.querySelector('#btn-sign-up');

    password.addEventListener('keyup', () => {
        if (password.value === confirmPassword.value) confirmPassword.setCustomValidity('');
        else confirmPassword.setCustomValidity('Passwords do not match');
    });

    confirmPassword.addEventListener('keyup', () => {
        if (password.value === confirmPassword.value) confirmPassword.setCustomValidity('');
        else confirmPassword.setCustomValidity('Passwords do not match');
    });

    confirm.addEventListener('click', (e) => {
        if (password.value === confirmPassword.value) confirmPassword.setCustomValidity('');
        else confirmPassword.setCustomValidity('Passwords do not match');
    });
}


window.addEventListener('load', () => {
    verifyPasswords();
});