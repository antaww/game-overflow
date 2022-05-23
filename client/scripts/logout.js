window.addEventListener('load', () => {
    const logout = document.querySelector('.logout-link');
    logout?.addEventListener('click', e => {
        e.preventDefault();
        localStorage.clear();
        window.location.href = logout.href;
    });
});