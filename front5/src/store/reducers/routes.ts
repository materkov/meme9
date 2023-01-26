import {Global, LoadingState} from "../store";

export interface SetRoute {
    type: 'routes/set'
    url: string
}

export function setRouteReducer(state: Global, data: SetRoute): Global {
    return {
        ...state,
        routing: {
            ...state.routing,
            url: data.url,
        }
    }
}

export interface SetFetchLocker {
    type: 'routes/setFetchLocker';
    key: string;
    state: LoadingState;
}

export function setFetchLocker(state: Global, data: SetFetchLocker): Global {
    return {
        ...state,
        routing: {
            ...state.routing,
            fetchLockers: {
                ...state.routing.fetchLockers,
                [data.key]: data.state,
            }
        }
    }
}
