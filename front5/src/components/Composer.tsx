import React from "react";
import styles from "./Composer.module.css";
import {emitCustomEvent} from "react-custom-events";

export function Composer() {
    const [text, setText] = React.useState('');

    const onSubmit = () => {
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