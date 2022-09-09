import React from "react";
import styles from "./Composer.module.css";
import {emitCustomEvent} from "react-custom-events";
import {apiHost} from "../store/types";

export function Composer() {
    const [text, setText] = React.useState('');

    const onSubmit = () => {
        const f = new FormData();
        f.set('text', text);

        fetch(apiHost + "/addPost", {
            method: 'POST',
            body: f,
            headers: {
                'authorization': 'Bearer ' + localStorage.getItem('authToken'),
            },
        })
            .then(r => r.json())
            .then(r => {
            })

        emitCustomEvent('postCreated', {
            text: text,
        })
        setText('');
    }

    return <>
        <textarea className={styles.text} value={text} onChange={(e) => setText(e.target.value)}/>
        <button className={styles.submit} onClick={onSubmit}>Опубликовать</button>
    </>;
}