import {Query} from "./types";

export type ApiPost = {
    type: "ApiPost";
    id: string;
    text: string;
    userId: string;
}

export type ApiUser = {
    type: "ApiUser";
    id: string;
    name: string;
}

export type ApiFeed = {
    type: "ApiFeed";
    posts: string[];
}

type ApiObject = ApiFeed | ApiPost | ApiUser;

type Store = {
    items: Query;
}

export const store: Store = {
    items: {},
};
