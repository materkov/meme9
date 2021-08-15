import React from "react";
import {UniversalRenderer} from "../types";
import {Feed} from "./Feed";
import {NewPost} from "./NewPost";
import {Post} from "./Post";
import {UserPage} from "./UserPage";
import {VkAuth} from "./VkAuth";

export const UniversalComponent = (props: { data: UniversalRenderer }) => {
    if (props.data.vkAuthRenderer) {
        return <VkAuth data={props.data.vkAuthRenderer}/>
    } else if (props.data.feedRenderer) {
        return <Feed data={props.data.feedRenderer}/>
    } else if (props.data.newPostRenderer) {
        return <NewPost data={props.data.newPostRenderer}/>
    } else if (props.data.postRenderer) {
        return <Post data={props.data.postRenderer}/>
    } else if (props.data.userPageRenderer) {
        return <UserPage data={props.data.userPageRenderer}/>
    } else {
        return null;
    }
}
