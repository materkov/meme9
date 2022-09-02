import React, {useEffect, useState} from "react";
import {Feed} from "./Feed";
import {PostPage} from "./PostPage";
import {UserPage} from "./UserPage";
import {BrowseResult} from "../store2/types";

export function Router() {
    const [url, setUrl] = React.useState(location.pathname);
    const [data, setData] = useState<BrowseResult>();

    useEffect(() => {
        document.addEventListener('urlChanged', () => {
            setUrl(location.pathname);
        });
    }, []);

    useEffect(() => {
        fetch("http://localhost:8000/browse?url=" + url)
            .then(r => r.json())
            .then((r: BrowseResult) => {
                setData(r);
            })
    }, [url]);

    if (!data) return null;
    if (data.feed) return <Feed data={data.feed}/>
    if (data.postPage) return <PostPage data={data.postPage}/>;
    if (data.userPage) return <UserPage data={data.userPage}/>;

    return <>404 page</>;
}
