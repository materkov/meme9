import {api} from "./types";
import {Global, store} from "../store2/store";

const PING_INTERVAL = 179 * 60 * 1000; // 2 min 59 sec

export function setOnline() {
    localStorage.setItem('onlineManagerLastPing', Math.floor(Date.now()).toString())

    const state = store.getState() as Global;
    if (state.routing.accessToken) {
        api('/setOnline', {});
    }
}

export function runForever() {
    const lastPing = parseInt(localStorage.getItem('onlineManagerLastPing') || '0');
    const currentTs = Math.floor(Date.now());

    let nextTime = PING_INTERVAL - (currentTs - lastPing);
    if (nextTime < 1000) {
        nextTime = 1000;
    }

    setTimeout(() => {
        setOnline();
        setInterval(setOnline, PING_INTERVAL);
    }, nextTime);
}
