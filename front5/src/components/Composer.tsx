import React from "react";
import styles from "./Composer.module.css";
import {emitCustomEvent} from "react-custom-events";
import {api} from "../store/types";

export function Composer() {
    const [text, setText] = React.useState('');

    const onSubmit = () => {
        api("/addPost", {text: text});

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