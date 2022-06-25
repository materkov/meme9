import React, {useEffect} from "react";
import {Post, QueryParams, User} from "../types";
import {Post as PostTT, PostQuery} from "./Post";
import {api} from "../api";
import {PostComposer} from "./PostComposer";
import {Spinner} from "./Spinner";

export function Feed() {
    const [viewer, setViewer] = React.useState<User | undefined>();
    const [feed, setFeed] = React.useState<Post[] | undefined>();
    const [isLoaded, setIsLoaded] = React.useState(false);

    useEffect(() => {
        const feedQuery: QueryParams = {
            feed: {
                inner: PostQuery,
            },
            viewer: {
                inner: {},
            }
        }

        api(feedQuery).then(data => {
            setFeed(data.feed);
            setViewer(data.viewer);
            setIsLoaded(true);
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
                //history.pushState(null, '', '/');
                location.href = '/';
            })
        }
    }, [])

    return <>
        {!isLoaded && <Spinner/>}

        {isLoaded && <>
            {viewer && <>
                <PostComposer/>
                {feed && feed.map(post => {
                    return <PostTT post={post} key={post.id}/>
                })}

                {!feed && <div>Лента новостей пуста</div>}
            </>}

            {!viewer && <div>Авторизуйтесь, чтобы посмотреть ленту</div>}
        </>}
    </>;
}


