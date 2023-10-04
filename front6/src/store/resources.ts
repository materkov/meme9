import {create} from "zustand";
import * as types from "../api/api";

export interface Resources {
    users: { [id: string]: types.User }
    setUser: (u: types.User) => void

    posts: { [id: string]: types.Post }
    setPost: (obj: types.Post) => void
}

export const useResources = create<Resources>()((set, get) => ({
    users: {},
    setUser: (obj: types.User) => {
        set({users: {...get().users, [obj.id]: obj}})
    },

    posts: {},
    setPost: (obj: types.Post) => {
        set({posts: {...get().posts, [obj.id]: obj}})
    }
}))
