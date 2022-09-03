import React, {useEffect, useState} from "react";
import {Feed} from "./Feed";
import {PostPage} from "./PostPage";
import {UserPage} from "./UserPage";
import {BrowseResult} from "../store2/types";

const dataCache: {[key: string]: BrowseResult} = {};

export function Router() {
    const [url, setUrl] = React.useState(location.pathname);
    const [data, setData] = useState<BrowseResult>();

    useEffect(() => {
        document.addEventListener('urlChanged', () => {
            setUrl(location.pathname);
        });
    }, []);

    useEffect(() => {
        if (dataCache[url]) {
            setData(dataCache[url]);
            return;
        }

        // Some preload
        if (data && data.feed && url.startsWith("/posts/")) {
            const post = data.feed.nodes?.posts?.find(p => p.id == url.substring(7));
            if (post) {
                setData({
                    postPage: {
                        pagePost: url.substring(7),
                        nodes: {
                            posts: [post]
                        }
                    }
                })
            }
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
            .then((r: BrowseResult) => {
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
