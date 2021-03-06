const banBtn = document.querySelector('.ban-option');
const unbanBtn = document.querySelector('.unban-option');
const dots = document.querySelector('.fa-ellipsis-vertical');
const options = document.querySelector('.ban-options');

dots?.addEventListener('click', () => {
    dots.classList.toggle('dots-clicked');
    options.classList.toggle('no-display');
});

banBtn?.addEventListener('click', e =>{
    let clicked = e.target;
    e.preventDefault();
    const confirmMessage = confirm('Are you sure you want to ban this user ?');
    if (confirmMessage) {
        const url = new URL(window.location.href);
        window.location.href = `/ban-user?id=${url.searchParams.get('id')}`;
    }
});

unbanBtn?.addEventListener('click', e => {
    let clicked = e.target;
    e.preventDefault();
    const confirmMessage = confirm('Are you sure you want to unban this user ?');
    if (confirmMessage) {
        const url = new URL(window.location.href);
        window.location.href = `/unban-user?id=${url.searchParams.get('id')}`;
    }
});


