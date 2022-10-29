import {api} from "./types";

export function setOnline() {
    if (localStorage.getItem('authToken')) {
        api("/setOnline", {});
    }
}

export function runForever() {
    setTimeout(() => {
        setOnline();
        setInterval(setOnline, 2.9 * 60 * 1000);
    }, 500);
}
