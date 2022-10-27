import React from "react";
import styles from "./Composer.module.css";
import {api, Edges, Post, PostLikeData} from "../store/types";
import {queryClient} from "../store/fetcher";

export function Composer() {
    const [text, setText] = React.useState('');
    const [success, setSuccess] = React.useState(false);
    const [err, setErr] = React.useState(false);

    const onSubmit = () => {
        setSuccess(false);
        setErr(false);

        api("/addPost", {text: text}).then((resp: Post) => {
            setSuccess(true);

            queryClient.setQueryData(["/posts/" + resp.id], resp);

            const feedData = queryClient.getQueryData<Edges>(["/feed"]);
            if (feedData) {
                queryClient.setQueryData(["/feed"], {...feedData, items: [resp.id, ...feedData.items]});
            }

            queryClient.invalidateQueries(["/feed"]);
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