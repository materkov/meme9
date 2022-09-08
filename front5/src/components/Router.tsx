import React, {useEffect, useState} from "react";
import {AddPostResponse, apiHost, BrowseResult, Post, User} from "../store/types";
import {emitCustomEvent, useCustomEventListener} from "react-custom-events";
import {Feed} from "./Feed";
import {PostPage} from "./PostPage";
import {UserPage} from "./UserPage";
import {Loader} from "./Loader";

const dataCache: { [key: string]: BrowseResult } = {};

const components: { [key: string]: any } = {
    "Feed": Feed,
    "PostPage": PostPage,
    "UserPage": UserPage,
};

export function Router() {
    const [url, setUrl] = React.useState(location.pathname + location.search);
    const [data, setData] = useState<BrowseResult>();

    useCustomEventListener('urlChanged', () => {
        setUrl(location.pathname + location.search);
    })

    useCustomEventListener('postCreated', (e) => {
        if (data && data.componentName == "Feed") {
            const clientId = '__client_' + Math.random();
            const dataCopy = JSON.parse(JSON.stringify(data));

            const post: Post = {
                id: clientId,
                fromId: "324825265",
                // @ts-ignore
                text: e.text,
                detailsURL: "/posts/1111",
            }
            dataCopy.componentData.nodes?.posts?.push(post);

            dataCopy.componentData.posts = [clientId, ...data.componentData.posts || []];
            setData(dataCopy);

            fetch(apiHost + "/browse?url=/posts/add&q=" + encodeURIComponent(JSON.stringify({text: post.text})), {
                headers: {
                    'authorization': localStorage.getItem('authToken') || "",
                }
            })
                .then(r => r.json())
                .then((r: AddPostResponse) => {
                    fetch(apiHost + "/browse?url=" + url)
                        .then(r => r.json())
                        .then((r) => {
                            setData(r);
                            dataCache[url] = r;
                        })
                })
        }
    })

    useEffect(() => {
        if (dataCache[url]) {
            setData(dataCache[url]);
            return;
        }

        // Some preload
        if (data && data.componentName == "Feed" && url.startsWith("/posts/")) {
            setData({
                componentName: "PostPage",
                componentData: [
                    url.substring(7),
                    data.componentData[0].find((post: any) => post.id == url.substring(7)),
                ],
            })
        }

        if (data && data.componentName == "Feed" && url.startsWith("/users/")) {
            let user: any = null;
            for (let post of data.componentData[0]) {
                if (post.from.id == url.substring(7)) {
                    user = post.from;
                }
            }

            if (user) {
                setData({
                    componentName: "UserPage",
                    componentData: [
                        user,
                        [],
                    ]
                })
            }
        }

        let q = '';
        if (url.startsWith('/vk-callback')) {
            q = JSON.stringify({'redirectUri': location.origin + location.pathname});
        }

        fetch(apiHost + "/browse?url=" + encodeURIComponent(url) + '&q=' + encodeURIComponent(q), {
            method: 'GET',
            headers: {
                'authorization': localStorage.getItem('authToken') || "",
            }
        })
            .then(r => r.json())
            .then((r: BrowseResult) => {
                setData(r);
                dataCache[url] = r;

                if (r.vkCallback) {
                    localStorage.setItem("authToken", r.vkCallback.authToken);
                    window.history.pushState(null, '', '/');
                    emitCustomEvent('urlChanged');
                    emitCustomEvent('onAuthorized');
                }
            })
    }, [url]);

    if (!data || !data.componentName) {
        return <Loader/>;
    }

    const C = components[data.componentName];
    if (!C) {
        return <>404 page</>;
    }

    return <C data={data.componentData}/>
}
