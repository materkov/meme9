import React from "react";
import * as schema from "../../api/renderer";
import {Feed} from "../Feed/Feed";
import {PostPage} from "../PostPage/PostPage";
import {Profile} from "../Profile/Profile";
import {Login} from "../Login/Login";
import {Sandbox} from "../Sandbox/Sandbox";

export function UniversalRenderer(props: { data: schema.UniversalRenderer }) {
    const renderer = props.data;

    if (renderer.feedRenderer) return <Feed data={renderer.feedRenderer}/>;
    if (renderer.postRenderer) return <PostPage data={renderer.postRenderer}/>
    if (renderer.profileRenderer) return <Profile data={renderer.profileRenderer}/>
    if (renderer.loginPageRenderer) return <Login data={renderer.loginPageRenderer}/>
    if (renderer.sandboxRenderer) return <Sandbox data={renderer.sandboxRenderer}/>

    return null;
}
