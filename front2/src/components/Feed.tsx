import React from "react";
import {FeedRenderer} from "../types";
import {Post} from "./Post";

export const Feed = (props: { data: FeedRenderer }) => {
    return (
        <div>
            <h1>Feed:</h1>
            {props.data.posts?.map(p => <Post key={p.id} data={p}/>)}
        </div>
    );
}
