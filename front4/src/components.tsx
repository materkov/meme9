import React, {ChangeEvent, useEffect} from "react";
import styles from "./index.module.css";
import {Post, QueryParams, User} from "./types";
import {api} from "./api";
import {Post as PostTT, PostQuery} from "./components/post";

export function FeedPage() {
    const [viewer, setViewer] = React.useState<User | undefined>();
    const [feed, setFeed] = React.useState<Post[] | undefined>();

    useEffect(() => {
        const feedQuery: QueryParams = {
            feed: {
                include: true,
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
                    include: true,
                    inner: {
                        vkAuthCallback: {
                            include: true,
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
                    include: true,
                    inner: {
                        name: {include: true}
                    }
                }
            }
            api(q).then(result => {
                setViewer(result.viewer);
            })
        }
    }, [])
    return <>
        {viewer?.id ? ('Вы авторизованы как ' + viewer.name) : ''}

        <VKAuth/>

        <PostComposer/>
        {feed && feed.map(post => {
            return <PostTT post={post} key={post.id}/>
        })}
    </>;
}

function VKAuth() {
    const [url, setURL] = React.useState('');
    useEffect(() => {
        if (url) {
            return
        }

        const query: QueryParams = {
            vkAuthUrl: {
                include: true
            }
        };
        api(query).then(result => {
            setURL(result.vkAuthUrl || "");
        })
    }, [])

    return <div>
        {!url ? 'Loading ...' : <a href={url}>Авторизоваться через ВК</a>}
    </div>
}


function PostComposer() {
    const [text, setText] = React.useState('');

    const onChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
        setText(e.target.value);
    }

    const onClick = () => {
        const addPostQuery: QueryParams = {
            mutation: {
                include: true,
                inner: {
                    addPost: {
                        include: true,
                        text: text,
                    }
                }
            }
        }
        api(addPostQuery).then(result => {
            alert('DONE');
        })
        setText('');
    };

    return <div>
        <textarea className={styles.PostArea} value={text} onChange={onChange}/>
        <button onClick={onClick}>Отправить</button>
    </div>
}

