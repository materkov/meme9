import * as types from "../../store/types";
import {Global} from "../store";

export interface SetPhoto {
    type: 'photos/set'
    photo: types.Photo;
}

export function setPhoto(state: Global, data: SetPhoto): Global {
    return {
        ...state,
        photos: {
            ...state.photos,
            byId: {
                ...state.photos.byId,
                [data.photo.id]: data.photo,
            }
        }
    }
}
