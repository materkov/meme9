import * as types from "../../types/types";
import React from "react";
import * as styles from "./Image.module.css";

export function Image(paragraph: types.ParagraphImage) {
    return <div className={styles.imageWrapper}>
        <img className={styles.image} src={paragraph.url}/>
    </div>
}

