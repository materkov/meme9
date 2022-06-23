import React, {ChangeEvent} from "react";
import {QueryParams} from "../types";
import {api} from "../api";
import styles from "./PostComposer.module.css";

export function PostComposer() {
    const [text, setText] = React.useState('');

    const onChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
        setText(e.target.value);
    }

    const onClick = () => {
        const addPostQuery: QueryParams = {
            mutation: {
                inner: {
                    addPost: {
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

    return (
        <div className={styles.composer}>
            <textarea className={styles.postArea} value={text} onChange={onChange} placeholder="Что у вас нового?"/>
            <div className={styles.bottomContainer}>
                <button onClick={onClick} className={styles.submit}>Отправить</button>
            </div>
        </div>
    )
}

