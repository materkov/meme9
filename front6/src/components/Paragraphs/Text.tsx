import * as types from "../../types/types";
import React from "react";
import * as styles from "./Text.module.css";
import {ArticlePage, useArticlePage} from "../../store/store";

export function Text(paragraph: types.ParagraphText) {
    if (!paragraph.text) return null;

    const setText = useArticlePage((state: ArticlePage) => state.setText);

    const [localText, setLocalText] = React.useState('');
    const onFocus = (e: React.FocusEvent) => {
        setLocalText(e.target.textContent || "");
    };

    const onBlur = (e: React.FocusEvent) => {
        if (localText !== e.target.textContent) {
            setText(paragraph.id, e.target.textContent);
            setLocalText(e.target.textContent || "");
        }
    };


    return (
        <p contentEditable={true} className={styles.text}
           onFocus={onFocus} onBlur={onBlur}
        >{paragraph.text}</p>
    )
}
