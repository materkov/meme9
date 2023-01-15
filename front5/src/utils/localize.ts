//import {emitCustomEvent} from "react-custom-events";
import {setOnline} from "../store/onlineManager";

export function localizeCounter(count: number, form1: string, form234: string, form567: string) {
    const mod = count % 10;

    if (mod == 1) {
        return form1;
    } else if (mod == 2 || mod == 3 || mod == 4) {
        return form234;
    } else {
        return form567;
    }
}

export function navigate(url: string) {
    //window.history.pushState(null, '', url);
    //emitCustomEvent('urlChanged');
}

export function authorize(token: string) {
    //emitCustomEvent('onAuthorized');

    setOnline();
}
