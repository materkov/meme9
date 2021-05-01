import React from "react";
import * as schema from "../../api/api2";
import {Post} from "../Feed/Post/Post";

export function PostPage(props: { data: schema.PostRenderer }) {
    if (!props.data.post) {
        return null;
    }

    return <Post data={props.data.post}/>
}
