import React from "react";
import styles from "./Composer.module.css";
import {useMutation} from "@tanstack/react-query";
import {addPostMutation} from "../store/addPostMutation";

export function Composer() {
    const [text, setText] = React.useState('');
    const addPostReq = useMutation(addPostMutation(text))

    const onSubmit = () => {
        addPostReq.mutate();
        setText('');
    }

    return <>
        <textarea className={styles.text} value={text} onChange={(e) => setText(e.target.value)}/>
        <button className={styles.submit} onClick={onSubmit}>Опубликовать</button>
    </>;
}