import React, {ChangeEvent, useEffect} from "react";
import {store} from "./store";
import {createRoot} from "react-dom/client";
import styles from "./index.module.css";
import {Post, Query, QueryParams, User} from "./types";

function api(query: QueryParams): Promise<Query> {
    return new Promise((resolve, reject) => {
        let origin = window.location.origin;
        if (origin == "http://localhost:3000") {
            origin = "http://localhost:8000";
        }

        fetch(origin + "/gql", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
            },
            body: JSON.stringify(query),
        })
            .then(data => data.json())
            .then(data => {
                resolve(data as Query);
            })
            .catch(() => {
                reject();
            })
    })
}

function FeedPage() {
    const [viewer, setViewer] = React.useState<User | undefined>();

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
    })
    return <>
        {viewer?.id ? ('Вы авторизованы как ' + viewer.name) : ''}

        <VKAuth/>

        <PostComposer/>
        {store.items.feed && store.items.feed.map(post => {
            return <Post post={post} key={post.id}/>
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
    })

    return <div>
        {!url ? 'Loading ...' : <a href={url}>Авторизоваться через ВК</a>}
    </div>
}

function Post(props: { post: Post }) {
    return <div>
        <div><b>Text: </b> {props.post.text}</div>
        <div><b>User: </b> {props.post.user?.name}</div>
        <hr/>
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

const feedQuery: QueryParams = {
    feed: {
        include: true,
        inner: {
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

api(feedQuery).then(data => {
    store.items = data;

    const root = document.getElementById('root');
    if (root) {
        createRoot(root).render(<FeedPage/>);
    }
})
