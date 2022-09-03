import React, {useEffect, useState} from "react";
import {Feed} from "./Feed";
import {PostPage} from "./PostPage";
import {UserPage} from "./UserPage";
import {BrowseResult, Post} from "../store2/types";
import {useCustomEventListener} from "react-custom-events";

const dataCache: { [key: string]: BrowseResult } = {};

export function Router() {
    const [url, setUrl] = React.useState(location.pathname);
    const [data, setData] = useState<BrowseResult>();

    useEffect(() => {
        document.addEventListener('urlChanged', () => {
            setUrl(location.pathname);
        });
    }, []);

    useCustomEventListener('postCreated', (e) => {
        if (data && data.feed) {
            const dataCopy = JSON.parse(JSON.stringify(data));

            const post: Post = {
                id: "1111",
                fromId: "50",
                // @ts-ignore
                text: e.text,
                detailsURL: "/posts/1111",
            }
            dataCopy.feed.nodes?.posts?.push(post);

            dataCopy.feed.posts = ["1111", ...data.feed.posts || []];
            setData(dataCopy);

            fetch("http://localhost:8000/posts.insert", {
                method: 'POST',
                body: JSON.stringify({text: post.text})
            })
                .then(r => r.json())
                .then((r) => {
                })
        }
    })

    useEffect(() => {
        if (dataCache[url]) {
            setData(dataCache[url]);
            return;
        }

        // Some preload
        if (data && data.feed && url.startsWith("/posts/")) {
            setData({
                postPage: {
                    pagePost: url.substring(7),
                    nodes: data.feed.nodes,
                }
            })
        }

        if (data && data.feed && url.startsWith("/users/")) {
            const user = data.feed.nodes?.users?.find(p => p.id == url.substring(7));
            if (user) {
                setData({
                    userPage: {
                        posts: [],
                        pageUser: url.substring(7),
                        nodes: {
                            users: [user]
                        }
                    }
                })
            }
        }

        fetch("http://localhost:8000/browse?url=" + url)
            .then(r => r.json())
            .then((r) => {
                setData(r);
                dataCache[url] = r;
            })
    }, [url]);

    if (!data) return null;
    if (data.feed) return <Feed data={data.feed}/>
    if (data.postPage) return <PostPage data={data.postPage}/>;
    if (data.userPage) return <UserPage data={data.userPage}/>;

    return <>404 page</>;
}
