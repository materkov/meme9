import React, {useEffect} from "react";
import {Post, QueryParams} from "./types";
import {api} from "./api";

export function Post(props: { post: Post }) {
    return <div>
        <div><b>Text: </b> {props.post.text}</div>
        <div><b>User: </b> <a href={"/users/" + props.post.user?.id}>{props.post.user?.name}</a></div>
        <div><a href={"/posts/" + props.post.id}>Link</a></div>
        <hr/>
    </div>
}

export function PostPage(props: { id: string }) {
    const [post, setPost] = React.useState<Post | undefined>(undefined);
    useEffect(() => {
        const q: QueryParams = {
            node: {
                include: true,
                id: props.id,
                inner: {
                    onPost: {
                        date: {include: true},
                        text: {include: true},
                        user: {
                            include: true,
                            inner: {
                                name: {include: true},
                            }
                        }
                    }
                }
            }
        }
        api(q).then(data => data.node?.type === "Post" && setPost(data.node))
    }, []);

    return <>{post && <Post post={post}/>}</>
}
