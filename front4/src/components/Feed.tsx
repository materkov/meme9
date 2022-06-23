import React, {useEffect} from "react";
import {Post, QueryParams, User} from "../types";
import {Post as PostTT, PostQuery} from "./post";
import {api} from "../api";
import {PostComposer} from "./PostComposer";

export function Feed() {
    const [viewer, setViewer] = React.useState<User | undefined>();
    const [feed, setFeed] = React.useState<Post[] | undefined>();

    useEffect(() => {
        const feedQuery: QueryParams = {
            feed: {
                inner: PostQuery,
            }
        }

        api(feedQuery).then(data => {
            setFeed(data.feed);
        })
    }, [])

    useEffect(() => {
        if (location.pathname == "/vk-callback") {
            const q: QueryParams = {
                mutation: {
                    inner: {
                        vkAuthCallback: {
                            url: location.href,
                        }
                    }
                }
            }
            api(q).then(result => {
                localStorage.setItem("authToken", result.mutation?.vkAuth?.token || '');
            })
            history.pushState(null, '', '/');
        }

        if (!viewer) {
            const q: QueryParams = {
                viewer: {
                    inner: {
                        name: {}
                    }
                }
            }
            api(q).then(result => {
                setViewer(result.viewer);
            })
        }
    }, [])

    return <>
        <PostComposer/>

        {feed && feed.map(post => {
            return <PostTT post={post} key={post.id}/>
        })}
    </>;
}


