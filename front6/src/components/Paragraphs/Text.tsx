import * as types from "../../types/types";
import React from "react";
import * as styles from "./Text.module.css";
import {useArticlePage} from "../../store/articlePage";

export function Text(props: {articleId: string, paragraph: types.ParagraphText}) {
    if (!props.paragraph.text) return null;

    //const setText = useArticlePage((state: ArticlePage) => state.setText);
    const articlePage = useArticlePage();

    const [localText, setLocalText] = React.useState('');
    const onFocus = (e: React.FocusEvent) => {
        const currentText = e.target.textContent || "";

        setLocalText(currentText);
    };

    const onBlur = (e: React.FocusEvent) => {
        const currentText = e.target.textContent || "";

        if (localText !== currentText) {
            articlePage.setText(props.articleId, props.paragraph.id, currentText);
            setLocalText(currentText);
        }
    };


    return (
        <p contentEditable={true} className={styles.text}
           onFocus={onFocus} onBlur={onBlur}
        >{props.paragraph.text}</p>
    )
}
