import React from "react";
import styles from "./Composer.module.css";
import {api, Edges, Post, uploadApi} from "../store/types";
import {queryClient} from "../store/fetcher";
import {Close} from "./icons/Close";

export function Composer() {
    const [text, setText] = React.useState('');
    const [success, setSuccess] = React.useState(false);
    const [err, setErr] = React.useState(false);

    const [photoAttachData, setPhotoAttachData] = React.useState('');
    const [photoAttachToken, setPhotoAttachToken] = React.useState('');

    const onSubmit = () => {
        if (photoAttachData && !photoAttachToken) {
            return;
        }

        setSuccess(false);
        setErr(false);

        api("/addPost", {text: text, photo: photoAttachToken}).then((resp: Post) => {
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

    const onFileSelected = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (!e.target.files || e.target.files.length < 1) {
            return;
        }

        const file = e.target.files[0];

        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => {
            setPhotoAttachData(String(reader.result));
        }

        uploadApi(file).then(uploadToken => {
            setPhotoAttachToken(uploadToken);
        })
    }

    const removeFile = () => {
        setPhotoAttachToken('');
        setPhotoAttachData('');
    }

    return <>
        <textarea className={styles.text} value={text} onChange={(e) => setText(e.target.value)}/>
        <button className={styles.submit} onClick={onSubmit}>Опубликовать</button>
        <br/>

        {photoAttachData && !photoAttachToken && <div>Загрузка файла...</div>}

        {photoAttachData &&
            <div className={styles.photoAttachContainer}>
                <img className={styles.photoAttach} src={photoAttachData}/>
                <Close onClick={removeFile} className={styles.photoAttachClose}/>
            </div>
        }

        {!photoAttachData &&
            <input type="file" placeholder="Добавить фотографию" onChange={onFileSelected}/>
        }

        {success && <div>Ваш пост сохранен</div>}
        {err && <div>Не удалось сохранить пост</div>}
        <hr/>
    </>;
}