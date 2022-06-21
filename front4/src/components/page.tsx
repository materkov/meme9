import React from "react";
import styles from "./page.module.css";

export function Page(props: { children: React.ReactNode }) {
    return <div className={styles.page}>
        <div className={styles.header}>
            <a href={"/"} className={styles.headerLink}>meme9</a>
        </div>
        <div>{props.children}</div>
    </div>
}
