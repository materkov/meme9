import React from "react";
import styles from "./Page.module.css";
import {Header} from "./Header";

export function Page(props: { children: React.ReactNode }) {
    return <div className={styles.page}>
        <Header/>
        <div>{props.children}</div>
    </div>
}
