import React from "react";
import {UserPageRenderer} from "../types";
import {Post} from "./Post";

export const UserPage = (props: { data: UserPageRenderer }) => {
    return (
        <div>
            User page <b>{props.data.userName}</b><br/><br/>
            ID: {props.data.userId}<br/>
            Name: {props.data.userName}<br/>
            FEED<br/><br/>
            {props.data.posts?.map(p => <Post key={p.id} data={p}/>)}
        </div>
    )
}
