import {ComponentPost} from "./Post";
import React from "react";
import {useQuery} from "@tanstack/react-query";
import {fetcher} from "../store/fetcher";

export function PostPage() {
    const postId = location.pathname.substring(7);
    const {data: post} = useQuery(["/posts/" + postId], fetcher);
    if (!post) {
        return <>Loading ...</>;
    }

    return <>
        <ComponentPost id={postId}/>
    </>

}
