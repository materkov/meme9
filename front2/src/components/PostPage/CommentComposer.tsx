import React from "react";
import styles from "./CommentComposer.module.css";
import {CommentComposerRenderer} from "../../api/posts";
import {GlobalStoreContext} from "../../Context";

export const CommentComposer = (props: { data: CommentComposerRenderer }) => {
    const [text, setText] = React.useState('');
    const [isOk, setIsOk] = React.useState(false);
    const [isSaving, setIsSaving] = React.useState(false);
    const store = React.useContext(GlobalStoreContext);

    const onKeydown = (e: React.KeyboardEvent) => {
        if (e.code == "Enter" && e.metaKey) {
            setIsSaving(true);
            setIsOk(false);

            store.addComment({
                text: text,
                postId: props.data.postId,
            })
                .then(r => {
                    setIsOk(true);
                    setIsSaving(false);
                    setText('');
                })
                .catch(() => {
                    setIsSaving(false);
                })
        }
    }

    const onChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        setText(e.target.value);
    }

    return (
        <>
            <textarea className={styles.Composer} placeholder={props.data.placeholder} value={text}
                      onKeyDown={onKeydown} onChange={onChange} disabled={isSaving}/>
            {isOk && <div className={styles.SuccessBox}>Комментарий добавлен</div>}
        </>
    )
}