import * as types from "../../types/types";
import React from "react";
import * as styles from "./Text.module.css";
import {ArticlePage, useArticlePage} from "../../store/articlePage";

export function Text(paragraph: types.ParagraphText) {
    if (!paragraph.text) return null;

    const setText = useArticlePage((state: ArticlePage) => state.setText);

    const [localText, setLocalText] = React.useState('');
    const onFocus = (e: React.FocusEvent) => {
        const currentText = e.target.textContent || "";

        setLocalText(currentText);
    };

    const onBlur = (e: React.FocusEvent) => {
        const currentText = e.target.textContent || "";

        if (localText !== currentText) {
            setText(paragraph.id, currentText);
            setLocalText(currentText);
        }
    };


    return (
        <p contentEditable={true} className={styles.text}
           onFocus={onFocus} onBlur={onBlur}
        >{paragraph.text}</p>
    )
}
