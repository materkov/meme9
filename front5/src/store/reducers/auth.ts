import {Global} from "../store";

export interface SetToken {
    type: 'auth/setToken';
    token: string;
}

export function setToken(state: Global, data: SetToken): Global {
    return {
        ...state,
        routing: {
            ...state.routing,
            accessToken: data.token,
        }
    }
}
