import {Global} from "../store";

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
