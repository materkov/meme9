import * as types from "../../store/types";
import {Global} from "../store";

export interface SetOnline {
    type: 'online/set'
    userId: string;
    online: types.Online
}

export function setOnline(state: Global, data: SetOnline): Global {
    return {
        ...state,
        online: {
            ...state.online,
            byId: {
                ...state.online.byId,
                [data.userId]: data.online,
            }
        }
    }
}
