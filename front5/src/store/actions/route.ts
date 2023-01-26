import {store} from "../store";

export function setRoute(url: string) {
    window.history.pushState(null, '', url);
    store.dispatch({type: 'routes/set', url: url});
}
