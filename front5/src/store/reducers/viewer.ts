import {Global} from "../store";

export interface SetViewer {
    type: 'viewer/set'
    userId: string
}

export function setViewer(state: Global, data: SetViewer): Global {
    return {
        ...state,
        viewer: {
            ...state.viewer,
            isLoaded: true,
            id: data.userId,
        }
    }
}
