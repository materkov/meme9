import React from "react";
import styles from "./Composer.module.css";
import {emitCustomEvent} from "react-custom-events";
import {api} from "../store/types";

export function Composer() {
    const [text, setText] = React.useState('');
    const [success, setSuccess] = React.useState(false);
    const [err, setErr] = React.useState(false);

    const onSubmit = () => {
        setSuccess(false);
        setErr(false);

        api("/addPost", {text: text}).then(() => {
            setSuccess(true);

            emitCustomEvent('postCreated', {
                text: text,
            })
        }).catch(() => {
            setErr(true);
        })

        setText('');
    }

    return <>
        <textarea className={styles.text} value={text} onChange={(e) => setText(e.target.value)}/>
        <button className={styles.submit} onClick={onSubmit}>Опубликовать</button>
        {success && <div>Ваш пост сохранен</div>}
        {err && <div>Не удалось сохранить пост</div>}
    </>;
}