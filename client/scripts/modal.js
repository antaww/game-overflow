import {follow} from "./follow.js";

const modals = document.querySelectorAll('.myModal');
const btn = document.querySelectorAll('.modal-profile-btn');
const span = document.querySelectorAll('.close');

const keys = {37: 1, 38: 1, 39: 1, 40: 1};

function preventDefault(e) {
    e.preventDefault();
}

function preventDefaultForScrollKeys(e) {
    if (keys[e.keyCode]) {
        preventDefault(e);
        return false;
    }
}

const supportsPassive = false;

const wheelOpt = supportsPassive ? {passive: false} : false;
const wheelEvent = 'onwheel' in document.createElement('div') ? 'wheel' : 'mousewheel';

function disableScroll() {
    window.addEventListener('DOMMouseScroll', preventDefault, false);
    window.addEventListener(wheelEvent, preventDefault, wheelOpt);
    window.addEventListener('touchmove', preventDefault, wheelOpt);
    window.addEventListener('keydown', preventDefaultForScrollKeys, false);
}

function enableScroll() {
    window.removeEventListener('DOMMouseScroll', preventDefault, false);
    window.removeEventListener(wheelEvent, preventDefault, wheelOpt);
    window.removeEventListener('touchmove', preventDefault, wheelOpt);
    window.removeEventListener('keydown', preventDefaultForScrollKeys, false);
}

btn.forEach(element => {
    element.addEventListener('click', e => {
        const clicked = e.currentTarget;
        const modal = [...modals].find(modal => {
            const modalUsername = modal.querySelector('.modal-name');
            const clickedUsername = clicked.parentElement.parentElement.parentElement.querySelector('.topic-user p');
            return clickedUsername.textContent.trim() === modalUsername.textContent.trim();
        });
        console.log(modal);
        modal.classList.add('modal-display');
        disableScroll();
    });
});

span.forEach(element => {
    element.addEventListener('click', () => {
        modals.forEach(modal => {
            modal.classList.remove('modal-display');
            enableScroll();
        });
    });
});

document.addEventListener('click', e => {
    modals.forEach(modal => {
        if (modal.classList.contains('modal-display') && !e.target.classList.contains('avatar')) {
            if (e.target.classList.contains('modal-display')) {
                modal.classList.remove('modal-display');
                enableScroll();
            }
        }
    });
});

document.addEventListener('keydown', e => {
    modals.forEach(modal => {
        if (modal.classList.contains('modal-display')) {
            if (e.keyCode === 27) {
                modal.classList.remove('modal-display');
                enableScroll();
            }
        }
    })
});

const followBtn = document.querySelectorAll('.follow-button');
followBtn.forEach(element => {
    follow(element);
});




