import React from "react";
import * as schema from "../../api/renderer";
import {Feed} from "../Feed/Feed";
import {PostPage} from "../PostPage/PostPage";
import {Profile} from "../Profile/Profile";
import {Login} from "../Login/Login";

export interface Props {
    data: schema.UniversalRenderer;
}

export function UniversalRenderer(props: Props) {
    const renderer = props.data;

    if (renderer.feedRenderer) return <Feed data={renderer.feedRenderer}/>;
    if (renderer.postRenderer) return <PostPage data={renderer.postRenderer}/>
    if (renderer.profileRenderer) return <Profile data={renderer.profileRenderer}/>
    if (renderer.loginPageRenderer) return <Login data={renderer.loginPageRenderer}/>

    return null;
}
