const modals = document.querySelectorAll('.myModal');
const btn = document.querySelectorAll('.modalProfileBtn');
const span = document.querySelectorAll('.close');


var keys = {37: 1, 38: 1, 39: 1, 40: 1};

function preventDefault(e) {
    e.preventDefault();
}

function preventDefaultForScrollKeys(e) {
    if (keys[e.keyCode]) {
        preventDefault(e);
        return false;
    }
}

// modern Chrome requires { passive: false } when adding event
var supportsPassive = false;
try {
    window.addEventListener("test", null, Object.defineProperty({}, 'passive', {
        get: function () {
            supportsPassive = true;
        }
    }));
} catch (e) {
}

var wheelOpt = supportsPassive ? {passive: false} : false;
var wheelEvent = 'onwheel' in document.createElement('div') ? 'wheel' : 'mousewheel';

// call this to Disable
function disableScroll() {
    window.addEventListener('DOMMouseScroll', preventDefault, false); // older FF
    window.addEventListener(wheelEvent, preventDefault, wheelOpt); // modern desktop
    window.addEventListener('touchmove', preventDefault, wheelOpt); // mobile
    window.addEventListener('keydown', preventDefaultForScrollKeys, false);
}

// call this to Enable
function enableScroll() {
    window.removeEventListener('DOMMouseScroll', preventDefault, false);
    window.removeEventListener(wheelEvent, preventDefault, wheelOpt);
    window.removeEventListener('touchmove', preventDefault, wheelOpt);
    window.removeEventListener('keydown', preventDefaultForScrollKeys, false);
}

//log everytime window is scrolled
window.addEventListener('scroll', e => {
    console.log(window.scrollY);
});

btn.forEach(element => {
    element.addEventListener('click', e => {
        let clicked = e.target;
        const modal = [...modals].find(modal => {
            return clicked.parentElement.parentElement.parentElement.querySelector('.topic-user p').innerText === modal.querySelector('.modal-name').innerText;
        });
        modal.classList.add('modal-display');
        disableScroll();
        console.log('modal opened');
    });
});

span.forEach(element => {
    element.addEventListener('click', () => {
        modals.forEach(modal => {
            modal.classList.remove('modal-display');
            enableScroll();
            console.log('modal closed');
        });
    });
});

document.addEventListener('click', e => {
    modals.forEach(modal => {
        if (modal.classList.contains('modal-display') && !e.target.classList.contains('avatar')) {
            if (e.target.classList.contains('modal-display')) {
                modal.classList.remove('modal-display');
                enableScroll();
                console.log('modal closed by outside');
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
                console.log('modal closed from esc');
            }
        }
    })
});



