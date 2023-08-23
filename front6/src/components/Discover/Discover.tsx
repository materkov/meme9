import React, {useEffect} from "react";
import {useDiscoverPage} from "../../store/discoverPage";
import * as styles from "./Discover.module.css";
import {useGlobals} from "../../store/globals";
import {postsAdd} from "../../store/api";
import {Post} from "../Post/Post";

export function Discover() {
    const discoverState = useDiscoverPage();
    const globalState = useGlobals();

    const [text, setText] = React.useState('');
    const [saving, setSaving] = React.useState(false);

    useEffect(() => {
        discoverState.fetch();
    }, []);

    const post = () => {
        if (!text) {
            return;
        }

        setSaving(true);
        postsAdd({text: text}).then(() => {
            setSaving(false);
            setText('');
            discoverState.refetch();
        })
    };

    return <div>
        <h1>Discover</h1>

        {globalState.viewerId &&
            <div className={styles.newPostContainer}>
                <textarea className={styles.newPost} placeholder="What's new today?" value={text}
                          onChange={e => setText(e.target.value)}/>
                <button disabled={saving} onClick={post}>Post</button>
                <hr/>
            </div>
        }

        {discoverState.posts.map(post => <Post post={post} key={post.id}/>)}
    </div>
}
