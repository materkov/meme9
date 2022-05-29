import React, {ChangeEvent} from "react";
import {store} from "./store";
import {createRoot} from "react-dom/client";
import styles from "./index.module.css";
import {Post, Query, QueryParams} from "./types";

function api(query: QueryParams): Promise<Query> {
    return new Promise((resolve, reject) => {
        fetch("http://127.0.0.1:8000/gql", {
            method: 'POST',
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
    return <>
        <PostComposer/>
        {store.items.feed.map(post => {
            return <Post post={post} key={post.id}/>
        })}
    </>;
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
        userId: 10,
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
