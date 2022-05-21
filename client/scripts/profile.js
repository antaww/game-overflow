import {follow} from "./follow.js";

window.addEventListener('load', () => {
    const followBtn = document.querySelector('.user-follow');
    if (!followBtn) return;

    follow(followBtn);
});

