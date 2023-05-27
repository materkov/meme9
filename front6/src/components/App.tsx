import React from "react";
import {ArticlePage} from "./ArticlePage";
import {Profile} from "./Profile/Profile";
import * as styles from "./App.module.css";
import {useNavigation} from "../store/navigation";

export function App() {
    let page = null;

    const navState = useNavigation(state => state);

    if (navState.url.startsWith("/article/")) {
        page = <ArticlePage/>
    } else if (navState.url.startsWith("/users/")) {
        page = <Profile/>
    } else {
        page = <div>404 page</div>;
    }

    return <div className={styles.app}>
        {page}
    </div>;
}
