import {store} from "../store";
import {SetRoute} from "../reducers";

export function setRoute(url: string) {
    window.history.pushState(null, '', url);
    store.dispatch({type: 'routes/set', url: url} as SetRoute);
}
