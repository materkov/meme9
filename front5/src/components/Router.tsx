import React from "react";
import {Feed} from "./Feed";
import {UserPage} from "./UserPage";
import {VkCallback} from "./VkCallback";
import {PostPage} from "./PostPage";
import {LoginPage} from "./LoginPage";
import {RegisterPage} from "./RegisterPage";
import {connect} from "react-redux";
import {Global} from "../store/store";

interface Props {
    url: string;
}

function Component(props: Props) {
    const url = props.url;

    if (url == "/") {
        return <Feed/>
    } else if (url.startsWith("/users/")) {
        return <UserPage/>
    } else if (url.startsWith("/posts/")) {
        return <PostPage/>
    } else if (url.startsWith("/register")) {
        return <RegisterPage/>
    } else if (url.startsWith("/login")) {
        return <LoginPage/>
    } else if (url.startsWith("/vk-callback")) {
        return <VkCallback/>
    } else {
        return <>404 page</>;
    }
}

export const Router = connect((state: Global) => {
    return {
        url: state.routing.url,
    }
})(Component);
