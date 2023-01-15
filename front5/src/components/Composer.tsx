import React from "react";
import styles from "./Composer.module.css";
import {api, api2, Post, PostsAdd, uploadApi} from "../store/types";
import {Close} from "./icons/Close";
import {actions} from "../store2/actions";

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

        api2("posts.add", {text: text, photo: photoAttachToken} as PostsAdd).then((resp: Post) => {
            actions.prependFeed(resp);
            setSuccess(true);
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