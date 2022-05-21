import {follow} from "./follow.js";

window.addEventListener('load', () => {
    const markdown = document.querySelector('.markdown');
    const innerHTML = markdown.innerHTML.replace(/&lt;/g, '<').replace(/&gt;/g, '>');
    markdown.innerHTML = window.marked.parse(innerHTML);

    const followBtn = document.querySelector('.user-follow');
    if (!followBtn) return;

    follow(followBtn);
});

